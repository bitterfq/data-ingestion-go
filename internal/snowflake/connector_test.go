package snowflake

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	// Just test that NewClient runs without panic
	db, err := NewClient("../../.env")

	if err != nil {
		t.Fatal("[FAILED] snowflake connector failed: ", err)
	}

	t.Log("[SUCCESS] snowflake connector succeeded")

	defer db.Close()
}

func TestClientHealthCheck(t *testing.T) {
	db, err := NewClient("../../.env")
	if err != nil {
		t.Fatal("Failed to create client:", err)
	}
	defer db.Close()

	rows, _ := db.Query("SELECT CURRENT_USER(), CURRENT_ROLE(), CURRENT_DATABASE(), CURRENT_SCHEMA(), CURRENT_WAREHOUSE()")
	defer rows.Close()
	for rows.Next() {
		var user, role, dbname, schema, wh string
		rows.Scan(&user, &role, &dbname, &schema, &wh)
		fmt.Println("User:", user, "Role:", role, "DB:", dbname, "Schema:", schema, "WH:", wh)
		fmt.Println("-----")
	}
	t.Log("Snowflake client test succeeded")
}

func TestSimpleQuery(t *testing.T) {
	db, err := NewClient("../../.env")
	if err != nil {
		t.Fatal("Failed to create client:", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM SUPPLY_CHAIN.PUBLIC.SUPPLIERS").Scan(&count)
	if err != nil {
		t.Fatal("Failed to query tables:", err)
	}
	fmt.Println("Total tables visible to this role:\n", count)
	fmt.Println("-----")
	t.Log("Simple query test succeeded")

}

func TestSelectQuery(t *testing.T) {
	db, err := NewClient("../../.env")
	if err != nil {
		t.Fatal("Failed to create client:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT SUPPLIER_ID, TENANT_ID, SUPPLIER_CODE FROM SUPPLY_CHAIN.PUBLIC.SUPPLIERS LIMIT 5")
	if err != nil {
		t.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var tenantId string
		var supplierCode string
		if err := rows.Scan(&id, &tenantId, &supplierCode); err != nil {
			t.Fatal("Failed to scan row:", err)
		}
		fmt.Println("Supplier ID:", id, "Tenant ID:", tenantId, "Supplier Code:", supplierCode)
		fmt.Println("-----")
	}
	t.Log("Select query test succeeded")
}

func TestLoadPrivateKey(t *testing.T) {
	// Adjust the path as necessary for your environment
	keyPath := "keys/rsa_key.p8"
	passphrase := []byte("") // Empty passphrase for this example

	privateKey, err := loadPrivateKey(keyPath, passphrase)
	if err != nil {
		t.Fatal("Failed to load private key:", err)
	}
	if privateKey == nil {
		t.Fatal("Private key is nil")
	}
	t.Log("Private key loaded successfully")
}
