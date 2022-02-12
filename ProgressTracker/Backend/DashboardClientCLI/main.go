package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var aa = flag.String("addr", "75.119.149.184:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dashboardURL := url.URL{Scheme: "ws", Host: *aa, Path: "/dashboard"}

	clientKontrolka, _, err := websocket.DefaultDialer.Dial(dashboardURL.String(), nil)
	if err != nil {
		log.Fatal("clientKontrolka dial:", err)
	}
	defer clientKontrolka.Close()

	go func() {
		<-interrupt
		err := clientKontrolka.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		os.Exit(0)
	}()

	for {
		_, message, _ := clientKontrolka.ReadMessage()
		log.Println(string(message))
	}
}
