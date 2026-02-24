package util

import (
	"context"
	"crypto/sha256"
	"hash"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
)

func TestEncrypt(t *testing.T) {
	type args struct {
		ctx   context.Context
		sha   hash.Hash
		data  []byte
		pass1 []byte
		pass2 []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				sha:   sha256.New(),
				data:  []byte("TEST TEST TEST TEST TEST TEST TEST TEST TEST TEST"),
				pass1: []byte("123456"),
				pass2: []byte("123456"),
			},
			want:    []byte("TEST TEST TEST TEST TEST TEST TEST TEST TEST TEST"),
			wantErr: false,
		},
		{
			name: "success 100kb",
			args: args{
				ctx:   context.Background(),
				sha:   sha256.New(),
				data:  []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				pass1: []byte("123456"),
				pass2: []byte("123456"),
			},
			want:    []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
			wantErr: false,
		},
		{
			name: "wrong password",
			args: args{
				ctx:   context.Background(),
				sha:   sha256.New(),
				data:  []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit"),
				pass1: []byte("Test123"),
				pass2: []byte("Test124"),
			},
			want:    []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.sha.Write(tt.args.pass1)
			pass1 := tt.args.sha.Sum(nil)

			encrypt, err := Encrypt(context.Background(), pass1, tt.args.data)
			if err != nil {
				t.Fatal(err)
			}
			tt.args.sha.Reset()
			tt.args.sha.Write(tt.args.pass2)
			pass2 := tt.args.sha.Sum(nil)
			decrypt, err := Decrypt(context.Background(), pass2, encrypt)

			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(decrypt, tt.want) {
				t.Errorf("Encrypt() got = %v, want %v", decrypt, tt.want)
			}
		})
	}
}

func Test_generateRandom(t *testing.T) {

	type args struct {
		size int
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				size: 10,
			},
			wantSize: 10,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := generateRandom(tt.args.size)
			got2, err := generateRandom(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateRandom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got1, got2) {
				t.Errorf("generateRandom() hash can't be the same. got1 = %v, got2 = %v", got1, got2)
			}
		})
	}
}

func TestAddContextUserCtx(t *testing.T) {

	t.Run("UserCtx", func(t *testing.T) {
		got := AddContextUserCtx(context.Background(), "username", "email", uuid.Nil)
		userCtx := got.Value("UserCtx").(*models.UserCtx)
		userWant := &models.UserCtx{
			Username: "username",
			Email:    "email",
			Id:       uuid.Nil,
		}
		if !reflect.DeepEqual(userCtx, userWant) {
			t.Errorf("AddContextUserCtx() = %v, want %v", userCtx, userWant)
		}
	})

}
