package schemamodel

type CreateSchemaFilterInput struct {
	ChainId         int64  `json:"chainId" validate:"required"`
	ContractAddress string `json:"contractAddress" validate:"required"`
	Sync            bool   `json:"sync"`
}
