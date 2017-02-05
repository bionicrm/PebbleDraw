package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"io/ioutil"
)

const TokenKey = "token"

type webWsHandler func(deviceInfo)

var tokenHandlerMap = map[string]func(deviceInfo, drawingInfo) error{}

type deviceInfo struct {
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Z            float64 `json:"z"`
	Tapped       bool    `json:"tapped"`
	ClickUp      bool    `json:"click_up"`
	ClickSelect  bool    `json:"click_select"`
	ClickDown    bool    `json:"click_down"`
	ClickBack    bool    `json:"click_back"`
}

type drawingInfo struct {
	Drawing bool `json:"drawing"`
	Hue     int  `json:"hue"`
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

	drawing := false
	hue := 0

	var incrementHue = func(amount int) {
		hue += amount
		if hue > 360 {
			hue -= 360
		} else if hue < 0 {
			hue += 360
		}
	}

	for {
		deviceI := deviceInfo{}
		if err := ws.ReadJSON(&deviceI); err != nil {
			log.Println(err)
			return
		}
		if deviceI.ClickUp {
			incrementHue(360 / 10)
		}
		if deviceI.ClickSelect {
			drawing = !drawing
		}
		if deviceI.ClickDown {
			incrementHue(-360 / 10)
		}
		if deviceI.ClickBack {
			drawing = false
		}
		drawingI := drawingInfo{
			Drawing: drawing,
			Hue:     hue,
		}
		if err := ws.WriteJSON(drawingI); err != nil {
			log.Println(err)
			return
		}
		if handler, ok := tokenHandlerMap[token]; ok {
			if err := handler(deviceI, drawingI); err != nil {
				log.Println(err)
				return
			}
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

	tokenHandlerMap[token] = func(deviceI deviceInfo, drawingI drawingInfo) error {
		return ws.WriteJSON(struct{
			DeviceI  deviceInfo  `json:"device_i"`
			DrawingI drawingInfo `json:"drawing_i"`
		}{
			DeviceI:  deviceI,
			DrawingI: drawingI,
		})
	}
}