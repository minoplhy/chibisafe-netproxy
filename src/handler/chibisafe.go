package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type UploadPostMeta struct {
	Name        string `json:"name"`
	FileSize    int64  `json:"size"`
	ContentType string `json:"contentType"`
}

type UploadResponseMeta struct {
	URL        string `json:"url"`
	Identifier string `json:"identifier"`
	PublicURL  string `json:"publicUrl"`
}

type UploadProcessMeta struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	ContentType string `json:"type"`
}

type UploadProcessResponseMeta struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}

func UploadPost(BasePath string, PostData UploadPostMeta, accessKey string) ([]byte, error) {
	URL := BasePath + "/api/upload"
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

	// Create a buffer to store the file contents
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create a PUT request with the file contents
	req, err := http.NewRequest("PUT", URL, &buffer)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer req.Body.Close()

	// Set appropriate headers for the file
	req.Header.Set("Content-Type", ContentType)

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
		return nil, nil
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
