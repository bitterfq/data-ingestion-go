package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bitterfq/data-ingestion-go/internal/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
)

func main() {

	// 1. connect to db
	conn, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create schema if it doesn't exist
	schema, _ := os.ReadFile("schema.sql")
	conn.Exec(string(schema)) // create tables from schema.sql

	q := db.New(conn)
	ctx := context.Background() //look into this

	fmt.Println("Starting server on :8080")
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// POST /suppliers
	mux.HandleFunc("/suppliers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Handle supplier creation

		var body struct {
			TenantID  string `json:"tenant_id"`
			LegalName string `json:"legal_name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

		sup, err := q.CreateSupplier(ctx, db.CreateSupplierParams{
			SupplierID:   ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(), // Generate a unique ID
			TenantID:     body.TenantID,
			SupplierCode: sql.NullString{String: "code", Valid: true},
			LegalName:    body.LegalName,
		})

		if err != nil {
			http.Error(w, "Failed to insert supplier", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sup)
	})

	// POST /parts
	mux.HandleFunc("/parts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Handle part creation

		var body struct {
			TenantID    string `json:"tenant_id"`
			PartNumber  string `json:"part_number"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
		part, err := q.CreatePart(ctx, db.CreatePartParams{
			PartID:      ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String(), // Generate a unique ID
			TenantID:    body.TenantID,
			PartNumber:  body.PartNumber,
			Description: body.Description,
		})

		if err != nil {
			http.Error(w, "Failed to insert part", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(part)
	})

	// DELETE /suppliers/{id}
	mux.HandleFunc("/suppliers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// grab everything after /suppliers/
		id := r.URL.Path[len("/suppliers/"):]
		if id == "" {
			http.Error(w, "Supplier ID is required", http.StatusBadRequest)
			return
		}

		if err := q.DeleteSupplier(r.Context(), id); err != nil {
			log.Println("failed to delete supplier:", err)
			http.Error(w, "Failed to delete supplier", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	// delete /parts
	mux.HandleFunc("/parts/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[len("/parts/"):]
		if id == "" {
			http.Error(w, "Part ID is required", http.StatusBadRequest)
			return
		}

		if err := q.DeletePart(ctx, id); err != nil {
			http.Error(w, "Failed to delete part", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
