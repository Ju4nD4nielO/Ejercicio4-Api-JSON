package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Message string `json:"message"`
}

type Health struct {
	Status string `json:"status"`
}

type Item struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Console   string `json:"console"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
	Genre     string `json:"genre"`
}

func main() {
	http.HandleFunc("/api/ping", pingHandler)
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/items", itemsHandler)

	log.Println("JSON API running on :24979")
	log.Fatal(http.ListenAndServe(":24979", nil))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := Message{
		Message: "pong",
	}

	writeJSON(w, http.StatusOK, response)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	response := Message{
		Message: "Hello from pure JSON API",
	}

	writeJSON(w, http.StatusOK, response)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {

	data, err := os.ReadFile("data/items.json")

	if err != nil {
		http.Error(w, "Could not read data", http.StatusInternalServerError)
		return
	}

	var items []Item

	err = json.Unmarshal(data, &items)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
