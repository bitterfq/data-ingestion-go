package main

import (
	"fmt"

	"github.com/bitterfq/data-ingestion-go/parts"
	"github.com/bitterfq/data-ingestion-go/suppliers"
)

func main() {
	s := suppliers.GenerateSuppliers("tenant_acme", 100000)
	p := parts.GenerateParts(100000)
	//fmt.Println("%+v\n", s)

	out := suppliers.SupplierWriter("data/suppliers.csv", s)
	fmt.Println("Suppliers generated:", len(s))
	out2 := parts.PartsWriter("data/parts.csv", p)
	fmt.Println("Parts generated:", len(p))

	if !out || !out2 {
		panic("error writing csv files")
	}
}
