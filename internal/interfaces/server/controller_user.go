package server

import (
	"context"
	"errors"

	"github.com/google/uuid"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Controller) Register(ctx context.Context, user *pb.User) (*pb.JwtToken, error) {
	log := log.WithFields(logrus.Fields{
		"method": "Register",
		"user":   user.Email,
	})

	uc := &models.UserCtx{
		Username: "not register",
		Email:    user.Email,
	}
	ctx = context.WithValue(ctx, "UserCtx", uc)
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("error generating password")
		return nil, status.Errorf(codes.Internal, "error hash password")
	}
	userCtx := util.AddContextUserCtx(ctx, "not register", user.Email, uuid.Nil)
	newUser, err := s.db.AddUser(userCtx, &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		log.WithError(err).Error("could not add user")
		return nil, status.Error(codes.Internal, err.Error())
	}
	jwt, err := as.CreateJwt(newUser)
	if err != nil {
		log.WithError(err).Error("could not create jwt")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.JwtToken{Token: jwt}, nil
}
func (s *Controller) Login(ctx context.Context, user *pb.User) (*pb.JwtToken, error) {
	log := log.WithFields(logrus.Fields{
		"method": "Login",
		"user":   user.Email,
	})

	userCtx := util.AddContextUserCtx(ctx, "not sig in", user.Email, uuid.Nil)
	getUser, err := s.db.GetUser(userCtx, user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		log.WithError(err).Error("Could not get user")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = bcrypt.CompareHashAndPassword(getUser.Password, []byte(user.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Error(codes.Unauthenticated, "Passwords is incorrect")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	jwt, err := as.CreateJwt(getUser)
	if err != nil {
		log.WithError(err).Error("could not create jwt")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.JwtToken{Token: jwt}, nil
}
