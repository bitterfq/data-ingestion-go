package suppliers

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/oklog/ulid/v2"
)

type Supplier struct {
	SupplierID           string
	TenantID             string
	SupplierCode         string
	LegalName            string
	Country              string
	LeadTimeDaysAvg      int
	LeadTimeDaysP95      int
	OnTimeDeliveryRate   float64
	DefectRatePPM        int
	CapacityUnitsPerWeek int
	RiskScore            float64
	FinancialRiskTier    string
	ApprovedStatus       string
	DataSource           string
}

func GenerateSupplier(tenant string) Supplier {

	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	onTime := gofakeit.Float64Range(60, 100)
	risk := 100 - onTime + gofakeit.Float64Range(0, 10)

	return Supplier{
		SupplierID:           ulid.MustNew(ulid.Timestamp(t), entropy).String(),
		TenantID:             tenant,
		SupplierCode:         gofakeit.LetterN(6), // ERP-like code
		LegalName:            gofakeit.Company(),
		Country:              gofakeit.CountryAbr(),
		LeadTimeDaysAvg:      gofakeit.Number(3, 90),
		LeadTimeDaysP95:      gofakeit.Number(7, 180),
		OnTimeDeliveryRate:   gofakeit.Float64Range(60, 100),
		DefectRatePPM:        gofakeit.Number(50, 1000),
		CapacityUnitsPerWeek: gofakeit.Number(100, 10000),
		RiskScore:            risk,
		FinancialRiskTier:    gofakeit.RandomString([]string{"LOW", "MEDIUM", "HIGH"}),
		ApprovedStatus:       gofakeit.RandomString([]string{"APPROVED", "PENDING", "SUSPENDED"}),
		DataSource:           "synthetic.v1",
	}

}

func GenerateSuppliers(tenant string, count int) []Supplier {
	suppliers := make([]Supplier, count)
	for i := 0; i < count; i++ {
		suppliers[i] = GenerateSupplier(tenant)
	}
	return suppliers
}
