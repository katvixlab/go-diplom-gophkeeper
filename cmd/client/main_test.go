package main

import (
	"path/filepath"
	"testing"
)

func TestInitLogger(t *testing.T) {
	logLevel = "info"
	logFile = filepath.Join(t.TempDir(), "client.log")

	initLogger()
	if appLogger == nil {
		t.Fatal("expected app logger to be initialized")
	}
}
