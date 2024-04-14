package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func ResponseBuild(w http.ResponseWriter, ContentType string, ContentData []byte) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(ContentData)
}

func ErrorResponseBuild(w http.ResponseWriter, StatusCode int64, Message string) {
	ErrorResponse := ErrorResponse{
		StatusCode: StatusCode,
		Message:    Message,
	}
	Response, _ := json.Marshal(ErrorResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(StatusCode))
	w.Write(Response)
}

func InfoLogBuilder(headers []string, message string) {
	logString := "[" + strings.Join(headers, "] [") + "]" + " : " + message
	log.Info().Msgf(logString)
}

func ErrorLogBuilder(headers []string, message string) {
	logString := "[" + strings.Join(headers, "] [") + "]" + " : " + message
	log.Error().Msgf(logString)
}
