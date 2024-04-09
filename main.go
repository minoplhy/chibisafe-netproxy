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

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Panic(err)
		return
	}
	defer file.Close()

	API_Key := r.Header.Get("x-api-key")
	if API_Key == "" {
		http.Error(w, "X-api-key is empty!", http.StatusBadRequest)
		log.Panicf("X-api-key is empty!")
		return
	}
	log.Printf("Received a successful POST from %s", r.RemoteAddr)
	tempfilepath := handler.GetTempFilename(fileHeader.Filename)
	log.Printf("Successfully obtained temporary Filename : %s", tempfilepath)
	handler.SaveFile(tempfilepath, file)
	handler.DiscardFile(file)

	PostData := handler.UploadPostMeta{
		ContentType: fileHeader.Header.Get("Content-Type"),
		Name:        fileHeader.Filename,
		FileSize:    fileHeader.Size,
	}

	chibisafe_post, _ := handler.UploadPost(Chibisafe_basepath, PostData, API_Key)
	var chibisafe_Response_Metadata handler.UploadResponseMeta
	err = json.Unmarshal(chibisafe_post, &chibisafe_Response_Metadata)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Printf("Successfully obtained PUT key with identifier: %s", chibisafe_Response_Metadata.Identifier)

	_, err = handler.NetworkStoragePut(chibisafe_Response_Metadata.URL, PostData.ContentType, tempfilepath)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Printf("Successfully PUT file to Network Storage with identifier: %s", chibisafe_Response_Metadata.Identifier)

	PostProcessData := handler.UploadProcessMeta{
		Name:        fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Identifier:  chibisafe_Response_Metadata.Identifier,
	}

	PostProcess, _ := handler.UploadProcessPost(Chibisafe_basepath, PostData.ContentType, chibisafe_Response_Metadata.Identifier, API_Key, PostProcessData)
	var PostProcessResponse handler.UploadProcessResponseMeta
	err = json.Unmarshal(PostProcess, &PostProcessResponse)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Printf("Successfully Processed Response with identifier: %s and UUID: %s", PostProcessResponse.Name, PostProcessResponse.UUID)

	err = handler.DeleteFile(tempfilepath)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Printf("Successfully Deleted Temporary file from local disk Filename: %s", tempfilepath)

	fmt.Fprintf(w, "%s", PostProcessResponse.URL)
}

func main() {
	Chibisafe_basepath := os.Getenv("CHIBISAFE_BASEPATH")
	Host := os.Getenv("HOST")

	if Chibisafe_basepath == "" {
		log.Fatal("CHIBISAFE_BASEPATH environment is not set!")
	}
	if Host == "" {
		Host = "127.0.0.1:4000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/upload", uploadHandler)

	if err := http.ListenAndServe(Host, mux); err != nil {
		log.Fatal(err)
	}
}
