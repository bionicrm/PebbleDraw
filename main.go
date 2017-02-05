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
	Data []struct{
		X            float64 `json:"x"`
		Y            float64 `json:"y"`
		Z            float64 `json:"z"`
		Vibe         bool    `json:"vibe"`
		Time         uint64  `json:"time"`
	}                    `json:"data"`
	Tapped       bool    `json:"tapped"`
	ClickUp      bool    `json:"click_up"`
	ClickSelect  bool    `json:"click_select"`
	ClickDown    bool    `json:"click_down"`
	ClickBack    bool    `json:"click_back"`
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ws/pebble", pebbleHandler)
	http.HandleFunc("/ws/web", webHandler)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalln(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := w.Write(b); err != nil {
		log.Println(err)
	}
}

var upgrader = websocket.Upgrader{}

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
			return
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