package main

import "testing"

func TestInitLogger(t *testing.T) {
	logLevel = "info"

	initLogger()
	if appLogger == nil {
		t.Fatal("expected app logger to be initialized")
	}
}
