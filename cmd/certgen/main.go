package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		outPath string
		bits    int
	)

	flag.StringVar(&outPath, "out", "private.pem", "output path for private key")
	flag.IntVar(&bits, "bits", 2048, "RSA key size in bits")
	flag.Parse()

	if bits < 2048 {
		fmt.Fprintln(os.Stderr, "bits must be >= 2048")
		os.Exit(1)
	}

	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Fprintln(os.Stderr, "generate key:", err)
		os.Exit(1)
	}

	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err = os.WriteFile(outPath, pem.EncodeToMemory(pemBlock), 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "write key:", err)
		os.Exit(1)
	}

	fmt.Println("private key generated:", outPath)
}
