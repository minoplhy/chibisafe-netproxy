package handler

import (
	"encoding/json"
	"log"
	"strings"
)

func ErrorResponseBuild(StatusCode int64, Message string) string {
	ErrorResponse := ErrorResponse{
		StatusCode: StatusCode,
		Message:    Message,
	}
	Response, _ := json.Marshal(ErrorResponse)
	return string(Response)
}

func LogBuilder(Type string, headers []string, message string) {
	logString := "[" + Type + "]" + " "
	logString += "[" + strings.Join(headers, "] [") + "]" + " : " + message
	log.Println(logString)
}
