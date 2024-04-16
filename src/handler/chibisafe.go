package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func Check_API_Key(Basepath string, accessKey string) bool {
	URL := Basepath + "/api/user/me"
	headers := map[string]string{
		"X-Api-Key": accessKey,
	}
	GETStruct := URLRequest{
		URL:    URL,
		Header: headers,
		Method: "GET",
	}
	resp, err := HTTPNoData(GETStruct)
	if err != nil {
		// Boolean is returned. So, no error was allowed to be returned
		log.Error().Msg(err.Error())
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func UploadPost(BasePath string, headers map[string]string, PostData UploadPostMeta) ([]byte, error) {
	URL := BasePath + "/api/upload"
	// Convert PostData to JSON
	PostDataJson, err := json.Marshal(PostData)
	if err != nil {
		return nil, err
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
		buffer := bytes.NewBuffer(nil)
		_, err := io.Copy(buffer, resp.Body)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
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
		buffer := bytes.NewBuffer(nil)
		_, err := io.Copy(buffer, resp.Body)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	} else {
		return nil, err
	}
}

func UploadProcessPost(BasePath string, headers map[string]string, PostData UploadProcessMeta) ([]byte, error) {
	URL := BasePath + "/api/upload/process"
	// Convert PostData to JSON
	PostDataJson, err := json.Marshal(PostData)
	if err != nil {
		return nil, err
	}

	POSTStruct := URLRequest{
		URL:    URL,
		Method: "POST",
		Header: headers,
	}

	resp, err := HTTPBytes(POSTStruct, bytes.NewBuffer(PostDataJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		buffer := bytes.NewBuffer(nil)
		_, err := io.Copy(buffer, resp.Body)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	} else {
		return nil, err
	}
}
