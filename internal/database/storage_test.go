package database

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testDs DataStorable
	tnow   time.Time

	uidS1 = uuid.New()
	uidS2 = uuid.New()
	uidS3 = uuid.New()
	uidU1 = uuid.New()
	uidU2 = uuid.New()

	secretData1 = models.SecretData{ID: uidS1, UserID: uidU1, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}
	secretData2 = models.SecretData{ID: uidS2, UserID: uidU1, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}
	secretData3 = models.SecretData{ID: uidS3, UserID: uidU2, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}

	list1     = []models.SecretData{secretData1, secretData2}
	list2     = []models.SecretData{secretData3}
	testUser1 = models.User{ID: uidU1, Username: "Test User", Password: []byte("Test Password"), Email: "user1@test.com", CreatedAt: &tnow, UpdatedAt: &tnow, SecretData: &list1}
	testUser2 = models.User{ID: uidU2, Username: "Test User", Password: []byte("Test Password"), Email: "user2@test.com", CreatedAt: &tnow, UpdatedAt: &tnow, SecretData: &list2}
)

func TestMain(m *testing.M) {
	testDb := "test.db"
	db, err := gorm.Open(sqlite.Open(testDb), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	l := logrus.New()
	log = logger.NewLogger(l)
	tnow = time.Now()
	testDs = *NewDataStore(log, db)
	err = testDs.Migrate()
	if err != nil {
		log.Fatal("failed to migrate database", err)
	}

	_, err = testDs.AddUser(addContext(context.Background(), uuid.New()), &testUser1)
	_, err = testDs.AddUser(addContext(context.Background(), uuid.New()), &testUser2)
	if err != nil {
		return
	}

	code := m.Run()

	err = os.Remove(testDb)
	if err != nil {
		log.Fatal("failed to remove test database", err)
	}
	os.Exit(code)
}

func TestDataStore_AddSecretData(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	type args struct {
		ctx  context.Context
		data models.SecretData
	}
	tests := []struct {
		name    string
		args    args
		want    *models.SecretData
		wantErr bool
	}{
		{
			name: "Success add secret data",
			args: args{
				ctx: addContext(context.Background(), uuid.Nil),
				data: models.SecretData{
					ID:     id1,
					UserID: uuid.Nil,
					Type:   "CARD",
					Name:   "Test Secret",
					Secret: []byte("Test Secret"),
				},
			},
			want: &models.SecretData{
				ID:     id1,
				UserID: uuid.Nil,
				Type:   "CARD",
				Name:   "Test Secret",
				Secret: []byte("Test Secret"),
			},
			wantErr: false,
		},
		{
			name: "unique id fails",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				data: models.SecretData{
					ID:     uidS1,
					UserID: id2,
					Type:   "CARD",
					Name:   "Test Secret",
					Secret: []byte("Test Secret"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.AddSecretData(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSecretData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataStore_AddUser(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	type args struct {
		ctx  context.Context
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "success create user",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: &models.User{
					ID:       id2,
					Username: "Test User",
					Password: []byte("Test Password"),
					Email:    "testt@test.com",
				},
			},
			want: &models.User{
				ID:         id2,
				Username:   "Test User",
				Password:   []byte("Test Password"),
				Email:      "testt@test.com",
				SecretData: nil,
			},
			wantErr: false,
		},
		{
			name: "unique id error",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: &models.User{
					ID:       uidU1,
					Username: "Test User",
					Password: []byte("Test Password"),
					Email:    "test1@test.com",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "email already exists",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: &models.User{
					ID:       id1,
					Username: "Test User",
					Password: []byte("Test Password"),
					Email:    "user1@test.com",
				},
			},
			want: &models.User{
				ID:         id1,
				Username:   "Test User",
				Password:   []byte("Test Password"),
				Email:      "test@test.com",
				SecretData: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.AddUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, got.ID, tt.want.ID, "AddUser() ID")
				assert.Equal(t, got.Username, tt.want.Username, "AddUser() Username")
				assert.Equal(t, got.Password, tt.want.Password, "AddUser() Password")
				assert.Equal(t, got.Email, tt.want.Email, "AddUser() Email")
				assert.Equal(t, got.SecretData, tt.want.SecretData, "AddUser() SecretData")
				assert.NotNil(t, got.CreatedAt, "AddUser() CreatedAt")
				assert.NotNil(t, got.UpdatedAt, "AddUser() UpdatedAt")
			}
		})
	}
}

func TestDataStore_GetSecretData(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *[]models.SecretData
		wantErr bool
	}{
		{
			name: "get secret data is user exist",
			args: args{
				ctx: addContext(context.Background(), uidU1),
			},
			want:    &[]models.SecretData{secretData1, secretData2},
			wantErr: false,
		},
		{
			name: "get secret data is user does not exist",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
			},
			want:    &[]models.SecretData{},
			wantErr: false,
		},
		{
			name: "get secret data without user context",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.GetSecretData(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSecretData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataStore_GetUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "get user is exist",
			args: args{
				ctx:   addContext(context.Background(), uuid.New()),
				email: "user1@test.com",
			},
			want:    &testUser1,
			wantErr: false,
		},
		{
			name: "get user is not exist",
			args: args{
				ctx:   addContext(context.Background(), uuid.New()),
				email: "t123@test.com",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.GetUser(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, got.ID, tt.want.ID, "AddUser() ID")
				assert.Equal(t, got.Username, tt.want.Username, "AddUser() Username")
				assert.Equal(t, got.Password, tt.want.Password, "AddUser() Password")
				assert.Equal(t, got.Email, tt.want.Email, "AddUser() Email")
				assert.NotNil(t, got.CreatedAt, "AddUser() CreatedAt")
				assert.NotNil(t, got.UpdatedAt, "AddUser() UpdatedAt")
			}
		})
	}
}

func TestDataStore_UpdateSecretData(t *testing.T) {
	type args struct {
		ctx  context.Context
		data models.SecretData
	}
	tests := []struct {
		name    string
		args    args
		want    *models.SecretData
		wantErr bool
	}{
		{
			name: "update only name",
			args: args{
				ctx: addContext(context.Background(), uidU2),
				data: models.SecretData{
					ID:   uidS3,
					Name: "TEST33",
				},
			},
			want: &models.SecretData{
				ID:     uidS3,
				UserID: uidU2,
				Type:   "CARD",
				Name:   "TEST33",
				Secret: []byte("Test Secret"),
			},
			wantErr: false,
		},
		{
			name: "update only secret",
			args: args{
				ctx: addContext(context.Background(), uidU2),
				data: models.SecretData{
					ID:     uidS3,
					Secret: []byte("TEST, TEST, TEST"),
				},
			},
			want: &models.SecretData{
				ID:     uidS3,
				UserID: uidU2,
				Type:   "CARD",
				Name:   "TEST33",
				Secret: []byte("TEST, TEST, TEST"),
			},
			wantErr: false,
		},
		{
			name: "update secret data is exist",
			args: args{
				ctx: addContext(context.Background(), uidU2),
				data: models.SecretData{
					ID:     uidS3,
					UserID: uidU1,
					Type:   "TEST3",
					Name:   "TEST3",
					Secret: []byte("TEST3, TEST3"),
				},
			},
			want: &models.SecretData{
				ID:     uidS3,
				UserID: uidU2,
				Type:   "CARD",
				Name:   "TEST3",
				Secret: []byte("TEST3, TEST3"),
			},
			wantErr: false,
		},
		{
			name: "update secret data is not exist",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				data: models.SecretData{
					ID:     uuid.New(),
					UserID: uidU1,
					Type:   "TEST3",
					Name:   "TEST3",
					Secret: []byte("TEST3, TEST3"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.UpdateSecretData(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateSecretData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataStore_UpdateUser(t *testing.T) {
	tnew := time.Now()
	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "update only password",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: models.User{
					ID:       uidU2,
					Password: []byte("Test Password23"),
				},
			},
			want: &models.User{
				ID:       uidU2,
				Username: "Test User",
				Password: []byte("Test Password23"),
				Email:    "user2@test.com",
			},
			wantErr: false,
		},
		{
			name: "update only username",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: models.User{
					ID:       uidU2,
					Username: "User2",
				},
			},
			want: &models.User{
				ID:       uidU2,
				Username: "User2",
				Password: []byte("Test Password23"),
				Email:    "user2@test.com",
			},
			wantErr: false,
		},
		{
			name: "update user is exist",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: models.User{
					ID:         uidU2,
					Username:   "Test User2",
					Password:   []byte("Test Password2"),
					Email:      "user22@test.com",
					SecretData: &list2,
				},
			},
			want: &models.User{
				ID:         uidU2,
				Username:   "Test User2",
				Password:   []byte("Test Password2"),
				Email:      "user2@test.com",
				SecretData: &list2,
			},
			wantErr: false,
		},
		{
			name: "update user is not exist",
			args: args{
				ctx: addContext(context.Background(), uuid.New()),
				user: models.User{
					ID:         uuid.New(),
					Username:   "Test User2",
					Password:   []byte("Test Password2"),
					Email:      "user22@test.com",
					SecretData: &list2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.UpdateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got.UpdatedAt = &tnew
				got.CreatedAt = &tnew
				tt.want.UpdatedAt = &tnew
				tt.want.CreatedAt = &tnew
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataStore_DeleteSecretData(t *testing.T) {
	type args struct {
		ctx          context.Context
		idSecretData uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "delete secret data is not exist",
			args: args{
				ctx:          addContext(context.Background(), uuid.New()),
				idSecretData: uuid.New(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "delete secret data is exist",
			args: args{
				ctx:          addContext(context.Background(), uidU1),
				idSecretData: uidS2,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.DeleteSecretData(tt.args.ctx, tt.args.idSecretData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteSecretData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataStore_DeleteUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "delete user is not exist",
			args: args{
				ctx:   addContext(context.Background(), uuid.New()),
				email: "222@test.com",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "delete user is exist",
			args: args{
				ctx:   addContext(context.Background(), uuid.New()),
				email: "user2@test.com",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDs.DeleteUser(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func addContext(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, "UserCtx", &models.UserCtx{
		Username: "Test",
		Email:    "test@test.com",
		Id:       userId,
	})
}
