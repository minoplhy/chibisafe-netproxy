package handler

import "encoding/json"

func ErrorResponseBuild(StatusCode int64, Message string) string {
	ErrorResponse := ErrorResponse{
		StatusCode: StatusCode,
		Message:    Message,
	}
	Response, _ := json.Marshal(ErrorResponse)
	return string(Response)
}
