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

type Part struct {
	part_id                string
	tenant_id              string
	part_number            string
	description            string
	category               string
	lifecycle_status       string
	uom                    string
	spec_hash              string
	bom_compatibility      []string
	default_supplier_id    string
	qualified_supplier_ids []string
	unit_cost              float64
	moq                    int
	lead_time_days_avg     int
	lead_time_days_p95     int
	quality_grade          string
	compliance_flags       []string
	hazard_class           string
	last_price_change      time.Time
	data_source            string
	source_timestamp       time.Time
	ingestion_timestamp    time.Time
	schema_version         string
}

func GeneratePart() Part {
	return Part{
		part_id:                gofakeit.UUID(),
		tenant_id:              "example_tenant_id",
		part_number:            "example_part_number",
		description:            "example_description",
		category:               "example_category",
		lifecycle_status:       "active",
		uom:                    "pcs",
		spec_hash:              "example_spec_hash",
		bom_compatibility:      []string{"compatibility_1", "compatibility_2"},
		default_supplier_id:    "example_supplier_id",
		qualified_supplier_ids: []string{"supplier_1", "supplier_2"},
		unit_cost:              10.5,
		moq:                    100,
		lead_time_days_avg:     30,
		lead_time_days_p95:     45,
		quality_grade:          "A",
		compliance_flags:       []string{"flag_1", "flag_2"},
		hazard_class:           "non-hazardous",
		last_price_change:      time.Now(),
		data_source:            "example_source",
		source_timestamp:       time.Now(),
		ingestion_timestamp:    time.Now(),
		schema_version:         "1.0.0",
	}
}

func GenerateParts(count int) []Part {
	parts := make([]Part, count)
	for i := 0; i < count; i++ {
		parts[i] = GeneratePart()
	}

	return parts
}

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
			part.part_id,
			part.tenant_id,
			part.part_number,
			part.description,
			part.category,
			part.lifecycle_status,
			part.uom,
			part.spec_hash,
			strings.Join(part.bom_compatibility, ";"),
			part.default_supplier_id,
			strings.Join(part.qualified_supplier_ids, ";"),
			fmt.Sprintf("%.2f", part.unit_cost),
			strconv.Itoa(part.moq),
			strconv.Itoa(part.lead_time_days_avg),
			strconv.Itoa(part.lead_time_days_p95),
			part.quality_grade,
			strings.Join(part.compliance_flags, ";"),
			part.hazard_class,
			part.last_price_change.Format(time.RFC3339),
			part.data_source,
			part.source_timestamp.Format(time.RFC3339),
			part.ingestion_timestamp.Format(time.RFC3339),
			part.schema_version,
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

	return true

}
