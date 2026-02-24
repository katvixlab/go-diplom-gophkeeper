package database

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataStorable interface {
	AddUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)

	AddSecretData(ctx context.Context, data models.SecretData) (*models.SecretData, error)
	GetSecretData(ctx context.Context) (*[]models.SecretData, error)
	UpdateSecretData(ctx context.Context, data models.SecretData) (*models.SecretData, error)
	DeleteSecretData(ctx context.Context, idSecretData uuid.UUID) (bool, error)
	Migrate() error
}

var (
	log  *logger.Logger
	once sync.Once
	ds   DataStorable
)

type DataStore struct {
	db *gorm.DB
}

func NewDataStore(logger *logger.Logger, db *gorm.DB) *DataStorable {
	once.Do(func() {
		log = logger
		ds = &DataStore{db: db}
	})
	return &ds
}

func (ds *DataStore) Migrate() error {
	log.Info("migrating database schema")
	return ds.db.AutoMigrate(&models.User{}, &models.SecretData{})
}

func (ds *DataStore) AddUser(ctx context.Context, user *models.User) (*models.User, error) {
	log := log.WithFields(logrus.Fields{
		"method": "AddUser",
		"user":   ctx.Value("UserCtx").(*models.UserCtx).Email,
	})

	log.Info("adding user")
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	tx := ds.db.Create(&user)
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return user, nil
}

func (ds *DataStore) GetUser(ctx context.Context, email string) (*models.User, error) {
	log := log.WithFields(logrus.Fields{
		"method": "GetUser",
		"user":   ctx.Value("UserCtx").(*models.UserCtx).Email,
	})

	log.Info("getting user")
	var user models.User
	tx := ds.db.Where("email = ?", email).Take(&user)
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func (ds *DataStore) DeleteUser(ctx context.Context, email string) (bool, error) {
	log := log.WithFields(logrus.Fields{
		"method": "DeleteUser",
		"user":   ctx.Value("UserCtx").(*models.UserCtx).Email,
	})

	log.Info("deleting user")
	tx := ds.db.Where("email = ?", email).First(&models.User{}).Delete(&models.User{})
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (ds *DataStore) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	log := log.WithFields(logrus.Fields{
		"method": "UpdateUser",
		"user":   ctx.Value("UserCtx").(*models.UserCtx).Email,
	})

	param := make(map[string]interface{})

	if user.Username != "" {
		param["username"] = user.Username
	}
	if user.Password != nil {
		param["password"] = user.Password
	}

	log.Info("updating user")
	tx := ds.db.Model(&user).Clauses(clause.Returning{}).Updates(param).First(&user)
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

func (ds *DataStore) AddSecretData(ctx context.Context, data models.SecretData) (*models.SecretData, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, ErrUserNotFound
	}
	log := log.WithFields(logrus.Fields{
		"method": "AddSecretData",
		"user":   userCtx.Email,
	})
	data.UserID = userCtx.Id
	log.Info("adding secret data")
	tx := ds.db.Create(&data)
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

func (ds *DataStore) GetSecretData(ctx context.Context) (*[]models.SecretData, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, ErrUserNotFound
	}
	log := log.WithFields(logrus.Fields{
		"method": "GetSecretData",
		"user":   userCtx.Email,
	})

	log.Info("getting secret data")
	var dataList []models.SecretData
	tx := ds.db.Where("user_id = ?", userCtx.Id).Find(&dataList)
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &dataList, nil
}

func (ds *DataStore) UpdateSecretData(ctx context.Context, data models.SecretData) (*models.SecretData, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, ErrUserNotFound
	}
	log := log.WithFields(logrus.Fields{
		"method": "UpdateSecretData",
		"user":   userCtx.Email,
	})

	param := make(map[string]interface{})

	if data.Name != "" {
		param["name"] = data.Name
	}
	if data.Secret != nil {
		param["secret"] = data.Secret
	}

	log.Info("updating secret data")
	tx := ds.db.Model(&data).Clauses(clause.Returning{}).Where("id = ?", data.ID).Where("user_id = ?", userCtx.Id).Updates(param).First(&data)
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &data, nil
}

func (ds *DataStore) DeleteSecretData(ctx context.Context, idSecretData uuid.UUID) (bool, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return false, ErrUserNotFound
	}
	log := log.WithFields(logrus.Fields{
		"method": "DeleteSecretData",
		"user":   userCtx.Email,
	})

	log.Info("deleting secret data")
	tx := ds.db.Where("id = ?", idSecretData).Where("user_id = ?", userCtx.Id).First(&models.SecretData{}).Delete(&models.SecretData{})
	if err := tx.Error; err != nil {
		log.Error(err.Error())
		return false, err
	}
	return true, nil
}

var ErrUserNotFound = errors.New("user not found")
