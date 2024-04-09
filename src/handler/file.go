package handler

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func GetTempFilename(Filename string) string {
	filename := fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(Filename))
	return filename
}

func SaveFile(filename string, file multipart.File) error {
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Panic(err)
		return err
	}

	dst, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
		return err
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

func DiscardFile(file multipart.File) error {
	// Clear the data from memory
	_, err := io.Copy(io.Discard, file)
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
