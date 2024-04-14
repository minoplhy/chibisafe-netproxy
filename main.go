package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/minoplhy/chibisafe-netproxy/src/handler"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check is already done in main()
	Chibisafe_basepath := os.Getenv("CHIBISAFE_BASEPATH")
	// Max Upload Size
	// if Enviroment is failed, will fallback to 10MB
	GetMaxUploadSize := os.Getenv("MAX_UPLOAD_SIZE")
	maxUploadSize, err := strconv.Atoi(GetMaxUploadSize)
	if err != nil {
		maxUploadSize = 10 * 1024 * 1024 // 10 MB
	}

	if r.Method != "POST" {
		handler.ErrorResponseBuild(w, http.StatusMethodNotAllowed, "Method not allowed")
		handler.ErrorLogBuilder([]string{r.RemoteAddr}, "Method not allowed")
		return
	}

	// Check for x-api-key header
	API_Key := r.Header.Get("x-api-key")
	if API_Key == "" {
		handler.ErrorResponseBuild(w, http.StatusBadRequest, "X-api-key is empty!")
		handler.ErrorLogBuilder([]string{r.RemoteAddr}, "X-api-key is empty!")
		return
	}

	// Validate x-api-key
	if !handler.Check_API_Key(Chibisafe_basepath, API_Key) {
		handler.ErrorResponseBuild(w, http.StatusUnauthorized, "Failure to validate X-API-Key")
		handler.ErrorLogBuilder([]string{r.RemoteAddr}, "Failure to validate X-API-Key")
		return
	}

	// Set Max Body read size
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize))

	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		switch err.Error() {
		case "http: request body too large":
			handler.ErrorResponseBuild(w, http.StatusRequestEntityTooLarge, "Request Body is too large!")
			handler.ErrorLogBuilder([]string{r.RemoteAddr}, "Request Body is too large!")
		default:
			handler.ErrorResponseBuild(w, http.StatusInternalServerError, "Something went wrong!")
			handler.ErrorLogBuilder([]string{r.RemoteAddr}, err.Error())
		}
		return
	}
	defer file.Close()

	handler.InfoLogBuilder([]string{r.RemoteAddr}, "Received a successful POST")
	tempfilepath := handler.GetTempFilename(fileHeader.Filename)
	handler.InfoLogBuilder([]string{r.RemoteAddr, tempfilepath}, "Successfully obtained temporary Filename")

	// Save Temporary file
	handler.SaveFile(tempfilepath, file)

	PostData := handler.UploadPostMeta{
		ContentType: fileHeader.Header.Get("Content-Type"),
		Name:        fileHeader.Filename,
		FileSize:    fileHeader.Size,
	}

	UploadHeaders := map[string]string{
		"X-Api-Key":    API_Key,
		"Content-Type": "application/json",
	}

	// Check if client sent X-Real-IP Header
	if r.Header.Get("X-Real-IP") != "" && handler.IsInternalIP(r.RemoteAddr) {
		UploadHeaders["X-Real-IP"] = r.Header.Get("X-Real-IP")
	}

	// POST to chibisafe's /api/upload
	chibisafe_post, err := handler.UploadPost(Chibisafe_basepath, UploadHeaders, PostData)
	if err != nil {
		switch err.Error() {
		default:
			handler.ErrorResponseBuild(w, http.StatusBadRequest, "Something went wrong!")
			handler.ErrorLogBuilder([]string{r.RemoteAddr, tempfilepath}, err.Error())
		}
		return
	}

	var chibisafe_Response_Metadata handler.UploadResponseMeta
	err = json.Unmarshal(chibisafe_post, &chibisafe_Response_Metadata)
	if err != nil {
		switch err.Error() {
		default:
			handler.ErrorResponseBuild(w, http.StatusInternalServerError, "Something went wrong!")
			handler.ErrorLogBuilder([]string{}, err.Error())
		}
		return
	}
	handler.InfoLogBuilder([]string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, "Successfully obtained PUT keys")

	// PUT File to Network Storage
	_, err = handler.NetworkStoragePut(chibisafe_Response_Metadata.URL, PostData.ContentType, tempfilepath)
	if err != nil {
		switch err.Error() {
		default:
			handler.ErrorResponseBuild(w, http.StatusInternalServerError, "Something went wrong!")
			handler.ErrorLogBuilder([]string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, err.Error())
		}
		return
	}
	handler.InfoLogBuilder([]string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, "Successfully PUT file to Network Storage")

	// Build Struct for PostProcess Json
	//
	// Name 	   -> original Filename
	// ContentType -> original Content-Type
	// Identifier  -> File Identifier ID
	PostProcessData := handler.UploadProcessMeta{
		Name:        fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Identifier:  chibisafe_Response_Metadata.Identifier,
	}

	ProcessHeaders := map[string]string{
		"X-Api-Key":    API_Key,
		"Content-Type": "application/json",
	}

	// Check if client sent X-Real-IP Header
	if r.Header.Get("X-Real-IP") != "" && handler.IsInternalIP(r.RemoteAddr) {
		ProcessHeaders["X-Real-IP"] = r.Header.Get("X-Real-IP")
	}

	// POST to chibisafe's /api/upload/process
	PostProcess, err := handler.UploadProcessPost(Chibisafe_basepath, ProcessHeaders, PostProcessData)
	if err != nil {
		switch err.Error() {
		default:
			handler.ErrorResponseBuild(w, http.StatusInternalServerError, "Something went wrong!")
			handler.ErrorLogBuilder([]string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, err.Error())
		}
		return
	}

	var PostProcessResponse handler.UploadProcessResponseMeta
	err = json.Unmarshal(PostProcess, &PostProcessResponse)
	if err != nil {
		switch err.Error() {
		default:
			handler.ErrorResponseBuild(w, http.StatusInternalServerError, "Something went wrong!")
			handler.ErrorLogBuilder([]string{}, err.Error())
		}
		return
	}
	handler.InfoLogBuilder([]string{r.RemoteAddr, PostProcessResponse.Name, tempfilepath}, fmt.Sprintf("Successfully Processed Response with UUID: %s", PostProcessResponse.UUID))

	// Delete Temporary file
	err = handler.DeleteFile(tempfilepath)
	if err != nil {
		switch err.Error() {
		default:
			handler.ErrorLogBuilder([]string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, err.Error())
		}
		return
	}
	handler.InfoLogBuilder([]string{r.RemoteAddr, tempfilepath}, "Successfully Deleted Temporary file from local disk")

	// Response JSON to client
	JsonResponse, _ := json.Marshal(PostProcessResponse)
	handler.ResponseBuild(w, "application/json", JsonResponse)
}

func main() {
	// Open or create a file for appending logs
	log_file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer log_file.Close()

	// Setup Logging Policy
	// Multi level writer on logfile and console
	//
	// Format : Console -> Human Readable
	// 			File 	-> Json
	logger := zerolog.New(zerolog.MultiLevelWriter(log_file, os.Stdout)).With().Timestamp().Logger()
	logger = logger.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}, log_file))
	// Set as Global logger :)
	log.Logger = logger

	Chibisafe_basepath := os.Getenv("CHIBISAFE_BASEPATH")
	Max_Upload_Size := os.Getenv("MAX_UPLOAD_SIZE")

	if Chibisafe_basepath == "" {
		log.Fatal().Msg("CHIBISAFE_BASEPATH environment is not set!")
	}
	if Max_Upload_Size != "" {
		_, err := strconv.Atoi(Max_Upload_Size)
		if err != nil {
			log.Fatal().Msg("MAX_UPLOAD_SIZE environment is invaild!")
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/upload", uploadHandler)

	if err := http.ListenAndServe(":4040", mux); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
