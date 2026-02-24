package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/sirupsen/logrus"
)

var (
	once sync.Once
	log  *logger.Logger
	as   *Service
)

type ServiceAuth interface {
	CreateJwt(user *models.User) (string, error)
	CreateUserCtx(token string) (*models.UserCtx, error)
}

type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAuthService(logger *logger.Logger, pathPrivateKey string) (*Service, error) {
	cert, err := getFile(pathPrivateKey)
	if err != nil {
		return nil, err
	}
	blocks, _ := pem.Decode(cert)
	key, err := x509.ParsePKCS1PrivateKey(blocks.Bytes)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		log = logger
		as = &Service{privateKey: key, publicKey: &key.PublicKey}
	})
	return as, nil
}

func (as Service) CreateJwt(user *models.User) (string, error) {
	log := log.WithFields(logrus.Fields{
		"method": "CreateJwt",
		"user":   user.Email,
	})

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"Id":       user.ID,
		"Username": user.Username,
		"Email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString(as.privateKey)
	if err != nil {
		log.WithError(err).Error("error signing token")
		return "", err
	}

	return signedToken, nil
}

func (as Service) CreateUserCtx(token string) (*models.UserCtx, error) {
	log := log.WithFields(logrus.Fields{
		"method": "CreateJwt",
		"token":  token[len(token)-5:],
	})

	jwtToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return as.publicKey, nil
	})
	if err != nil {
		log.WithError(err).Error("error parsing token")
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		log.WithError(err).Error("invalid token")
		return nil, err
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		log.WithError(err).Error("error parsing expiration time")
		return nil, err
	}
	if exp != nil && exp.Before(time.Now()) {
		log.Error("Token expired")
		return nil, ErrTokenExpired
	}

	id, err := uuid.Parse(claims["Id"].(string))
	if err != nil {
		log.WithError(err).Error("error parsing id")
		return nil, err
	}
	return &models.UserCtx{
		Id:       id,
		Username: claims["Username"].(string),
		Email:    claims["Email"].(string),
	}, nil
}

func getFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("error closing file")
		}
	}(file)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

var ErrTokenExpired = errors.New("token expired")
