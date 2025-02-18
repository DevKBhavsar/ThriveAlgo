package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Holiday struct {
	ID          string    `json:"id" bson:"id"`
	Date        time.Time `json:"date" bson:"date"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) error {
	id := path.Base(r.URL.Path)
	holiday, err := retrieve(id)
	if err != nil {
		return err
	}
	if holiday.ID == "" {
		http.Error(w, "Holiday not found", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(holiday)
	return nil
}

func handlePost(w http.ResponseWriter, r *http.Request) error {
	var holiday Holiday
	err := json.NewDecoder(r.Body).Decode(&holiday)
	if err != nil {
		return err
	}

	err = holiday.create()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(holiday)
	return nil
}

func handlePut(w http.ResponseWriter, r *http.Request) error {
	id := path.Base(r.URL.Path)
	holiday, err := retrieve(id)
	if err != nil {
		return err
	}
	if holiday.ID == "" {
		http.Error(w, "Holiday not found", http.StatusNotFound)
		return nil
	}

	err = json.NewDecoder(r.Body).Decode(&holiday)
	if err != nil {
		return err
	}

	err = holiday.update()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func handleDelete(w http.ResponseWriter, r *http.Request) error {
	id := path.Base(r.URL.Path)
	err := deleteHoliday(id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func main() {
	client = InitMongoClient()

	http.HandleFunc("/api/holidays/", handleRequest)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
