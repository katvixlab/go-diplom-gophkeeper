package server

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/database"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/mocks"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	tnow time.Time

	uidS1 = uuid.New()
	uidS2 = uuid.New()
	uidS3 = uuid.New()
	uidU1 = uuid.New()
	uidU2 = uuid.New()

	secretData1    = models.SecretData{ID: uidS1, UserID: uidU1, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}
	secretData1upd = models.SecretData{ID: uidS1, UserID: uuid.Nil, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}
	secretData2    = models.SecretData{ID: uidS2, UserID: uidU1, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}
	secretData3    = models.SecretData{ID: uidS3, UserID: uidU2, Type: "CARD", Name: "Test Secret", Secret: []byte("Test Secret")}

	list1     = []models.SecretData{secretData1, secretData2}
	list2     = []models.SecretData{secretData3}
	testUser1 = models.User{ID: uidU1, Username: "Test User", Password: []byte("Test Password"), Email: "user1@test.com", CreatedAt: &tnow, UpdatedAt: &tnow, SecretData: &list1}
	testUser2 = models.User{ID: uidU2, Username: "Test User", Password: []byte("Test Password"), Email: "user2@test.com", CreatedAt: &tnow, UpdatedAt: &tnow, SecretData: &list2}

	userCtx1 = addContext(context.Background(), uidU1)
	userCtx2 = addContext(context.Background(), uidU2)

	note1 = pb.Note{Id: uidS1.String(), Name: "Test Secret", Type: "CARD", SecretData: []byte("Test Secret")}
	note2 = pb.Note{Id: uidS2.String(), Name: "Test Secret", Type: "CARD", SecretData: []byte("Test Secret")}
)

func TestMain(m *testing.M) {
	log = logger.NewLogger(logrus.New())
	code := m.Run()
	os.Exit(code)
}

