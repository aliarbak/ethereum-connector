package schemamodel

import "github.com/aliarbak/ethereum-connector/model"

type SchemaOutput struct {
	Id     string              `json:"id"`
	Name   string              `json:"name"`
	Abi    []model.Abi         `json:"abi"`
	Events []model.SchemaEvent `json:"events"`
}
