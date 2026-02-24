package main

import (
	"log"
	"os"

	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/mvc"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/services/ui"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var appLogger *logger.Logger

func main() {
	parseFlags()
	initLogger()

	conn, err := grpc.NewClient(connAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		if closeErr := conn.Close(); closeErr != nil {
			appLogger.Error(closeErr)
		}
	}(conn)

	uiService := ui.NewUIService(appLogger, conn)
	controller := mvc.NewUIController(appLogger, uiService)
	controller.AddItemInfoList("The application started successfully. Welcome!")
	if err = controller.Run(); err != nil {
		log.Fatal("failed to start UI controller", err)
	}
}

func initLogger() {
	logOutput := os.Stdout
	if file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666); err == nil {
		logOutput = file
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	rawLogger := &logrus.Logger{
		Out:   logOutput,
		Level: level,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	}

	appLogger = logger.NewLogger(rawLogger)
}
