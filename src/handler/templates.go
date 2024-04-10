package handler

import (
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
)

func ErrorResponseBuild(StatusCode int64, Message string) string {
	ErrorResponse := ErrorResponse{
		StatusCode: StatusCode,
		Message:    Message,
	}
	Response, _ := json.Marshal(ErrorResponse)
	return string(Response)
}

func InfoLogBuilder(headers []string, message string) {
	logString := "[" + strings.Join(headers, "] [") + "]" + " : " + message
	log.Info().Msgf(logString)
}

func ErrorLogBuilder(headers []string, message string) {
	logString := "[" + strings.Join(headers, "] [") + "]" + " : " + message
	log.Error().Msgf(logString)
}
