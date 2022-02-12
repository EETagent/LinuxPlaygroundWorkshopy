package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var server = flag.String("server", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{}

var database = []ChallengeClient{}
var databaseMutex = sync.Mutex{}
var dbChanged = make(chan bool)

var clients []*websocket.Conn
var clientsMutex = sync.Mutex{}

type ChallengeData struct {
	Hostname string `json:"hostname"`
	Command  string `json:"command"`
}
type ChallengeClient struct {
	Challenge string        `json:"challenge"`
	Data      ChallengeData `json:"data"`
}

func getResponseJson() []byte {
	responseJSON, _ := json.Marshal(database)
	return responseJSON
}

func dashboardWebSocketBroadcast() {
	for {
		<-dbChanged
		clientsMutex.Lock()
		for i := range clients {
			if clients[i] != nil {
				err := clients[i].WriteMessage(websocket.TextMessage, getResponseJson())
				if err != nil {
					log.Printf(err.Error())
					clients = append(clients[:i], clients[i+1:]...)
				}
			}
		}
		clientsMutex.Unlock()
	}
}

func dashboardWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}
	clientsMutex.Lock()
	clients = append(clients, c)
	id := len(clients) - 1
	clientsMutex.Unlock()
	defer func(c *websocket.Conn) {
		clientsMutex.Lock()
		clients = append(clients[:id], clients[id+1:]...)
		clientsMutex.Unlock()
		err := c.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(c)
	_ = c.WriteMessage(websocket.TextMessage, getResponseJson())
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func progressReporterWebSocket(w http.ResponseWriter, r *http.Request, challenge string) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	_, hostname, err := c.ReadMessage()
	if err != nil {
		return
	}
	hostnameString := string(hostname)
	if len(hostnameString) > 30 {
		hostnameString = hostnameString[:30]
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			databaseMutex.Lock()
			for i := range database {
				if database[i].Challenge == challenge {
					if database[i].Data.Hostname == hostnameString {
						database = append(database[:i], database[i+1:]...)
						break
					}
				}
			}
			dbChanged <- true
			databaseMutex.Unlock()
			break
		}
		output := string(message)
		isAlreadyInDatabase := false
		databaseMutex.Lock()
		for i := range database {
			if database[i].Challenge == challenge {
				if database[i].Data.Hostname == hostnameString {
					//if database[i].Data.Command == output {
					isAlreadyInDatabase = true
					database[i].Data.Command = output
					dbChanged <- true
					break
					//}
				}
			}
		}
		if isAlreadyInDatabase != true {
			if output != "" {
				challengeClient := ChallengeClient{Challenge: challenge, Data: ChallengeData{
					Hostname: hostnameString,
					Command:  output,
				}}
				database = append(database, challengeClient)
				dbChanged <- true
			}
		}
		databaseMutex.Unlock()
		//log.Println(output)
	}
}

func main() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	flag.Parse()
	log.SetFlags(0)
	go dashboardWebSocketBroadcast()
	http.HandleFunc("/report/kontrolka", func(w http.ResponseWriter, r *http.Request) {
		progressReporterWebSocket(w, r, "kontrolka")
	})
	http.HandleFunc("/report/ukladatko", func(w http.ResponseWriter, r *http.Request) {
		progressReporterWebSocket(w, r, "ukladatko")
	})
	http.HandleFunc("/dashboard", dashboardWebSocket)
	log.Fatal(http.ListenAndServe(*server, nil))
}
