package snowflake

import (
	"testing"
)

func TestClient(t *testing.T) {
	// Just test that NewClient runs without panic
	db, err := NewClient()

	if err != nil {
		t.Fatal("[FAILED] snowflake connector failed: ", err)
	}

	t.Log("[SUCCESS] snowflake connector succeeded")

	defer db.Close()
}
