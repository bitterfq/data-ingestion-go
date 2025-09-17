package main

import (
	"fmt"

	"github.com/bitterfq/data-ingestion-go/suppliers"
)

func main() {
	s := suppliers.GenerateSuppliers("tenant_acme", 50)
	fmt.Println("%+v\n", s)
}
