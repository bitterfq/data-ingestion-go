package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bitterfq/data-ingestion-go/internal/db"
	"github.com/bitterfq/data-ingestion-go/internal/parts"
	"github.com/bitterfq/data-ingestion-go/internal/suppliers"
	_ "github.com/mattn/go-sqlite3"
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
	tenant := "tenant_acme"

	// 2. generate suppliers
	sups := suppliers.GenerateSuppliers(tenant, 10000)
	suppliers.SupplierWriter("suppliers.csv", sups)

	// 3. insert suppliers in a transaction
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	qtx := q.WithTx(tx)

	for _, sup := range sups {
		_, err := qtx.CreateSupplier(ctx, db.CreateSupplierParams{
			SupplierID:        sup.SupplierID,
			SupplierCode:      sql.NullString{String: sup.SupplierCode, Valid: sup.SupplierCode != ""},
			TenantID:          sup.TenantID,
			LegalName:         sup.LegalName,
			DbaName:           sql.NullString{String: sup.DBAName, Valid: sup.DBAName != ""},
			Country:           sql.NullString{String: sup.Country, Valid: sup.Country != ""},
			Region:            sql.NullString{String: sup.Region, Valid: sup.Region != ""},
			AddressLine1:      sql.NullString{String: sup.AddressLine1, Valid: sup.AddressLine1 != ""},
			AddressLine2:      sql.NullString{String: sup.AddressLine2, Valid: sup.AddressLine2 != ""},
			City:              sql.NullString{String: sup.City, Valid: sup.City != ""},
			State:             sql.NullString{String: sup.State, Valid: sup.State != ""},
			PostalCode:        sql.NullString{String: sup.PostalCode, Valid: sup.PostalCode != ""},
			ContactEmail:      sql.NullString{String: sup.ContactEmail, Valid: sup.ContactEmail != ""},
			ContactPhone:      sql.NullString{String: sup.ContactPhone, Valid: sup.ContactPhone != ""},
			PreferredCurrency: sql.NullString{String: sup.PreferredCurrency, Valid: sup.PreferredCurrency != ""},
			Incoterms:         sql.NullString{String: sup.Incoterms, Valid: sup.Incoterms != ""},
			LeadTimeDaysAvg:   sql.NullInt64{Int64: int64(sup.LeadTimeDaysAvg), Valid: true},
			LeadTimeDaysP95:   sql.NullInt64{Int64: int64(sup.LeadTimeDaysP95), Valid: true},
			OnTimeDeliveryRate: sql.NullFloat64{
				Float64: sup.OnTimeDeliveryRate,
				Valid:   true,
			},
			DefectRatePpm:        sql.NullInt64{Int64: int64(sup.DefectRatePPM), Valid: true},
			CapacityUnitsPerWeek: sql.NullInt64{Int64: int64(sup.CapacityUnitsPerWeek), Valid: true},
			RiskScore:            sql.NullFloat64{Float64: sup.RiskScore, Valid: true},
			FinancialRiskTier:    sql.NullString{String: sup.FinancialRiskTier, Valid: sup.FinancialRiskTier != ""},
			Certifications:       sql.NullString{String: fmt.Sprintf("%v", sup.Certifications), Valid: len(sup.Certifications) > 0},
			ComplianceFlags:      sql.NullString{String: fmt.Sprintf("%v", sup.ComplianceFlags), Valid: len(sup.ComplianceFlags) > 0},
			ApprovedStatus:       sql.NullString{String: sup.ApprovedStatus, Valid: sup.ApprovedStatus != ""},
			Contracts:            sql.NullString{String: fmt.Sprintf("%v", sup.Contracts), Valid: len(sup.Contracts) > 0},
			TermsVersion:         sql.NullString{String: sup.TermsVersion, Valid: sup.TermsVersion != ""},
			Lat:                  sql.NullFloat64{Float64: sup.GeoCoords.Lat, Valid: true},
			Lon:                  sql.NullFloat64{Float64: sup.GeoCoords.Lon, Valid: true},
			DataSource:           sql.NullString{String: sup.DataSource, Valid: sup.DataSource != ""},
			SourceTimestamp:      sql.NullTime{Time: sup.SourceTimestamp, Valid: true},
			IngestionTimestamp:   sql.NullTime{Time: sup.IngestionTimestamp, Valid: true},
			SchemaVersion:        sql.NullString{String: sup.SchemaVersion, Valid: sup.SchemaVersion != ""},
		})
		if err != nil {
			log.Fatal("failed to insert supplier:", sup.SupplierID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted suppliers:", len(sups))

	// 4. collect supplier IDs
	var supplierIDs []string
	for _, sup := range sups {
		supplierIDs = append(supplierIDs, sup.SupplierID)
	}

	// 5. generate parts
	partsList := parts.GenerateParts(10000, tenant, supplierIDs)
	parts.PartsWriter("parts.csv", partsList)

	// 6. insert parts in a transaction
	tx, err = conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	qtx = q.WithTx(tx)

	for _, part := range partsList {
		_, err := qtx.CreatePart(ctx, db.CreatePartParams{
			PartID:               part.PartID,
			TenantID:             part.TenantID,
			PartNumber:           part.PartNumber,
			Description:          part.Description,
			Category:             sql.NullString{String: part.Category, Valid: part.Category != ""},
			LifecycleStatus:      sql.NullString{String: part.LifecycleStatus, Valid: part.LifecycleStatus != ""},
			Uom:                  sql.NullString{String: part.Uom, Valid: part.Uom != ""},
			SpecHash:             sql.NullString{String: part.SpecHash, Valid: part.SpecHash != ""},
			BomCompatibility:     sql.NullString{String: fmt.Sprintf("%v", part.BomCompatibility), Valid: len(part.BomCompatibility) > 0},
			DefaultSupplierID:    sql.NullString{String: part.DefaultSupplierID, Valid: part.DefaultSupplierID != ""},
			QualifiedSupplierIds: sql.NullString{String: fmt.Sprintf("%v", part.QualifiedSupplierIDs), Valid: len(part.QualifiedSupplierIDs) > 0},
			UnitCost:             sql.NullFloat64{Float64: part.UnitCost, Valid: true},
			Moq:                  sql.NullInt64{Int64: int64(part.Moq), Valid: true},
			LeadTimeDaysAvg:      sql.NullInt64{Int64: int64(part.LeadTimeDaysAvg), Valid: true},
			LeadTimeDaysP95:      sql.NullInt64{Int64: int64(part.LeadTimeDaysP95), Valid: true},
			QualityGrade:         sql.NullString{String: part.QualityGrade, Valid: part.QualityGrade != ""},
			ComplianceFlags:      sql.NullString{String: fmt.Sprintf("%v", part.ComplianceFlags), Valid: len(part.ComplianceFlags) > 0},
			HazardClass:          sql.NullString{String: part.HazardClass, Valid: part.HazardClass != ""},
			LastPriceChange:      sql.NullTime{Time: part.LastPriceChange, Valid: true},
			DataSource:           sql.NullString{String: part.DataSource, Valid: part.DataSource != ""},
			SourceTimestamp:      sql.NullTime{Time: part.SourceTimestamp, Valid: true},
			IngestionTimestamp:   sql.NullTime{Time: part.IngestionTimestamp, Valid: true},
			SchemaVersion:        sql.NullString{String: part.SchemaVersion, Valid: part.SchemaVersion != ""},
		})
		if err != nil {
			log.Fatal("failed to insert part:", part.PartID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted parts:", len(partsList))
}
