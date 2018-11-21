package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var err error
	var content []byte
	if content, err = ioutil.ReadFile("./index.html"); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(content))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
