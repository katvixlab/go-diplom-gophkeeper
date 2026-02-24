package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	var (
		addr     string
		username string
		email    string
		password string
	)

	flag.StringVar(&addr, "addr", "localhost:3200", "gRPC server address")
	flag.StringVar(&username, "username", "demo-user", "demo username")
	flag.StringVar(&email, "email", "demo@example.com", "demo email")
	flag.StringVar(&password, "password", "DemoPass123!", "demo password")
	flag.Parse()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	userClient := pb.NewUserServicesClient(conn)
	noteClient := pb.NewNoteServicesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jwtToken, err := userClient.Register(ctx, &pb.User{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("register failed: %v", err)
	}

	notes := []models.Noteable{
		models.CredentialNote{
			Username: "demo-login",
			Password: "s3cr3t-password",
			BaseNote: models.BaseNote{
				Id:         uuid.New(),
				NameRecord: "GitHub account",
				Created:    time.Now().Unix(),
				Type:       models.CREDENTIAL,
				MetaInfo:   []string{"portfolio", "demo"},
			},
		},
		models.TextNote{
			Text: "Demo recovery code: 1234-5678-ABCD",
			BaseNote: models.BaseNote{
				Id:         uuid.New(),
				NameRecord: "Recovery code",
				Created:    time.Now().Unix(),
				Type:       models.TEXT,
				MetaInfo:   []string{"temporary", "sample"},
			},
		},
		models.BankCardNote{
			Bank:         "Demo Bank",
			Number:       "4111 1111 1111 1111",
			Expiration:   "12/30",
			Cardholder:   "DEMO USER",
			SecurityCode: "123",
			BaseNote: models.BaseNote{
				Id:         uuid.New(),
				NameRecord: "Demo card",
				Created:    time.Now().Unix(),
				Type:       models.CARD,
				MetaInfo:   []string{"test only", "not real"},
			},
		},
	}

	key := getHash(email, password)
	for _, note := range notes {
		if err = addNote(context.Background(), noteClient, jwtToken.Token, key, note); err != nil {
			log.Fatalf("seed note failed: %v", err)
		}
	}

	fmt.Printf("seed completed for %s (%s): %d notes\n", username, email, len(notes))
}

func addNote(ctx context.Context, noteClient pb.NoteServicesClient, jwtToken string, key []byte, note models.Noteable) error {
	payload, err := json.Marshal(note)
	if err != nil {
		return err
	}

	encrypted, err := util.Encrypt(ctx, key, payload)
	if err != nil {
		return err
	}

	req := &pb.Note{
		Id:         note.GetID().String(),
		Name:       note.GetName(),
		Type:       note.GetType().String(),
		SecretData: encrypted,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"token": jwtToken}))
	_, err = noteClient.AddNote(ctx, req)
	return err
}

func getHash(email, password string) []byte {
	sum := sha256.Sum256([]byte(email + password))
	return sum[:]
}
