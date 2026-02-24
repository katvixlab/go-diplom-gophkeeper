package main

import (
	"log"
	"net"
	"os"

	"github.com/katvixlab/go-diplom-gophkeeper/internal/database"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/server"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/services/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var appLogger *logger.Logger

func main() {
	parseFlags()
	initLogger()

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	store := *database.NewDataStore(appLogger, db)
	if err = store.Migrate(); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	authService, err := auth.NewAuthService(appLogger, crtFile)
	if err != nil {
		log.Fatal("failed to initialize auth service", err)
	}

	controller := server.NewController(appLogger, store, *authService)
	listener, err := net.Listen("tcp", srvAddr)
	if err != nil {
		log.Fatal("failed to start listener", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.TokenInterceptor))
	pb.RegisterNoteServicesServer(grpcServer, controller)
	pb.RegisterUserServicesServer(grpcServer, controller)

	appLogger.Info("server started")
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal("failed to start gRPC server", err)
	}
}

func initLogger() {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	rawLogger := &logrus.Logger{
		Out:   os.Stdout,
		Level: level,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "@level",
				logrus.FieldKeyMsg:   "@message",
				logrus.FieldKeyFunc:  "@caller",
			},
		},
	}

	appLogger = logger.NewLogger(rawLogger)
}
