package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var photographer *Photographer

func timeHandler(w http.ResponseWriter, r *http.Request) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s\n", ts)
}

func camHandler(w http.ResponseWriter, r *http.Request) {
	response := photographer.GetPic()
	//http.ServeFile(w, r, filePath)
	fmt.Fprintf(w, "%v\n", response)
}

func main() {
	photographer = NewPhotographer()
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/image.jpg", camHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
