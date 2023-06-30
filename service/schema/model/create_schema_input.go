package schemamodel

import (
	"github.com/aliarbak/ethereum-connector/model"
)

type CreateSchemaInput struct {
	Name   string              `json:"name" validate:"required"`
	Abi    []model.Abi         `json:"abi" validate:"required"`
	Events []model.SchemaEvent `json:"events" validate:"required"`
}
