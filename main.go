package main

import (
	"fmt"
	"hugo/extend/isreflect"
)

type order struct {
	OrdId      int
	CustomerId struct{Sdf string}
}

func main() {
	o := order{
		OrdId:      456,
		CustomerId: struct{ Sdf string }{Sdf: "fdfdaas"},
	}
	m:=isreflect.ReadStruct(o)
	fmt.Println(m)
}