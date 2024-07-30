package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net/http"

	"github.com/gorilla/mux"        // For routing
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Item struct to represent data in the database and API responses
type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Database connection (you'll initialize it later)
var db *sql.DB

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)

		// body, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, "Error reading request body", http.StatusInternalServerError)
		// 	return
		// }

		// // Use the 'body' byte slice
		// log.Println("Request body: %+v\n", string(body))

		// ... Call the next handler, which can be another middleware in the chain, or the final handler.

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Open/Create SQLite database
	var err error
	db, err = sql.Open("sqlite3", "../items.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, name TEXT, description TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP router
	router := mux.NewRouter()

	// Define API endpoints (handlers)
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/items", createItem).Methods("POST")
	router.HandleFunc("/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	router.Use(loggingMiddleware)

	// Start server
	log.Println("Server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Handler functions for API endpoints:
func getItems(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "../items.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		err := rows.Scan(&i.Id, &i.Name, &i.Description)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to encode DB items as JSON", http.StatusInternalServerError)
		return
	}
}

func createItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Now you have the JSON data bound to the 'item' struct
	log.Printf("Received Item: %+v\n", item)

	db, err = sql.Open("sqlite3", "../items.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("INSERT INTO items(id, name, description) VALUES(%v, '%v', '%v')", item.Id, item.Name, item.Description))
	if err != nil {
		log.Fatal(err)
	}

	// ... (Do something with the item data)
}
func getItem(w http.ResponseWriter, r *http.Request)    { /* ... */ }
func updateItem(w http.ResponseWriter, r *http.Request) { /* ... */ }
func deleteItem(w http.ResponseWriter, r *http.Request) { /* ... */ }
