package main

import (
	"path/filepath"
	"testing"
)

func TestReadConfigNonExistent(t *testing.T) {
	config, err := readConfig(filepath.Join("test", "non-existent-file-name.config"))
	if err == nil || config != nil {
		t.Fatalf("Expected an error because of non-existent config file, but didn't get one.")
	}
}

func TestReadConfigInvalid(t *testing.T) {
	_, err := readConfig(filepath.Join("test", "invalid.conf"))
	if err == nil {
		t.Fatalf("Expected to err on invalid json format.")
	}
}

func TestReadConfig(t *testing.T) {
	config, err := readConfig(filepath.Join("test", "valid.conf"))
	if err != nil {
		t.Fatalf("Got unexpected error: %s\n", err.Error())
	}
	if config.Url != "http://localhost:8080/sabnzbd/api" {
		t.Errorf("Unexpected URL")
	}
	if config.Key != "YourSABnzbdApiOrNzbKey" {
		t.Errorf("Unexpected Key")
	}
}
