package main

import (
	"fmt"

	"github.com/bitterfq/data-ingestion-go/suppliers"
)

func main() {
	s := suppliers.GenerateSuppliers("tenant_acme", 5000)
	fmt.Println(len(s))
}
