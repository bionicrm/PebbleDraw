package main

import (
	"net/http"
	"io"
	"encoding/json"
	"log"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}

type homeJson struct {
	SomeKey string `json:"some_key"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(homeJson{"some value"})
	if err != nil {
		log.Println(err)
	}
	io.WriteString(w, string(b))
}