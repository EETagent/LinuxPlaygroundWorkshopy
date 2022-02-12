package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebsocketDashboard(t *testing.T) {
	dashboardWebsocket := httptest.NewServer(http.HandlerFunc(dashboardWebSocket))
	defer dashboardWebsocket.Close()
	dashboardURL := "ws" + strings.TrimPrefix(dashboardWebsocket.URL, "http")

	dashboardOnlyClient, _, err := websocket.DefaultDialer.Dial(dashboardURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dashboardOnlyClient.Close()

	err = dashboardOnlyClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestWebsocketReporter(t *testing.T) {
	kontrolkaWebsocket := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		progressReporterWebSocket(w, r, "kontrolka")
	}))
	ukladatkoWebsocket := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		progressReporterWebSocket(w, r, "ukladatko")
	}))
	defer kontrolkaWebsocket.Close()
	defer ukladatkoWebsocket.Close()

	kontrolkaURL := "ws" + strings.TrimPrefix(kontrolkaWebsocket.URL, "http")
	ukladatkoURL := "ws" + strings.TrimPrefix(ukladatkoWebsocket.URL, "http")

	kontrolkaOnlyClient, _, err := websocket.DefaultDialer.Dial(kontrolkaURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer kontrolkaOnlyClient.Close()
	ukladatkoOnlyClient, _, err := websocket.DefaultDialer.Dial(ukladatkoURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ukladatkoOnlyClient.Close()

	err = ukladatkoOnlyClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = kontrolkaOnlyClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestWebsocketBackend(t *testing.T) {
	t.Run("kontrolka", func(t *testing.T) {
		CheckWebsocketBackend(t, "hostname", "test", "kontrolka")
	})
	t.Run("ukladatko", func(t *testing.T) {
		CheckWebsocketBackend(t, "hostname", "test", "ukladatko")
	})
	t.Run("utf8", func(t *testing.T) {
		CheckWebsocketBackend(t, "aahdd/ ýčýč __kdlmd", "oufuuáěuárěrnjsdj nf sdkf jn xffjwsm jbm fxjc", "kontrolka")
	})
	t.Run("loop", func(t *testing.T) {
		for i := 0; i < 1500; i++ {
			CheckWebsocketBackend(t, "NekonecnaBaterie", "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aliquam erat volutpat. Integer in sapien. Etiam ligula pede, sagittis quis, interdum ultricies, scelerisque eu. Fusce tellus odio, dapibus id fermentum quis, suscipit id erat. Fusce dui leo, imperdiet in, aliquam sit amet, feugiat eu, orci. Proin mattis lacinia justo. Etiam posuere lacus quis dolor. Etiam egestas wisi a erat. Mauris suscipit, ligula sit amet pharetra semper, nibh ante cursus purus, vel sagittis velit mauris vel metus. Nulla accumsan, elit sit amet varius semper, nulla mauris mollis quam, tempor suscipit diam nulla vel leo.", "kontrolka")
		}
	})
}

func CheckWebsocketBackend(t *testing.T, hostname string, output string, challenge string) {
	go dashboardWebSocketBroadcast()
	dashboardWebsocket := httptest.NewServer(http.HandlerFunc(dashboardWebSocket))
	defer dashboardWebsocket.Close()
	dashboardURL := "ws" + strings.TrimPrefix(dashboardWebsocket.URL, "http")

	dashboardOnlyClient, _, err := websocket.DefaultDialer.Dial(dashboardURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dashboardOnlyClient.Close()

	challengeReporterWebsocket := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		progressReporterWebSocket(w, r, challenge)
	}))
	defer challengeReporterWebsocket.Close()
	challengeReporterURL := "ws" + strings.TrimPrefix(challengeReporterWebsocket.URL, "http")

	challengeReporterOnlyClient, _, err := websocket.DefaultDialer.Dial(challengeReporterURL, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer challengeReporterOnlyClient.Close()

	_, response, err := dashboardOnlyClient.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(response) != "[]" {
		log.Fatalf("Invalid response: %s", string(response))
	}

	err = challengeReporterOnlyClient.WriteMessage(websocket.TextMessage, []byte(hostname))
	if err != nil {
		t.Errorf("%v", err)
	}

	err = challengeReporterOnlyClient.WriteMessage(websocket.TextMessage, []byte(output))
	if err != nil {
		t.Errorf("%v", err)
	}

	_, response, err = dashboardOnlyClient.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(response) != "[{\"challenge\":\""+challenge+"\",\"data\":{\"hostname\":\""+hostname+"\",\"command\":\""+output+"\"}}]" {
		log.Fatalf("Invalid response: %s", string(response))
	}

	alwaysValidClientStruct := ChallengeClient{Challenge: challenge, Data: ChallengeData{
		Hostname: hostname,
		Command:  output,
	}}
	currentClientStruct := database[0]
	if alwaysValidClientStruct.Challenge != currentClientStruct.Challenge || alwaysValidClientStruct.Data.Command != currentClientStruct.Data.Command || alwaysValidClientStruct.Data.Hostname != currentClientStruct.Data.Hostname {
		t.Fatalf("Structs does not match %v %v", alwaysValidClientStruct, currentClientStruct)
	}

	err = challengeReporterOnlyClient.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, response, err = dashboardOnlyClient.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if string(response) != "[]" {
		log.Fatalf("Invalid response: %s", string(response))
	}
}
