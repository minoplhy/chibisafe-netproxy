package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func UploadPost(BasePath string, PostData UploadPostMeta, accessKey string) ([]byte, error) {
	URL := BasePath + "/api/upload"
	// Convert PostData to JSON
	PostDataJson, err := json.Marshal(PostData)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-Api-Key":    accessKey,
		"Content-Type": "application/json",
	}

	POSTStruct := URLRequest{
		URL:         URL,
		ContentType: PostData.ContentType,
		Method:      "POST",
		Header:      headers,
	}

	resp, err := HTTPBytes(POSTStruct, bytes.NewBuffer(PostDataJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		BodyRead, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return BodyRead, nil
	} else {
		return nil, err
	}
}

// PUT to Network Storage(S3)
// URL argument is given from UploadPost Response!
func NetworkStoragePut(URL string, ContentType string, filepath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Set ContentLenght
	filestat, _ := file.Stat()

	var filesize *int = new(int)
	*filesize = int(filestat.Size())

	headers := map[string]string{
		"Content-Type": ContentType,
	}

	PUTStruct := URLRequest{
		URL:           URL,
		ContentType:   ContentType,
		Method:        "PUT",
		Header:        headers,
		ContentLength: filesize,
	}

	// Send the request
	resp, err := HTTPOSFile(PUTStruct, file)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		BodyRead, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return BodyRead, nil
	} else {
		return nil, err
	}
}

func UploadProcessPost(BasePath string, ContentType string, Identifier string, accessKey string, PostData UploadProcessMeta) ([]byte, error) {
	URL := BasePath + "/api/upload/process"
	// Convert PostData to JSON
	PostDataJson, err := json.Marshal(PostData)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-Api-Key":    accessKey,
		"Content-Type": "application/json",
	}

	POSTStruct := URLRequest{
		URL:         URL,
		ContentType: ContentType,
		Method:      "POST",
		Header:      headers,
	}

	resp, err := HTTPBytes(POSTStruct, bytes.NewBuffer(PostDataJson))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		BodyRead, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return BodyRead, nil
	} else {
		return nil, err
	}
}
