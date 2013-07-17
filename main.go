package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const CONFIG_FILE = "$HOME/.config/queuenzb.conf"

func main() {
	config := readConfig()

	// Verify program arguments.
	if len(os.Args) < 2 {
		log.Fatalf("Please specify the nzb file to submit.\n")
	}

	// Open nzb file
	fileHandle, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Cannot open specified file: %s\n", err.Error())
	}
	defer fileHandle.Close()
	fileName := filepath.Base(os.Args[1])

	// Create multipart message.
	var buffer bytes.Buffer
	multiPartBuffer := multipart.NewWriter(&buffer)
	multiPartBuffer.WriteField("output", "json")
	multiPartBuffer.WriteField("mode", "addfile")
	multiPartBuffer.WriteField("apikey", config.Key)
	multiPartBuffer.WriteField("nzbname", fileName)
	part, err := multiPartBuffer.CreateFormFile("nzbfile", fileName)
	if err != nil {
		log.Fatalf("Failed to create nzbfile entry: %s\n", err.Error())
	}
	io.Copy(part, fileHandle)
	if err = multiPartBuffer.Close(); err != nil {
		log.Fatalf("Error finishing up multipart message: %s\n", err.Error())
	}

	// Create POST request
	request, err := http.NewRequest("POST", config.Url, &buffer)
	if err != nil {
		log.Fatalf("Failed to create new request: %s\n", err.Error())
	}
	request.Header.Set("Content-Type", "multipart/form-data; charset=UTF-8; boundary="+multiPartBuffer.Boundary())

	// Send multipart message by POST request.
	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error during POST: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Unexpected HTTP status returned: %d: %s\n", resp.StatusCode, resp.Status)
	}
}

type Config struct {
	Url string
	Key string
}

func readConfig() *Config {
	path := os.ExpandEnv(CONFIG_FILE)
	// Read config file
	configFile, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		log.Fatalf("Cannot read config file '%s': %s\n", path, err.Error())
	}
	defer configFile.Close()
	// Decode json content
	config := new(Config)
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode json config file '%s': %s\n", path, err.Error())
	}
	return config
}
