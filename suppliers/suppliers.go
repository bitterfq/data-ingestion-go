package suppliers

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/oklog/ulid/v2"
)

type GeoCoords struct {
	Lat float64
	Lon float64
}

type Supplier struct {
	// Identity
	SupplierID   string
	TenantID     string
	SupplierCode string

	// Names & location
	LegalName string
	DBAName   string
	Country   string
	Region    string

	// Address
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string

	// Contacts
	ContactEmail string
	ContactPhone string

	// Commercial
	PreferredCurrency string
	Incoterms         string

	// Performance & risk
	LeadTimeDaysAvg      int
	LeadTimeDaysP95      int
	OnTimeDeliveryRate   float64
	DefectRatePPM        int
	CapacityUnitsPerWeek int
	RiskScore            float64
	FinancialRiskTier    string

	// Certifications & compliance
	Certifications  []string
	ComplianceFlags []string

	// Status & contracts
	ApprovedStatus string
	Contracts      []string
	TermsVersion   string

	// Geo
	GeoCoords *GeoCoords

	// Lineage / metadata
	DataSource         string
	SourceTimestamp    time.Time
	IngestionTimestamp time.Time
	SchemaVersion      string
}

func GenerateSupplier(tenant string) Supplier {

	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	onTime := gofakeit.Float64Range(60, 100)
	risk := 100 - onTime + gofakeit.Float64Range(0, 10)

	// Uniform stubs for now
	countries := []string{"US", "CN", "DE", "MX", "IN", "VN", "PL", "JP", "KR"}
	regions := []string{"EMEA", "APAC", "AMERICAS"}
	incoterms := []string{"DDP", "FOB", "CIF", "EXW"}
	tiers := []string{"LOW", "MEDIUM", "HIGH"}
	statuses := []string{"APPROVED", "PENDING", "SUSPENDED"}
	currencies := []string{"USD", "CNY", "EUR", "INR", "JPY"}
	certs := []string{"ISO9001", "IATF16949", "AS9100", "ISO14001"}
	flags := []string{"ITAR", "REACH", "ROHS"}

	return Supplier{
		// Identity
		SupplierID:   ulid.MustNew(ulid.Timestamp(t), entropy).String(),
		TenantID:     tenant,
		SupplierCode: gofakeit.LetterN(1) + gofakeit.Numerify("######"),

		// Names & location
		LegalName: gofakeit.Company(),
		DBAName:   gofakeit.CompanySuffix(),
		Country:   gofakeit.RandomString(countries),
		Region:    gofakeit.RandomString(regions),

		// Address
		AddressLine1: gofakeit.Street(),
		AddressLine2: "",
		City:         gofakeit.City(),
		State:        gofakeit.StateAbr(),
		PostalCode:   gofakeit.Zip(),

		// Contacts
		ContactEmail: gofakeit.Email(),
		ContactPhone: gofakeit.Phone(),

		// Commercial
		PreferredCurrency: gofakeit.RandomString(currencies),
		Incoterms:         gofakeit.RandomString(incoterms),

		// Performance & risk
		LeadTimeDaysAvg:      gofakeit.Number(3, 90),
		LeadTimeDaysP95:      gofakeit.Number(7, 180),
		OnTimeDeliveryRate:   onTime,
		DefectRatePPM:        gofakeit.Number(50, 1000),
		CapacityUnitsPerWeek: gofakeit.Number(100, 10000),
		RiskScore:            risk,
		FinancialRiskTier:    gofakeit.RandomString(tiers),

		// Certifications & compliance
		Certifications:  []string{gofakeit.RandomString(certs)},
		ComplianceFlags: []string{gofakeit.RandomString(flags)},

		// Status & contracts
		ApprovedStatus: gofakeit.RandomString(statuses),
		Contracts:      []string{"CONTRACT_" + gofakeit.Numerify("####")},
		TermsVersion:   gofakeit.Numerify("#.#"),

		// Geo
		GeoCoords: &GeoCoords{
			Lat: gofakeit.Latitude(),
			Lon: gofakeit.Longitude(),
		},

		// Lineage / metadata
		DataSource:         "synthetic.v1",
		SourceTimestamp:    time.Now().Add(-time.Hour * time.Duration(gofakeit.Number(1, 72))),
		IngestionTimestamp: time.Now(),
		SchemaVersion:      "1.0.0",
	}

}

func GenerateSuppliers(tenant string, count int) []Supplier {
	suppliers := make([]Supplier, count)
	for i := 0; i < count; i++ {
		suppliers[i] = GenerateSupplier(tenant)
	}
	return suppliers
}

func SupplierWriter(filename string, suppliers []Supplier) bool {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("error creating file:", err)
		return false
	}
	defer file.Close()

	w := csv.NewWriter(file)

	// Header row
	header := []string{
		"supplier_id", "tenant_id", "supplier_code",
		"legal_name", "dba_name", "country", "region",
		"address_line1", "address_line2", "city", "state", "postal_code",
		"contact_email", "contact_phone",
		"preferred_currency", "incoterms",
		"lead_time_days_avg", "lead_time_days_p95", "on_time_delivery_rate",
		"defect_rate_ppm", "capacity_units_per_week", "risk_score", "financial_risk_tier",
		"certifications", "compliance_flags",
		"approved_status", "contracts", "terms_version",
		"lat", "lon",
		"data_source", "source_timestamp", "ingestion_timestamp", "schema_version",
	}
	if err := w.Write(header); err != nil {
		fmt.Println("error writing header:", err)
		return false
	}

	// Data rows
	for _, sup := range suppliers {
		row := []string{
			sup.SupplierID,
			sup.TenantID,
			sup.SupplierCode,
			sup.LegalName,
			sup.DBAName,
			sup.Country,
			sup.Region,
			sup.AddressLine1,
			sup.AddressLine2,
			sup.City,
			sup.State,
			sup.PostalCode,
			sup.ContactEmail,
			sup.ContactPhone,
			sup.PreferredCurrency,
			sup.Incoterms,
			fmt.Sprintf("%d", sup.LeadTimeDaysAvg),
			fmt.Sprintf("%d", sup.LeadTimeDaysP95),
			fmt.Sprintf("%.2f", sup.OnTimeDeliveryRate),
			fmt.Sprintf("%d", sup.DefectRatePPM),
			fmt.Sprintf("%d", sup.CapacityUnitsPerWeek),
			fmt.Sprintf("%.2f", sup.RiskScore),
			sup.FinancialRiskTier,
			fmt.Sprintf("%v", sup.Certifications),
			fmt.Sprintf("%v", sup.ComplianceFlags),
			sup.ApprovedStatus,
			fmt.Sprintf("%v", sup.Contracts),
			sup.TermsVersion,
		}
		if sup.GeoCoords != nil {
			row = append(row,
				fmt.Sprintf("%.6f", sup.GeoCoords.Lat),
				fmt.Sprintf("%.6f", sup.GeoCoords.Lon),
			)
		} else {
			row = append(row, "", "")
		}
		row = append(row,
			sup.DataSource,
			sup.SourceTimestamp.Format(time.RFC3339),
			sup.IngestionTimestamp.Format(time.RFC3339),
			sup.SchemaVersion,
		)

		if err := w.Write(row); err != nil {
			fmt.Println("error writing row:", err)
			return false
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Println("error flushing writer:", err)
		return false
	}

	fmt.Println("wrote file:", filename)
	return true
}
