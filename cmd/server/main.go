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
	"strings"
	"time"

	"github.com/bitterfq/data-ingestion-go/internal/database/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"

	"github.com/bitterfq/data-ingestion-go/internal/snowflake"
)

func main() {

	// 1. connect to db
	conn, err := sql.Open("sqlite3", "data/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// create schema if it doesn't exist
	schema, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Exec(string(schema)); err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatal(err)
	}

	q := db.New(conn)
	ctx := context.Background() //look into this

	log.Printf("Starting server on :8080")
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
			log.Printf("failed to delete supplier: %v", err)
			http.Error(w, "Failed to delete supplier", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	// delete /parts
	mux.HandleFunc("/parts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Path[len("/parts/"):]
		if id == "" {
			http.Error(w, "Part ID is required", http.StatusBadRequest)
			return
		}
		if err := q.DeletePart(r.Context(), id); err != nil {
			log.Printf("failed to delete part %s: %v", id, err)
			http.Error(w, "Failed to delete part", http.StatusInternalServerError)
			return
		}
		log.Printf("deleted part id=%s", id)
		w.WriteHeader(http.StatusNoContent)
	})

	// fetch from snowflake and insert into db
	mux.HandleFunc("/fetch-and-insert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Connect to Snowflake
		sfDB, err := snowflake.NewClient("")
		if err != nil {
			http.Error(w, "Failed to connect to Snowflake: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer sfDB.Close()

		rows, err := sfDB.Query("SELECT SUPPLIER_ID, TENANT_ID, SUPPLIER_CODE, LEGAL_NAME, DBA_NAME,COUNTRY, REGION, ADDRESS_LINE1, ADDRESS_LINE2, CITY, STATE, POSTAL_CODE FROM SUPPLY_CHAIN.PUBLIC.SUPPLIERS")
		if err != nil {
			http.Error(w, "Failed to query Snowflake: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		//var insertedSuppliers []db.Supplier
		for rows.Next() {
			var id, tenantID, legalName string
			var supplierCode, dbaName, country, region, addressLine1,
				addressLine2, city, state, postalCode sql.NullString

			if err := rows.Scan(
				&id, &tenantID, &supplierCode, &legalName,
				&dbaName, &country, &region, &addressLine1,
				&addressLine2, &city, &state, &postalCode,
			); err != nil {
				http.Error(w, "Failed to scan row: "+err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("inserting supplier %s (%s)", id, legalName)

			_, err := q.CreateSupplier(ctx, db.CreateSupplierParams{
				SupplierID:   id,
				TenantID:     tenantID,
				SupplierCode: supplierCode,
				LegalName:    legalName,
				DbaName:      dbaName,
				Country:      country,
				Region:       region,
				AddressLine1: addressLine1,
				AddressLine2: addressLine2,
				City:         city,
				State:        state,
				PostalCode:   postalCode,
			})
			if err != nil {
				http.Error(w, "Failed to insert supplier into local DB: "+err.Error(), http.StatusInternalServerError)
				return
			}

			//insertedSuppliers = append(insertedSuppliers, sup)

		}

	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
