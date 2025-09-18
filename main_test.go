package main

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/bitterfq/data-ingestion-go/internal/db"
	"github.com/bitterfq/data-ingestion-go/parts"
	"github.com/bitterfq/data-ingestion-go/suppliers"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, *db.Queries) {
	t.Helper()

	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("read schema.sql: %v", err)
	}
	if _, err := conn.Exec(string(schema)); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	return conn, db.New(conn)
}

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	conn, q := setupTestDB(t)
	tenant := "tenant_test"

	// 1. generate + insert suppliers
	sups := suppliers.GenerateSuppliers(tenant, 5)
	for _, sup := range sups {
		_, err := q.CreateSupplier(ctx, db.CreateSupplierParams{
			SupplierID:         sup.SupplierID,
			TenantID:           sup.TenantID,
			LegalName:          sup.LegalName,
			SchemaVersion:      sql.NullString{String: sup.SchemaVersion, Valid: true},
			IngestionTimestamp: sql.NullTime{Time: sup.IngestionTimestamp, Valid: true},
			SourceTimestamp:    sql.NullTime{Time: sup.SourceTimestamp, Valid: true},
		})
		if err != nil {
			t.Fatalf("insert supplier: %v", err)
		}
	}

	// 2. collect supplier IDs + insert parts
	var supplierIDs []string
	for _, s := range sups {
		supplierIDs = append(supplierIDs, s.SupplierID)
	}
	partsList := parts.GenerateParts(10, tenant, supplierIDs)
	for _, part := range partsList {
		_, err := q.CreatePart(ctx, db.CreatePartParams{
			PartID:             part.PartID,
			TenantID:           part.TenantID,
			PartNumber:         part.PartNumber,
			Description:        part.Description,
			SchemaVersion:      sql.NullString{String: part.SchemaVersion, Valid: true},
			IngestionTimestamp: sql.NullTime{Time: part.IngestionTimestamp, Valid: true},
			SourceTimestamp:    sql.NullTime{Time: part.SourceTimestamp, Valid: true},
		})
		if err != nil {
			t.Fatalf("insert part: %v", err)
		}
	}

	// 3. verify counts with raw SQL
	var supCount, partCount int
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM dim_supplier_v1").Scan(&supCount); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM dim_part_v1").Scan(&partCount); err != nil {
		t.Fatal(err)
	}

	if supCount != len(sups) {
		t.Errorf("expected %d suppliers, got %d", len(sups), supCount)
	}
	if partCount != len(partsList) {
		t.Errorf("expected %d parts, got %d", len(partsList), partCount)
	}
}
