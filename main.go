package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const CONFIG_FILE = "$HOME/.config/enqueuenzb.conf"

func main() {
	config, err := readConfig(CONFIG_FILE)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	// Verify program arguments.
	if len(os.Args) < 2 {
		log.Fatalf("Please specify the nzb file to submit.\n")
	}

	for _, filePath := range os.Args[1:] {
		// Submit NZB file.
		if err = submitNzbFile(config, filePath); err != nil {
			log.Fatalf("Error during submission: %s\n", err.Error())
		}
		// Optionally remove NZB file.
		if config.Delete {
			if err = os.Remove(filePath); err != nil {
				log.Fatalf("Error while removing NZB file: %s\n", err.Error())
			}
		}
	}
}

func submitNzbFile(config *Config, filePath string) error {
	// Open nzb file
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Cannot open specified file: %s\n", err.Error())
	}
	defer fileHandle.Close()
	fileName := filepath.Base(filePath)
	// Generate message body
	boundary, content, err := createApiMessage(config, fileName, fileHandle)
	if err != nil {
		return fmt.Errorf("Error while creating API message: %s\n", err.Error())
	}
	// Create POST request
	request, err := http.NewRequest("POST", config.Url, content)
	if err != nil {
		return fmt.Errorf("Failed to create new request: %s\n", err.Error())
	}
	// Set the necessary header
	request.Header.Set("Content-Type", "multipart/form-data; charset=UTF-8; boundary="+boundary)
	// Send multipart message by POST request.
	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Error during POST: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Unexpected HTTP status returned: %d: %s\n", resp.StatusCode, resp.Status)
	}
	return nil
}

func createApiMessage(config *Config, name string, file io.Reader) (string, *bytes.Buffer, error) {
	if config == nil {
		return "", nil, fmt.Errorf("Cannot use nil config.")
	}
	var buffer bytes.Buffer
	// Create multipart message.
	multiPartBuffer := multipart.NewWriter(&buffer)
	multiPartBuffer.WriteField("output", "json")
	multiPartBuffer.WriteField("mode", "addfile")
	multiPartBuffer.WriteField("apikey", config.Key)
	multiPartBuffer.WriteField("nzbname", name)
	part, err := multiPartBuffer.CreateFormFile("nzbfile", name)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to create nzbfile entry: %s\n", err.Error())
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", nil, fmt.Errorf("Error while reading nzb file contents: %s\n", err.Error())
	}
	if err = multiPartBuffer.Close(); err != nil {
		return "", nil, fmt.Errorf("Error finishing up multipart message: %s\n", err.Error())
	}
	return multiPartBuffer.Boundary(), &buffer, nil
}
