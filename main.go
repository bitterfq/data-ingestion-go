package main

import (
	"github.com/bitterfq/data-ingestion-go/parts"
	"github.com/bitterfq/data-ingestion-go/suppliers"
)

func main() {
	s := suppliers.GenerateSuppliers("tenant_acme", 10)
	p := parts.GenerateParts(10)
	//fmt.Println("%+v\n", s)

	out := suppliers.SupplierWriter("data/suppliers.csv", s)
	out2 := parts.PartsWriter("data/parts.csv", p)

	//fmt.Println("%+v\n", p)
	print(out)
	print(out2)
}
