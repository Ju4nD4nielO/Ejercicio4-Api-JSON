package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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
	http.HandleFunc("/api/items/create", createItemHandler)

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

	// lee query parameter
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		writeJSON(w, http.StatusOK, items)
		return
	}

	// convertir id a int
	id, err := strconv.Atoi(idParam)

	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// buscar item
	for _, item := range items {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)

}

func createItemHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// leer archivo
	data, err := os.ReadFile("data/items.json")

	if err != nil {
		http.Error(w, "Could not read data", http.StatusInternalServerError)
		return
	}

	var items []Item

	err = json.Unmarshal(data, &items)

	if err != nil {
		http.Error(w, "Invalid JSON file", http.StatusInternalServerError)
		return
	}

	// leer JSON enviado por cliente
	var newItem Item

	err = json.NewDecoder(r.Body).Decode(&newItem)

	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// generar nuevo id
	newItem.ID = items[len(items)-1].ID + 1

	// agregar a lista
	items = append(items, newItem)

	// convertir lista a JSON
	updatedData, err := json.MarshalIndent(items, "", " ")

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// guardar archivo
	err = os.WriteFile("data/items.json", updatedData, 0644)

	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, newItem)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
