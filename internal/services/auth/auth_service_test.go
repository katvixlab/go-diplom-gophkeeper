package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	uidU1     = uuid.New()
	testUser1 = models.User{ID: uidU1, Username: "Test User", Password: []byte([]byte("Test Password")), Email: "user1@test.com"}
)

func TestMain(m *testing.M) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.Fatal(err)
	}
	pemBytes, err := creteCert(privateKey)
	blocks, _ := pem.Decode(pemBytes)
	key, err := x509.ParsePKCS1PrivateKey(blocks.Bytes)

	as = &Service{
		privateKey: key,
		publicKey:  &privateKey.PublicKey,
	}

	log = logger.NewLogger(logrus.New())
	code := m.Run()
	os.Exit(code)
}

func TestCreateJwt(t *testing.T) {
	gotToken, err := as.CreateJwt(&testUser1)
	if err != nil {
		t.Errorf("CreateJwt() error = %v", err)
	}

	gotClaims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(gotToken, gotClaims, func(token *jwt.Token) (interface{}, error) {
		return as.publicKey, nil
	})
	if err != nil {
		t.Errorf("CreateJwt() error = %v", err)
	}

	wantClaims := jwt.MapClaims{
		"Id":       testUser1.ID.String(),
		"Username": testUser1.Username,
		"Email":    testUser1.Email,
		"exp":      token.Claims.(jwt.MapClaims)["exp"].(float64),
	}
	assert.Equal(t, token.Claims, wantClaims,
		"CreateJwt() gotToken = %v, want %v", token.Claims, wantClaims)

	wantUserCtx := &models.UserCtx{
		Id:       testUser1.ID,
		Username: testUser1.Username,
		Email:    testUser1.Email,
	}
	gotUserCtx, err := as.CreateUserCtx(gotToken)
	if err != nil {
		t.Errorf("CreateUserCtx() error = %v", err)
	}
	assert.Equal(t, gotUserCtx, wantUserCtx,
		"CreateUserCtx() gotUserCtx = %v, want %v", gotUserCtx, wantUserCtx)
}

func Test_getFile(t *testing.T) {
	testFile := "test.txt"
	wantBytes := []byte("test")

	file, err := os.OpenFile(testFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Errorf("getFile() error = %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		err = os.Remove(file.Name())
		if err != nil {
			t.Errorf("getFile() error = %v", err)
		}
	}(file)

	_, err = file.Write(wantBytes)
	if err != nil {
		t.Errorf("getFile() error = %v", err)
	}

	gotBytes, err := getFile(testFile)
	assert.Equalf(t, wantBytes, gotBytes, "getFile(%v)", testFile)

}

func TestNewAuthService(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("NewAuthService() error = %v", err)
	}

	testFile := "test.pem"
	pemBytes, err := creteCert(privateKey)

	file, err := os.OpenFile(testFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Errorf("getFile() error = %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		err = os.Remove(file.Name())
		if err != nil {
			t.Errorf("getFile() error = %v", err)
		}
	}(file)

	_, err = file.Write(pemBytes)
	if err != nil {
		t.Errorf("getFile() error = %v", err)
	}

	_, err = NewAuthService(logger.NewLogger(logrus.New()), testFile)

	assert.NoError(t, err, "NewAuthService() error = %v", err)
}

func creteCert(privateKey *rsa.PrivateKey) ([]byte, error) {

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Example Domain",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage: x509.KeyUsageDigitalSignature,
	}

	_, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	return pemBytes, err
}
