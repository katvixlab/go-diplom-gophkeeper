package interfaces

import (
	"github.com/google/uuid"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
)

func DtoToEntity(note *pb.Note) (models.SecretData, error) {
	uid, err := uuid.Parse(note.Id)
	if err != nil {
		return models.SecretData{}, err
	}
	return models.SecretData{
		ID:     uid,
		Type:   note.Type,
		Name:   note.Name,
		Secret: note.SecretData,
	}, nil
}
