package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Holiday struct {
	ID    string `json:"id" bson:"_id"`
	Date  string `json:"date" bson:"date"`
	Title string `json:"title" bson:"title"`
}

func getAllHoliday(w http.ResponseWriter, r *http.Request) {

	holidays, err := retrieveAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holidays)
}

func createHolidayHandler(w http.ResponseWriter, r *http.Request) {
	var holiday Holiday
	if err := json.NewDecoder(r.Body).Decode(&holiday); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a new unique ID using MongoDB's ObjectID
	holiday.ID = primitive.NewObjectID().Hex()

	if err := holiday.create(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(holiday)
}

func deleteHolidayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := deleteHoliday(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	client = InitMongoClient()
	r := mux.NewRouter()
	r.HandleFunc("/api/holidays", getAllHoliday).Methods("GET")
	r.HandleFunc("/api/holidays", createHolidayHandler).Methods("POST")
	r.HandleFunc("/api/holidays/{id}", deleteHolidayHandler).Methods("DELETE")

	// Create a CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Add your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          true, // Enable Debugging for testing, consider disabling in production
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler)) // Use the wrapped handler instead of r
}
