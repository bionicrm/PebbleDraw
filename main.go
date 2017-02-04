package main

import (
	"net/http"
	"io"
	"encoding/json"
	"log"
)

func main() {
	http.HandleFunc("/", homeHandler)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalln(err)
	}
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