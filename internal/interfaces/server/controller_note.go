package server

import (
	"context"
	"errors"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/database"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/services/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	log  *logger.Logger
	once sync.Once
	cs   *Controller
	as   auth.ServiceAuth
)

type Controller struct {
	pb.UnimplementedNoteServicesServer
	pb.UnimplementedUserServicesServer
	db database.DataStorable
}

func NewController(logger *logger.Logger, db database.DataStorable, authService auth.ServiceAuth) *Controller {
	once.Do(func() {
		log = logger
		as = authService
		cs = &Controller{db: db}
	})
	return cs
}

func (s *Controller) AddNote(ctx context.Context, note *pb.Note) (*empty.Empty, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, status.Error(codes.Internal, "User not found")
	}
	log := log.WithFields(logrus.Fields{
		"method": "AddNote",
		"user":   userCtx.Email,
	})

	sd, err := interfaces.DtoToEntity(note)
	if err != nil {
		log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	sd.UserID = userCtx.Id

	_, err = s.db.AddSecretData(ctx, sd)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			log.Warn("Context not found")
			return nil, status.Error(codes.Unauthenticated, "User not authenticated")
		}
		log.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *Controller) DeleteNote(ctx context.Context, req *pb.NoteRequest) (*empty.Empty, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, status.Error(codes.Internal, "User not found")
	}
	log := log.WithFields(logrus.Fields{
		"method": "DeleteNote",
		"user":   userCtx.Email,
	})

	parse, err := uuid.Parse(req.IdNote)
	if err != nil {
		log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err = s.db.DeleteSecretData(ctx, parse)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			log.Warn("Context not found")
			return nil, status.Error(codes.Unauthenticated, "User not authenticated")
		}
		log.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !ok {
		log.Warn("User not found")
		return nil, status.Error(codes.NotFound, req.IdNote)
	}
	return &empty.Empty{}, nil
}

func (s *Controller) UpdateNote(ctx context.Context, note *pb.Note) (*empty.Empty, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, status.Error(codes.Internal, "User not found")
	}
	log := log.WithFields(logrus.Fields{
		"method": "UpdateNote",
		"user":   userCtx.Email,
	})

	sd, err := interfaces.DtoToEntity(note)
	if err != nil {
		log.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.db.UpdateSecretData(ctx, sd)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			log.Warn("Context not found")
			return nil, status.Error(codes.Unauthenticated, "User not authenticated")
		}
		log.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (s *Controller) GetNotes(ctx context.Context, _ *pb.NoteRequest) (*pb.NoteList, error) {
	userCtx, ok := ctx.Value("UserCtx").(*models.UserCtx)
	if !ok {
		return nil, status.Error(codes.Internal, "User not found")
	}
	log := log.WithFields(logrus.Fields{
		"method": "GetNotes",
		"user":   userCtx.Email,
	})

	sd, err := s.db.GetSecretData(ctx)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			log.Warn("Context not found")
			return nil, status.Error(codes.Unauthenticated, "User not authenticated")
		}
		log.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	notes := make([]*pb.Note, 0)
	log.Info("Create list of notes")
	for _, data := range *sd {
		notes = append(notes, &pb.Note{
			Id:         data.ID.String(),
			Name:       data.Name,
			Type:       data.Type,
			SecretData: data.Secret},
		)
	}
	return &pb.NoteList{Notes: notes}, nil
}

func TokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == pb.UserServices_Register_FullMethodName ||
		info.FullMethod == pb.UserServices_Login_FullMethodName {
		return handler(ctx, req)
	}
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
	}
	if len(token) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	userCtx, err := as.CreateUserCtx(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return handler(context.WithValue(ctx, "UserCtx", userCtx), req)
}
