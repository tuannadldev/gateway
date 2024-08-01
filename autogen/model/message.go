package model

type Message struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name        string // REST param name must match gRPC message field name
	Type        string
	Bit         int64  // number of bit for int types (8, 16, 32, 64)
	Location    string // used for GET request, allowed values: Param, Query
	Required    bool
	Repeated    bool
	Description string
}
