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

func LogBuilder(Level string, headers []string, message string) {
	logString := "[" + Level + "]" + " "
	logString += "[" + strings.Join(headers, "] [") + "]" + " : " + message
	log.Println(logString)
}
