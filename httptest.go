package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, "%s\n", ts)
}

func main() {
	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, r)
			}))
	resp, _ := http.Get(srv.URL)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
