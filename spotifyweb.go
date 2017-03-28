package main

import (
	"fmt"
	"github.com/vansante/go-spotify-control"
	"log"
	"net/http"
	"time"
	"encoding/json"
)

const PORT = 1337

var control *spotifycontrol.SpotifyControl

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleStatus)
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/pause", handlePause)
	mux.HandleFunc("/play", handlePlay)
	mux.HandleFunc("/restart", handleRestart)

	addr := fmt.Sprintf(":%d", PORT)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	var err error
	control, err = spotifycontrol.NewSpotifyControl("", 1 * time.Second)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error starting control: %v", err)))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("OK!"))
}

func handlePause(w http.ResponseWriter, r *http.Request) {
	if control == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	status, err := control.SetPauseState(r.URL.Query().Get("paused") == "true")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error (un)pausing: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error unmarshalling: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	if control == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, err := control.Play(url)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error playing: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error unmarshalling: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if control == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	status, err := control.GetStatus()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error getting status: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error unmarshalling: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}