func TestController_AddNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockDataStorable(ctrl)

	m.EXPECT().AddSecretData(userCtx1, secretData1).Return(&secretData1, nil)
	m.EXPECT().AddSecretData(userCtx2, gomock.Any()).Return(nil, database.ErrUserNotFound)

	type fields struct {
		UnimplementedNoteServicesServer pb.UnimplementedNoteServicesServer
		UnimplementedUserServicesServer pb.UnimplementedUserServicesServer
		db                              database.DataStorable
	}
	type args struct {
		ctx  context.Context
		note *pb.Note
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx:  userCtx1,
				note: &note1,
			},
			want:    &empty.Empty{},
			wantErr: false,
		},
		{
			name: "User not found",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx:  userCtx2,
				note: &note2,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong UUID",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx2,
				note: &pb.Note{
					Id:         "wrong",
					Name:       "Test Secret",
					Type:       "CARD",
					SecretData: []byte("Test Secret"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong Ctx",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx:  context.Background(),
				note: &note1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Controller{
				UnimplementedNoteServicesServer: tt.fields.UnimplementedNoteServicesServer,
				UnimplementedUserServicesServer: tt.fields.UnimplementedUserServicesServer,
				db:                              tt.fields.db,
			}
			got, err := s.AddNote(tt.args.ctx, tt.args.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddNote() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_DeleteNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockDataStorable(ctrl)

	m.EXPECT().DeleteSecretData(userCtx1, uidS1).Return(true, nil)
	m.EXPECT().DeleteSecretData(userCtx2, uidS2).Return(false, database.ErrUserNotFound)
	m.EXPECT().DeleteSecretData(userCtx1, uidS2).Return(false, nil)

	type fields struct {
		UnimplementedNoteServicesServer pb.UnimplementedNoteServicesServer
		UnimplementedUserServicesServer pb.UnimplementedUserServicesServer
		db                              database.DataStorable
	}
	type args struct {
		ctx context.Context
		req *pb.NoteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx1,
				req: &pb.NoteRequest{
					IdNote: uidS1.String(),
				},
			},
			want:    &empty.Empty{},
			wantErr: false,
		},
		{
			name: "User not found",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx2,
				req: &pb.NoteRequest{
					IdNote: uidS2.String(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong UUID",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx2,
				req: &pb.NoteRequest{
					IdNote: "UUID",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong Ctx",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.NoteRequest{
					IdNote: uidS1.String(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Note not found",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx1,
				req: &pb.NoteRequest{
					IdNote: uidS2.String(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Controller{
				UnimplementedNoteServicesServer: tt.fields.UnimplementedNoteServicesServer,
				UnimplementedUserServicesServer: tt.fields.UnimplementedUserServicesServer,
				db:                              tt.fields.db,
			}
			got, err := s.DeleteNote(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteNote() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_GetNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockDataStorable(ctrl)

	m.EXPECT().GetSecretData(userCtx1).Return(&[]models.SecretData{secretData1, secretData2}, nil)
	m.EXPECT().GetSecretData(userCtx2).Return(nil, database.ErrUserNotFound)

	type fields struct {
		UnimplementedNoteServicesServer pb.UnimplementedNoteServicesServer
		UnimplementedUserServicesServer pb.UnimplementedUserServicesServer
		db                              database.DataStorable
	}
	type args struct {
		ctx context.Context
		in1 *pb.NoteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.NoteList
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx1,
				in1: &pb.NoteRequest{},
			},
			want:    &pb.NoteList{Notes: []*pb.Note{&note1, &note2}},
			wantErr: false,
		},
		{
			name: "User not found",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx2,
				in1: &pb.NoteRequest{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong Ctx",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: context.Background(),
				in1: &pb.NoteRequest{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Controller{
				UnimplementedNoteServicesServer: tt.fields.UnimplementedNoteServicesServer,
				UnimplementedUserServicesServer: tt.fields.UnimplementedUserServicesServer,
				db:                              tt.fields.db,
			}
			got, err := s.GetNotes(tt.args.ctx, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_UpdateNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockDataStorable(ctrl)

	m.EXPECT().UpdateSecretData(userCtx1, secretData1upd).Return(&secretData1, nil)
	m.EXPECT().UpdateSecretData(userCtx2, secretData1upd).Return(nil, database.ErrUserNotFound)

	type fields struct {
		UnimplementedNoteServicesServer pb.UnimplementedNoteServicesServer
		UnimplementedUserServicesServer pb.UnimplementedUserServicesServer
		db                              database.DataStorable
	}
	type args struct {
		ctx  context.Context
		note *pb.Note
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx:  userCtx1,
				note: &note1,
			},
			want:    &empty.Empty{},
			wantErr: false,
		},
		{
			name: "User not found",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx:  userCtx2,
				note: &note1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong UUID",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx: userCtx2,
				note: &pb.Note{
					Id:         "UUID",
					Name:       "Test Note",
					Type:       models.CARD.String(),
					SecretData: []byte("Test Secret"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong Ctx",
			fields: fields{
				UnimplementedNoteServicesServer: pb.UnimplementedNoteServicesServer{},
				db:                              m,
			},
			args: args{
				ctx:  context.Background(),
				note: &note1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Controller{
				UnimplementedNoteServicesServer: tt.fields.UnimplementedNoteServicesServer,
				UnimplementedUserServicesServer: tt.fields.UnimplementedUserServicesServer,
				db:                              tt.fields.db,
			}
			got, err := s.UpdateNote(tt.args.ctx, tt.args.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateNote() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	md := mocks.NewMockDataStorable(ctrl)
	ms := mocks.NewMockServiceAuth(ctrl)
	as = ms

	testUserCrpt1 := testUser1
	testUserCrpt2 := testUser2
	testUserCrpt1.Password, _ = bcrypt.GenerateFromPassword(testUser1.Password, bcrypt.DefaultCost)
	testUserCrpt2.Password, _ = bcrypt.GenerateFromPassword(testUser2.Password, bcrypt.DefaultCost)

	md.EXPECT().GetUser(gomock.Any(), testUser1.Email).Return(&testUserCrpt1, nil)
	ms.EXPECT().CreateJwt(&testUserCrpt1).Return("test token", nil)

	md.EXPECT().GetUser(gomock.Any(), "userNotFound@test.com").Return(nil, gorm.ErrRecordNotFound)

	md.EXPECT().GetUser(gomock.Any(), testUser2.Email).Return(&testUserCrpt2, nil)
	ms.EXPECT().CreateJwt(&testUserCrpt2).Return("", errors.New("test error"))

	md.EXPECT().GetUser(gomock.Any(), "password@test.com").Return(&testUserCrpt2, nil)

	type fields struct {
		UnimplementedNoteServicesServer pb.UnimplementedNoteServicesServer
		UnimplementedUserServicesServer pb.UnimplementedUserServicesServer
		db                              database.DataStorable
	}
	type args struct {
		ctx  context.Context
		user *pb.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.JwtToken
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				UnimplementedUserServicesServer: pb.UnimplementedUserServicesServer{},
				db:                              md,
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Username: testUser1.Username,
					Password: string(testUser1.Password),
					Email:    testUser1.Email,
				},
			},
			want:    &pb.JwtToken{Token: "test token"},
			wantErr: false,
		},
		{
			name: "User not found",
			fields: fields{
				UnimplementedUserServicesServer: pb.UnimplementedUserServicesServer{},
				db:                              md,
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Username: "username",
					Password: "password",
					Email:    "userNotFound@test.com",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Can't create jwt",
			fields: fields{
				UnimplementedUserServicesServer: pb.UnimplementedUserServicesServer{},
				db:                              md,
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Username: testUser2.Username,
					Password: string(testUser2.Password),
					Email:    testUser2.Email,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong password",
			fields: fields{
				UnimplementedUserServicesServer: pb.UnimplementedUserServicesServer{},
				db:                              md,
			},
			args: args{
				ctx: context.Background(),
				user: &pb.User{
					Username: testUser2.Username,
					Password: "password",
					Email:    "password@test.com",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Controller{
				UnimplementedNoteServicesServer: tt.fields.UnimplementedNoteServicesServer,
				UnimplementedUserServicesServer: tt.fields.UnimplementedUserServicesServer,
				db:                              tt.fields.db,
			}
			got, err := s.Login(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	md := mocks.NewMockDataStorable(ctrl)
	ms := mocks.NewMockServiceAuth(ctrl)
	as = ms

	newUser := pb.User{Username: testUser1.Username, Password: string(testUser1.Password), Email: testUser1.Email}

	md.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(&testUser1, nil)
	ms.EXPECT().CreateJwt(&testUser1).Return("test token", nil)

	type fields struct {
		UnimplementedNoteServicesServer pb.UnimplementedNoteServicesServer
		UnimplementedUserServicesServer pb.UnimplementedUserServicesServer
		db                              database.DataStorable
	}
	type args struct {
		ctx  context.Context
		user *pb.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.JwtToken
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				UnimplementedUserServicesServer: pb.UnimplementedUserServicesServer{},
				db:                              md,
			},
			args: args{
				ctx:  context.Background(),
				user: &newUser,
			},
			want:    &pb.JwtToken{Token: "test token"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Controller{
				UnimplementedNoteServicesServer: tt.fields.UnimplementedNoteServicesServer,
				UnimplementedUserServicesServer: tt.fields.UnimplementedUserServicesServer,
				db:                              tt.fields.db,
			}
			got, err := s.Register(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	md := mocks.NewMockDataStorable(ctrl)
	ms := mocks.NewMockServiceAuth(ctrl)
	t.Run("New controller", func(t *testing.T) {
		got := NewController(log, md, ms)
		assert.NotNil(t, got)
	})
}

func TestTokenInterceptor(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     interface{}
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TokenInterceptor(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenInterceptor() got = %v, want %v", got, tt.want)
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
