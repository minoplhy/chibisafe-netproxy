package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/minoplhy/chibisafe_netstorage_middleman/src/handler"
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

	// Open or create a file for appending logs
	log_file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	if r.Method != "POST" {
		http.Error(w, handler.ErrorResponseBuild(http.StatusMethodNotAllowed, "Method not allowed"), http.StatusMethodNotAllowed)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr}, "Method not allowed")
		return
	}

	// Check for x-api-key header
	API_Key := r.Header.Get("x-api-key")
	if API_Key == "" {
		http.Error(w, handler.ErrorResponseBuild(http.StatusBadRequest, "X-api-key is empty!"), http.StatusBadRequest)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr}, "X-api-key is empty!")
		return
	}

	// Validate x-api-key
	if !handler.Check_API_Key(Chibisafe_basepath, API_Key) {
		http.Error(w, handler.ErrorResponseBuild(http.StatusUnauthorized, "Failure to validate X-API-Key"), http.StatusUnauthorized)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr}, "Failure to validate X-API-Key")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize))

	// ParseMultipartForm parses a request body as multipart/form-data
	err = r.ParseMultipartForm(int64(maxUploadSize))
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, handler.ErrorResponseBuild(http.StatusRequestEntityTooLarge, "Request Body is too large!"), http.StatusRequestEntityTooLarge)
			handler.LogBuilder("ERROR", []string{r.RemoteAddr}, "Request Body is too large!")
			return
		}
		http.Error(w, handler.ErrorResponseBuild(http.StatusInternalServerError, "Something went wrong!"), http.StatusInternalServerError)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr}, err.Error())
		return
	}

	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, handler.ErrorResponseBuild(http.StatusInternalServerError, "Something went wrong!"), http.StatusInternalServerError)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr}, err.Error())
		return
	}
	defer file.Close()

	handler.LogBuilder("INFO", []string{r.RemoteAddr}, "Received a successful POST")
	tempfilepath := handler.GetTempFilename(fileHeader.Filename)
	handler.LogBuilder("INFO", []string{r.RemoteAddr, tempfilepath}, "Successfully obtained temporary Filename")
	handler.SaveFile(tempfilepath, file)
	handler.DiscardFile(file)

	PostData := handler.UploadPostMeta{
		ContentType: fileHeader.Header.Get("Content-Type"),
		Name:        fileHeader.Filename,
		FileSize:    fileHeader.Size,
	}

	chibisafe_post, err := handler.UploadPost(Chibisafe_basepath, API_Key, PostData)
	if err != nil {
		http.Error(w, handler.ErrorResponseBuild(http.StatusBadRequest, "Something went wrong!"), http.StatusBadRequest)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr, tempfilepath}, err.Error())
		return
	}

	var chibisafe_Response_Metadata handler.UploadResponseMeta
	err = json.Unmarshal(chibisafe_post, &chibisafe_Response_Metadata)
	if err != nil {
		http.Error(w, handler.ErrorResponseBuild(http.StatusInternalServerError, "Something went wrong!"), http.StatusInternalServerError)
		handler.LogBuilder("ERROR", []string{}, err.Error())
		return
	}
	handler.LogBuilder("INFO", []string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, "Successfully obtained PUT keys")

	_, err = handler.NetworkStoragePut(chibisafe_Response_Metadata.URL, PostData.ContentType, tempfilepath)
	if err != nil {
		http.Error(w, handler.ErrorResponseBuild(http.StatusInternalServerError, "Something went wrong!"), http.StatusInternalServerError)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, err.Error())
		return
	}
	handler.LogBuilder("INFO", []string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, "Successfully PUT file to Network Storage")

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

	PostProcess, err := handler.UploadProcessPost(Chibisafe_basepath, API_Key, PostProcessData)
	if err != nil {
		http.Error(w, handler.ErrorResponseBuild(http.StatusInternalServerError, "Something went wrong!"), http.StatusInternalServerError)
		handler.LogBuilder("ERROR", []string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, err.Error())
		return
	}

	var PostProcessResponse handler.UploadProcessResponseMeta
	err = json.Unmarshal(PostProcess, &PostProcessResponse)
	if err != nil {
		http.Error(w, handler.ErrorResponseBuild(http.StatusInternalServerError, "Something went wrong!"), http.StatusInternalServerError)
		handler.LogBuilder("ERROR", []string{}, err.Error())

		return
	}
	handler.LogBuilder("INFO", []string{r.RemoteAddr, PostProcessResponse.Name, tempfilepath}, fmt.Sprintf("Successfully Processed Response with UUID: %s", PostProcessResponse.UUID))

	err = handler.DeleteFile(tempfilepath)
	if err != nil {
		handler.LogBuilder("ERROR", []string{r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath}, err.Error())
		return
	}
	handler.LogBuilder("INFO", []string{r.RemoteAddr, tempfilepath}, "Successfully Deleted Temporary file from local disk")
	JsonResponse, _ := json.Marshal(PostProcessResponse)
	fmt.Fprintf(w, "%s", JsonResponse)
}

func main() {
	Chibisafe_basepath := os.Getenv("CHIBISAFE_BASEPATH")
	Max_Upload_Size := os.Getenv("MAX_UPLOAD_SIZE")

	if Chibisafe_basepath == "" {
		log.Fatal("CHIBISAFE_BASEPATH environment is not set!")
	}
	if Max_Upload_Size != "" {
		_, err := strconv.Atoi(Max_Upload_Size)
		if err != nil {
			log.Fatal("MAX_UPLOAD_SIZE environment is invaild!")
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/upload", uploadHandler)

	if err := http.ListenAndServe(":4040", mux); err != nil {
		log.Fatal(err)
	}
}
