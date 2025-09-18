// Package parts provides data structures and functions for generating and exporting synthetic part data.
package parts

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// Part represents a part entity with identity, description, supplier, cost, compliance, and metadata fields.
type Part struct {
	PartID               string
	TenantID             string
	PartNumber           string
	Description          string
	Category             string
	LifecycleStatus      string
	Uom                  string
	SpecHash             string
	BomCompatibility     []string
	DefaultSupplierID    string
	QualifiedSupplierIDs []string
	UnitCost             float64
	Moq                  int
	LeadTimeDaysAvg      int
	LeadTimeDaysP95      int
	QualityGrade         string
	ComplianceFlags      []string
	HazardClass          string
	LastPriceChange      time.Time
	DataSource           string
	SourceTimestamp      time.Time
	IngestionTimestamp   time.Time
	SchemaVersion        string
}

// GeneratePart creates and returns a single synthetic Part with example data.
// The data is randomly generated for testing or demo purposes.
func GeneratePart(tenant string, supplierIDs []string) Part {

	categories := []string{"ELECTRICAL", "MECHANICAL", "RAW_MATERIAL", "OTHER"}
	lifecycle_status := []string{"NEW", "ACTIVE", "NRND", "EOL"}
	uoms := []string{"EA", "KG", "M"}
	grades := []string{"A", "B", "C"}
	flags := []string{"ROHS", "REACH", "ITAR"}
	hazards := []string{"", "flammable", "toxic", "corrosive"}

	default_supplier_id := ""
	qualified_supplier_ids := []string{}

	if len(supplierIDs) > 0 {
		default_supplier_id = gofakeit.RandomString(supplierIDs)
		qualified_supplier_ids = append(qualified_supplier_ids, default_supplier_id)

		if len(supplierIDs) > 1 {
			qualified_supplier_ids = append(qualified_supplier_ids, gofakeit.RandomString(supplierIDs))
		}

	}

	return Part{
		PartID:               gofakeit.UUID(),
		TenantID:             tenant,
		PartNumber:           "P-" + gofakeit.Numerify("######"),
		Description:          gofakeit.Sentence(5),
		Category:             gofakeit.RandomString(categories),
		LifecycleStatus:      gofakeit.RandomString(lifecycle_status),
		Uom:                  gofakeit.RandomString(uoms),
		SpecHash:             gofakeit.UUID(),
		BomCompatibility:     []string{gofakeit.LetterN(3), gofakeit.LetterN(3)},
		DefaultSupplierID:    default_supplier_id,
		QualifiedSupplierIDs: qualified_supplier_ids,
		UnitCost:             gofakeit.Price(1, 1000), // Random price between 1 and 1000 -- keep it simple
		Moq:                  gofakeit.Number(1, 500),
		LeadTimeDaysAvg:      gofakeit.Number(2, 60),
		LeadTimeDaysP95:      gofakeit.Number(5, 90),
		QualityGrade:         gofakeit.RandomString(grades),

		// GenerateParts creates and returns a slice of synthetic Parts.
		// The number of parts generated is specified by count.
		ComplianceFlags:    flags,
		HazardClass:        gofakeit.RandomString(hazards),
		LastPriceChange:    time.Now(),
		DataSource:         "synthetic.v1",
		SourceTimestamp:    time.Now().Add(-time.Hour * time.Duration(gofakeit.Number(1, 72))),
		IngestionTimestamp: time.Now(),
		SchemaVersion:      "1.0.0",
	}
}

func GenerateParts(count int, tenant string, supplierIDs []string) []Part {
	parts := make([]Part, count)
	for i := 0; i < count; i++ {
		parts[i] = GeneratePart(tenant, supplierIDs)
	}
	return parts
}

// PartsWriter writes a slice of Part records to a CSV file with the given filename.
// Returns true if the file was written successfully, false otherwise.
func PartsWriter(filename string, parts []Part) bool {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("error creating file:", err)
		return false
	}
	defer file.Close()

	w := csv.NewWriter(file)

	// Header row
	header := []string{
		"part_id", "tenant_id", "part_number",
		"description", "category", "lifecycle_status",
		"uom", "spec_hash", "bom_compatibility",
		"default_supplier_id", "qualified_supplier_ids",
		"unit_cost", "moq",
		"lead_time_days_avg", "lead_time_days_p95",
		"quality_grade", "compliance_flags", "hazard_class",
		"last_price_change",
		"data_source", "source_timestamp", "ingestion_timestamp", "schema_version",
	}
	if err := w.Write(header); err != nil {
		fmt.Println("error writing header:", err)
		return false
	}

	// Data rows
	for _, part := range parts {
		row := []string{
			part.PartID,
			part.TenantID,
			part.PartNumber,
			part.Description,
			part.Category,
			part.LifecycleStatus,
			part.Uom,
			part.SpecHash,
			strings.Join(part.BomCompatibility, ";"),
			part.DefaultSupplierID,
			strings.Join(part.QualifiedSupplierIDs, ";"),
			fmt.Sprintf("%.2f", part.UnitCost),
			strconv.Itoa(part.Moq),
			strconv.Itoa(part.LeadTimeDaysAvg),
			strconv.Itoa(part.LeadTimeDaysP95),
			part.QualityGrade,
			strings.Join(part.ComplianceFlags, ";"),
			part.HazardClass,
			part.LastPriceChange.Format(time.RFC3339),
			part.DataSource,
			part.SourceTimestamp.Format(time.RFC3339),
			part.IngestionTimestamp.Format(time.RFC3339),
			part.SchemaVersion,
		}
		if err := w.Write(row); err != nil {
			fmt.Println("error writing row:", err)
			return false
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Println("error flushing data to file:", err)
		return false
	}

	fmt.Println("Parts csv file:", filename)
	return true

}
