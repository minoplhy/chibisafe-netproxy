package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func UploadPost(BasePath string, PostData UploadPostMeta, accessKey string) ([]byte, error) {
	URL := BasePath + "/api/upload"
	// Convert PostData to JSON
	PostDataJson, err := json.Marshal(PostData)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	// Create a new request with POST method and request body
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(PostDataJson))
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", accessKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		BodyRead, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		return BodyRead, nil
	} else {
		log.Panicf("Output from %s : %d", URL, resp.StatusCode)
		return nil, nil
	}
}

// PUT to Network Storage(S3)
// URL argument is given from UploadPost Response!
func NetworkStoragePut(URL string, ContentType string, filepath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer file.Close()

	// Create an HTTP client
	client := &http.Client{}

	// Create a PUT request with the file contents
	req, err := http.NewRequest(http.MethodPut, URL, file)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer req.Body.Close()

	// Set appropriate headers for the file
	req.Header.Set("Content-Type", ContentType)

	// Set ContentLenght
	filestat, _ := file.Stat()
	req.ContentLength = filestat.Size()

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		BodyRead, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		return BodyRead, nil
	} else {
		log.Panicf("Output from %s : %d", URL, resp.StatusCode)
		return nil, err
	}
}

func UploadProcessPost(BasePath string, ContentType string, Identifier string, accessKey string, PostData UploadProcessMeta) ([]byte, error) {
	URL := BasePath + "/api/upload/process"
	// Convert PostData to JSON
	PostDataJson, err := json.Marshal(PostData)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	// Create a new request with POST method and request body
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(PostDataJson))
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", accessKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		BodyRead, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		return BodyRead, nil
	} else {
		log.Panicf("Output from %s : %d", URL, resp.StatusCode)
		return nil, nil
	}
}
