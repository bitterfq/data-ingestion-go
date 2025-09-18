package main

import (
	"github.com/bitterfq/data-ingestion-go/suppliers"
)

func main() {
	s := suppliers.GenerateSuppliers("tenant_acme", 5000)
	//fmt.Println("%+v\n", s)

	out := suppliers.SupplierWriter("data/suppliers.csv", s)
	print(out)
}
