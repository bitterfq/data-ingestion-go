package suppliers

import (
	"os"
	"testing"
)

func TestGenerateSuppliers(t *testing.T) {
	tenant := "tenant_acme"
	num := 10
	sups := GenerateSuppliers(tenant, num)
	if len(sups) != num {
		t.Errorf("expected %d suppliers, got %d", num, len(sups))
	}
	for _, sup := range sups {
		if sup.TenantID != tenant {
			t.Errorf("expected tenant ID %s, got %s", tenant, sup.TenantID)
		}
		if sup.SupplierID == "" {
			t.Error("expected non-empty SupplierID")
		}
		if sup.LegalName == "" {
			t.Error("expected non-empty LegalName")
		}
	}
}

func TestWriteSuppliersToCSV(t *testing.T) {
	tenant := "tenant_acme"
	num := 5
	sups := GenerateSuppliers(tenant, num)
	filename := "test_suppliers.csv"
	success := SupplierWriter(filename, sups)
	if !success {
		t.Error("expected WriteSuppliersToCSV to return true")
	}
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", filename)
	}
	// Clean up
	os.Remove(filename)
}
