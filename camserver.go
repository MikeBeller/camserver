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

func takePicture() (string, error) {
	// Need to include a throttle that checks age of file and
	// doesn't retake if it's less than minimumPhotoInterval old
	cmd := exec.Command("raspistill", "-o", tempFilePath)
	_, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	err = os.Rename(tempFilePath, imageFilePath)
	if err != nil {
		return "", err
	}
	return imageFilePath, nil
}

func photographer() {
	for {
		<-photoRequests
		path, err := takePicture()
		if err != nil {
			log.Println("Error taking picture: ", err)
			photoResponses <- errorImageFilePath
		}

		photoResponses <- path
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
