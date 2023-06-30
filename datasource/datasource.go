package datasource

import (
	"context"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"math/big"
	"strings"
)

const (
	sourcePostgres = "postgres"
)

type Source interface {
	init(ctx context.Context) error
	CreateChain(ctx context.Context, chain model.Chain) error
	GetChain(ctx context.Context, chainId int64) (model.Chain, error)
	GetChains(ctx context.Context) ([]model.Chain, error)
	UpdateChain(ctx context.Context, chain model.Chain) error
	UpdateBlock(ctx context.Context, block model.Block) error
	CreateSchema(ctx context.Context, schema model.Schema) error
	GetSchema(ctx context.Context, id string) (model.Schema, error)
	GetSchemas(ctx context.Context) ([]model.Schema, error)
	GetSchemasWithFilters(ctx context.Context) (schemas []model.Schema, err error)
	CreateSchemaFilter(ctx context.Context, filter model.SchemaFilter) (err error)
	GetOutOfSyncSchemaFilters(ctx context.Context) ([]model.SchemaFilter, error)
	GetSchemaFilter(ctx context.Context, schemaFilterId string) (filter model.SchemaFilter, err error)
	GetSchemaFiltersBySchemaId(ctx context.Context, schemaId string) ([]model.SchemaFilter, error)
	UpdateSchemaFilterSync(ctx context.Context, filterId string, sync model.SchemaFilterSync) error
	RawTransactionLogsPersisted() bool
	Close(ctx context.Context) error
	GetRawTransactionBlockNumbersForContract(ctx context.Context, chainId int64, contractAddress string) (blockNumbers []*big.Int, err error)
	GetRawTransactionsForContract(ctx context.Context, chainId int64, blockNumber *big.Int, contractAddress string) ([]model.RawTransactionLog, error)
}

type Factory interface {
	CreateDataSource(ctx context.Context) (Source, error)
}

func NewFactory(ctx context.Context, config configs.DataSourceConfig) (Factory, error) {
	var factory Factory
	switch strings.ToLower(config.Source) {
	case sourcePostgres:
		factory = newPostgresFactory(config.Postgres)
		break
	default:
		return nil, errors.InvalidInput("invalid datasource type: %s", config.Source)
	}

	source, err := factory.CreateDataSource(ctx)
	if err != nil {
		return nil, err
	}

	defer source.Close(ctx)
	err = source.init(ctx)
	return factory, err
}
