package interfaces

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
)

func TestDtoToEntity(t *testing.T) {
	type args struct {
		note *pb.Note
	}
	tests := []struct {
		name    string
		args    args
		want    models.SecretData
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				note: &pb.Note{
					Id:         uuid.Nil.String(),
					Name:       "Test Note",
					Type:       models.CARD.String(),
					SecretData: []byte{1, 2, 3},
				},
			},
			want: models.SecretData{
				ID:     uuid.Nil,
				Name:   "Test Note",
				Type:   models.CARD.String(),
				Secret: []byte{1, 2, 3},
			},
			wantErr: false,
		},
		{
			name: "Wrong Id",
			args: args{
				note: &pb.Note{
					Id:         "uuid",
					Name:       "Test Note",
					Type:       models.CARD.String(),
					SecretData: []byte{1, 2, 3},
				},
			},
			want:    models.SecretData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DtoToEntity(tt.args.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("DtoToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DtoToEntity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
