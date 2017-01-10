package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const tempFilePath = "/tmp/tmp_image.jpg"
const imageFilePath = "/tmp/image.jpg"
const errorImageFilePath = "/tmp/error_image.jpg"
const minimumPhotoInterval = 10 * time.Second

var photoRequests = make(chan int)
var photoResponses = make(chan string)

func takePicture(path string) error {
	cmd := exec.Command("raspistill", "-o", path)
	_, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	err = os.Rename(tempFilePath, imageFilePath)
	if err != nil {
		return err
	}
	return nil
}

func photographer() {
	for {
		<-photoRequests
		// If there is a very recent photo, just return it
		fi, err := os.Stat(imageFilePath)
		if err == nil && time.Since(fi.ModTime()) < minimumPhotoInterval {
			photoResponses <- imageFilePath
			continue
		}

		// Else take a new one into a temp file
		if err := takePicture(tempFilePath); err != nil {
			log.Println("Error taking picture: ", err)
			photoResponses <- errorImageFilePath
			continue
		}

		// Rename the temp file to the image file (unix atomic replace)
		if err := os.Rename(tempFilePath, imageFilePath); err != nil {
			log.Println("Error renaming temp file: ", err)
			photoResponses <- errorImageFilePath
			continue
		}

		photoResponses <- imageFilePath
	}
}

func getPicFromPhotographer() string {
	photoRequests <- 1
	resp := <-photoResponses
	return resp
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s\n", ts)
}

func camHandler(w http.ResponseWriter, r *http.Request) {
	filePath := getPicFromPhotographer()
	http.ServeFile(w, r, filePath)
}

func main() {
	go photographer()
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/image.jpg", camHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
