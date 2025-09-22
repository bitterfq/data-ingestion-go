package parts

import (
	"os"
	"testing"
)

func TestGenerateParts(t *testing.T) {
	tenant := "tenant_acme"
	supplierIDs := []string{"sup1", "sup2", "sup3"}
	num := 10
	parts := GenerateParts(num, tenant, supplierIDs)
	if len(parts) != num {
		t.Errorf("expected %d parts, got %d", num, len(parts))
	}
	for _, part := range parts {
		if part.TenantID != tenant {
			t.Errorf("expected tenant ID %s, got %s", tenant, part.TenantID)
		}
		if part.PartID == "" {
			t.Error("expected non-empty PartID")
		}
		if part.PartNumber == "" {
			t.Error("expected non-empty PartNumber")
		}
		if part.DefaultSupplierID == "" {
			t.Error("expected non-empty DefaultSupplierID")
		}
		found := false
		for _, sid := range supplierIDs {
			if part.DefaultSupplierID == sid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected DefaultSupplierID to be one of %v, got %s", supplierIDs, part.DefaultSupplierID)
		}
	}
}
func TestWritePartsToCSV(t *testing.T) {
	tenant := "tenant_acme"
	supplierIDs := []string{"sup1", "sup2", "sup3"}
	num := 5
	parts := GenerateParts(num, tenant, supplierIDs)
	filename := "test_parts.csv"
	success := PartsWriter(filename, parts)
	if !success {
		t.Error("expected PartsWriter to return true")
	}
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", filename)
	}
	// Clean up
	os.Remove(filename)
}
