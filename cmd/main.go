package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pankajhirway/ludo-server/internal/players"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var playerMap map[string]players.Player

func init() {
	playerMap = make(map[string]players.Player)
}

var upgrader = websocket.Upgrader{} // use default options

func register(w http.ResponseWriter, r *http.Request) {
	userid := uuid.New()
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("upgrade:", err)
		return
	}
	player := Player{}
	player.Set(c, userid.String())
	go player.StartListening()
	playerMap[userid.String()] = player
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func unregister(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userid := string(body)
	if val, ok := playerMap[userid]; ok {
		val.c.Close()
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/register", register)
	http.HandleFunc("/unregister", unregister)
	log.Printf("Server started on : %X", addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
