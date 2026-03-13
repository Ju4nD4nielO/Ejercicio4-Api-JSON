package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	http.HandleFunc("/api/items/", itemByIDHandler)
	http.HandleFunc("/api/items/delete/", deleteItemHandler)

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

	// lee query parameter combinados
	console := r.URL.Query().Get("console")
	yearParam := r.URL.Query().Get("year")
	idParam := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idParam)

	// if err != nil {
	// 	http.Error(w, "Invalid id", http.StatusBadRequest)
	// 	return
	// }

	// aplicar filtros
	filtered := []Item{}

	for _, item := range items {

		if console != "" && item.Console != console {
			continue
		}

		if yearParam != "" {
			year, err := strconv.Atoi(yearParam)

			if err == nil && item.Year != year {
				continue
			}
		}

		filtered = append(filtered, item)
	}

	if console != "" || yearParam != "" {
		writeJSON(w, http.StatusOK, filtered)
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

	newItem.ID = items[len(items)-1].ID + 1

	items = append(items, newItem)

	updatedData, err := json.MarshalIndent(items, "", " ")

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile("data/items.json", updatedData, 0644)

	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, newItem)
}

func itemByIDHandler(w http.ResponseWriter, r *http.Request) {

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

	// extraer id del path
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idParam := parts[3]

	id, err := strconv.Atoi(idParam)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, item := range items {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// obtener id del path
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idParam := parts[4]

	id, err := strconv.Atoi(idParam)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// buscar y eliminar
	index := -1

	for i, item := range items {
		if item.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// eliminar elemento
	items = append(items[:index], items[index+1:]...)

	// convertir a JSON
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

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Item deleted successfully",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
