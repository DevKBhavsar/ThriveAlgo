// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Holiday struct {
	ID    string    `json:"id" bson:"_id,omitempty"`
	Date  time.Time `json:"date" bson:"date"`
	Title string    `json:"title" bson:"title"`
}

func getAllHoliday(w http.ResponseWriter, r *http.Request) {
	holidays, err := retrieveAll()
	if err != nil {
		log.Printf("Error retrieving holidays: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holidays)
}

func createHolidayHandler(w http.ResponseWriter, r *http.Request) {
	var holidayRequest struct {
		Title string `json:"title"`
		Date  string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&holidayRequest); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Parse the date from string
	parsedDate, err := time.Parse("2006-01-02", holidayRequest.Date)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	holiday := Holiday{
		ID:    primitive.NewObjectID().Hex(),
		Date:  parsedDate,
		Title: holidayRequest.Title,
	}

	if err := holiday.create(); err != nil {
		log.Printf("Error creating holiday: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(holiday)
}

func deleteHolidayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := deleteHoliday(id); err != nil {
		log.Printf("Error deleting holiday: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Holiday deleted successfully"}`))
}

func main() {
	// Initialize MongoDB client

	client = InitMongoClient()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/api/holidays", getAllHoliday).Methods("GET")
	r.HandleFunc("/api/holidays", createHolidayHandler).Methods("POST")
	r.HandleFunc("/api/holidays/{id}", deleteHolidayHandler).Methods("DELETE")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap router with CORS handler
	handler := c.Handler(r)

	// Start server
	fmt.Println("Server running on por number", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
