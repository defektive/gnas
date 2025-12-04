package files

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetURLReader(parsedFile *url.URL) (io.ReadCloser, error) {

	if parsedFile.Scheme == "http" || parsedFile.Scheme == "https" {
		req, err := http.NewRequest("GET", parsedFile.String(), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Mobile/15E148 Safari/604.1")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		return resp.Body, err
	}
	return nil, errors.New("unsupported protocol scheme " + parsedFile.Scheme)

}

func GetFileReader(fileToGet string) (io.ReadCloser, error) {

	parsedFile, err := url.Parse(fileToGet)
	if err != nil {
		return nil, err
	}

	return GetURLReader(parsedFile)
}

func DownloadURL(parsedFile *url.URL, savePath string) error {

	fileReader, err := GetURLReader(parsedFile)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	// Create the file to save the content
	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close() // Ensure the file is closed

	// Write the response body to the file
	_, err = io.Copy(out, fileReader)
	if err != nil {
		return fmt.Errorf("failed to copy data to file: %w", err)
	}

	return nil
}

func DownloadFile(fileURL string, savePath string) error {
	parsedFile, err := url.Parse(fileURL)
	if err != nil {
		return err
	}

	return DownloadURL(parsedFile, savePath)
}

func Download(fileURL string) error {
	parsedFile, err := url.Parse(fileURL)
	if err != nil {
		return err
	}

	savePath := filepath.Base(parsedFile.Path)

	return DownloadURL(parsedFile, savePath)
}

func UploadFile(remoteServer, authToken, filePath string) error {

	if !Exists(filePath) {
		return errors.New("file does not exist")
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, remoteServer, file)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Content-Type header (optional, but good practice)
	// You might need a more specific Content-Type based on your file type
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.New("error uploading file: " + resp.Status)
	}

	return nil
}
