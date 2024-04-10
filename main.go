package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/minoplhy/chibisafe_netstorage_middleman/src/handler"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check is already done in main()
	Chibisafe_basepath := os.Getenv("CHIBISAFE_BASEPATH")

	// Open or create a file for appending logs
	log_file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("[%s] : Method not allowed", r.RemoteAddr)
		return
	}

	// Check for x-api-key header
	API_Key := r.Header.Get("x-api-key")
	if API_Key == "" {
		http.Error(w, "X-api-key is empty!", http.StatusBadRequest)
		log.Printf("[%s] : X-api-key is empty!", r.RemoteAddr)
		return
	}

	// Set a limit on the request body size
	maxUploadSize := 10 * 1024 * 1024 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize))

	// ParseMultipartForm parses a request body as multipart/form-data
	err = r.ParseMultipartForm(int64(maxUploadSize))
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "File too large.", http.StatusRequestEntityTooLarge)
			log.Printf("[%s] : Request Body is too large!", r.RemoteAddr)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("[%s] : %s", r.RemoteAddr, err)
		return
	}
	defer file.Close()

	log.Printf("[%s] : Received a successful POST", r.RemoteAddr)
	tempfilepath := handler.GetTempFilename(fileHeader.Filename)
	log.Printf("[%s] [%s] : Successfully obtained temporary Filename", r.RemoteAddr, tempfilepath)
	handler.SaveFile(tempfilepath, file)
	handler.DiscardFile(file)

	PostData := handler.UploadPostMeta{
		ContentType: fileHeader.Header.Get("Content-Type"),
		Name:        fileHeader.Filename,
		FileSize:    fileHeader.Size,
	}

	chibisafe_post, err := handler.UploadPost(Chibisafe_basepath, PostData, API_Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("[%s] [%s] : %s", r.RemoteAddr, tempfilepath, err)
		return
	}

	var chibisafe_Response_Metadata handler.UploadResponseMeta
	err = json.Unmarshal(chibisafe_post, &chibisafe_Response_Metadata)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("[%s] [%s] [%s] : Successfully obtained PUT keys", r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath)

	_, err = handler.NetworkStoragePut(chibisafe_Response_Metadata.URL, PostData.ContentType, tempfilepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[%s] [%s] [%s] : %s", r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath, err)
		return
	}
	log.Printf("[%s] [%s] [%s] : Successfully PUT file to Network Storage", r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath)

	PostProcessData := handler.UploadProcessMeta{
		Name:        fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Identifier:  chibisafe_Response_Metadata.Identifier,
	}

	PostProcess, err := handler.UploadProcessPost(Chibisafe_basepath, PostData.ContentType, chibisafe_Response_Metadata.Identifier, API_Key, PostProcessData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("[%s] [%s] [%s] : %s", r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath, err)
		return
	}

	var PostProcessResponse handler.UploadProcessResponseMeta
	err = json.Unmarshal(PostProcess, &PostProcessResponse)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("[%s] [%s] [%s] : Successfully Processed Response with UUID: %s", r.RemoteAddr, PostProcessResponse.Name, tempfilepath, PostProcessResponse.UUID)

	err = handler.DeleteFile(tempfilepath)
	if err != nil {
		log.Printf("[%s] [%s] [%s] : %s", r.RemoteAddr, chibisafe_Response_Metadata.Identifier, tempfilepath, err)
		return
	}
	log.Printf("[%s] [%s] : Successfully Deleted Temporary file from local disk", r.RemoteAddr, tempfilepath)
	JsonResponse, _ := json.Marshal(PostProcessResponse)
	fmt.Fprintf(w, "%s", JsonResponse)
}

func main() {
	Chibisafe_basepath := os.Getenv("CHIBISAFE_BASEPATH")

	if Chibisafe_basepath == "" {
		log.Fatal("CHIBISAFE_BASEPATH environment is not set!")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/upload", uploadHandler)

	if err := http.ListenAndServe(":4040", mux); err != nil {
		log.Fatal(err)
	}
}
