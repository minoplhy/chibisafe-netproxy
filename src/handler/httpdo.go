package handler

import (
	"bytes"
	"net/http"
	"os"
)

func HTTPBytes(RequestStruct URLRequest, RequestData *bytes.Buffer) (*http.Response, error) {
	// Create a new request with POST method and request body
	req, err := http.NewRequest(RequestStruct.Method, RequestStruct.URL, RequestData)
	if err != nil {
		return nil, err
	}

	if RequestStruct.ContentLength != nil {
		req.ContentLength = int64(*RequestStruct.ContentLength)
	}

	for key, value := range RequestStruct.Header {
		req.Header.Set(key, value)
	}

	resp, err := HTTPClientDo(req)
	if err != nil {
		return nil, err
	}
	defer RequestData.Reset()
	return resp, nil
}

func HTTPOSFile(RequestStruct URLRequest, RequestData *os.File) (*http.Response, error) {
	// Create a new request with POST method and request body
	req, err := http.NewRequest(RequestStruct.Method, RequestStruct.URL, RequestData)
	if err != nil {
		return nil, err
	}

	if RequestStruct.ContentLength != nil {
		req.ContentLength = int64(*RequestStruct.ContentLength)
	}

	for key, value := range RequestStruct.Header {
		req.Header.Set(key, value)
	}

	resp, err := HTTPClientDo(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// HTTP With no Data sent
func HTTPNoData(RequestStruct URLRequest) (*http.Response, error) {
	req, err := http.NewRequest(RequestStruct.Method, RequestStruct.URL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range RequestStruct.Header {
		req.Header.Set(key, value)
	}

	resp, err := HTTPClientDo(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func HTTPClientDo(Request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(Request)
	if err != nil {
		return nil, err
	}
	return response, err
}
