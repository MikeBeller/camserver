package main

import (
	"fmt"
	"log"
	"net/http"
	_ "os"
	"time"
)

const tempFilePath = "/tmp/tmp_image.jpg"
const imageFilePath = "/tmp/image.jpg"

var photoRequests = make(chan int)
var photoResponses = make(chan string)

func photographer() {
	for {
		<-photoRequests
		time.Sleep(time.Second * 5)
		photoResponses <- imageFilePath
	}
}

func takePic() string {
	photoRequests <- 1
	resp := <-photoResponses
	return resp
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s\n", ts)
}

func camHandler(w http.ResponseWriter, r *http.Request) {
	filePath := takePic()
	http.ServeFile(w, r, filePath)
}

func main() {
	go photographer()
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/image.jpg", camHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
