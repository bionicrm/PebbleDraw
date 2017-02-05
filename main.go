package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"io/ioutil"
)

const TokenKey = "token"

type webWsHandler func(wsData)

var tokenHandlerMap = map[string]func(wsData){}

type wsData struct {
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Z            float64 `json:"z"`
	TappedX      bool    `json:"tapped_x"`
	TappedY      bool    `json:"tapped_y"`
	TappedZ      bool    `json:"tapped_z"`
	TopButton    bool    `json:"top_button"`
	MiddleButton bool    `json:"middle_button"`
	BottomButton bool    `json:"bottom_button"`
	LeftButton   bool    `json:"left_button"`
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/script.js", jsHandler)
	http.HandleFunc("/ws/pebble", pebbleHandler)
	http.HandleFunc("/ws/web", webHandler)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalln(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	serveFile(w, "web/index.html")
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	serveFile(w, "web/script.js")
}

var upgrader = websocket.Upgrader{}

func serveFile(w http.ResponseWriter, filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := w.Write(b); err != nil {
		log.Println(err)
	}
}

func pebbleHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}
	token := r.Form.Get(TokenKey)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		data := wsData{}
		if err := ws.ReadJSON(data); err != nil {
			log.Println(err)
			continue
		}
		if handler, ok := tokenHandlerMap[token]; ok {
			handler(data)
		}
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}
	token := r.Form.Get(TokenKey)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	tokenHandlerMap[token] = func(data wsData) {
		if err := ws.WriteJSON(data); err != nil {
			log.Println(err)
		}
	}
}