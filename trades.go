package main

type Data struct {
	Symbol     string   `json:"s"`
	Price      float64  `json:"p"`
	Timestamp  uint     `json:"t"`
	Conditions []string `json:"c" omitempty:"true"`
}

type Response struct {
	Type string `json:"type"`
	Data []Data `json:"data"`
}

type Subscribe struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}
