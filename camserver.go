package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(time.Second * 5)
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s\n", ts)
}

func main() {
	http.HandleFunc("/time", timeHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
