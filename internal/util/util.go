package util

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/sirupsen/logrus"
)

func Encrypt(ctx context.Context, key []byte, data []byte) ([]byte, error) {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"method": "Encrypt",
		"key":    base64.StdEncoding.EncodeToString(key[len(key)-5:]),
	})

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	nonce, err := generateRandom(gcm.NonceSize())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	log.Info("payload encrypted is successfully")
	return ciphertext, nil
}

func Decrypt(ctx context.Context, key []byte, data []byte) ([]byte, error) {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"method": "Encrypt",
		"key":    base64.StdEncoding.EncodeToString(key[len(key)-5:]),
	})

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	nonce := data[:aesgcm.NonceSize()]
	ciphertext := data[aesgcm.NonceSize():]

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("payload decrypted is successfully")
	return decrypted, nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func AddContextUserCtx(ctx context.Context, username string, email string, id uuid.UUID) context.Context {
	return context.WithValue(ctx, "UserCtx", &models.UserCtx{
		Username: username,
		Email:    email,
		Id:       id,
	})
}
