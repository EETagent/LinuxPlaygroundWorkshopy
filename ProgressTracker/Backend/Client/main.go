package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "75.119.149.184:8080", "http service address")

func checkProgram(name string, args []string) (string, error) {
	result, err := exec.Command(name, args...).Output()
	return string(result), err
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	hostname, _ := os.Hostname()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	kontrolkaURL := url.URL{Scheme: "ws", Host: *addr, Path: "/report/kontrolka"}
	ukladatkoURL := url.URL{Scheme: "ws", Host: *addr, Path: "/report/ukladatko"}

	clientKontrolka, _, err := websocket.DefaultDialer.Dial(kontrolkaURL.String(), nil)
	if err != nil {
		log.Fatal("clientKontrolka dial:", err)
	}
	defer clientKontrolka.Close()

	clientUkladatko, _, err := websocket.DefaultDialer.Dial(ukladatkoURL.String(), nil)
	if err != nil {
		log.Fatal("clientUkladatko dial:", err)
	}
	defer clientUkladatko.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	err = clientKontrolka.WriteMessage(websocket.TextMessage, []byte(hostname))
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	err = clientUkladatko.WriteMessage(websocket.TextMessage, []byte(hostname))
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			kontrolkaOutput, _ := checkProgram("/bin/kontrolka", nil)
			ukladatkoOutput, _ := checkProgram("/bin/ukladatko", nil)

			err := clientKontrolka.WriteMessage(websocket.TextMessage, []byte(kontrolkaOutput))
			if err != nil {
				log.Println("write:", err)
				return
			}
			err = clientUkladatko.WriteMessage(websocket.TextMessage, []byte(ukladatkoOutput))
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-interrupt:
			err := clientKontrolka.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			clientKontrolka.Close()
			if err != nil {
				log.Println("write close:", err)
				return
			}
			err = clientUkladatko.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			clientUkladatko.Close()
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
