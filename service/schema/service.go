package schema

import (
	"context"
	"github.com/aliarbak/ethereum-connector/datasource"
	"github.com/aliarbak/ethereum-connector/destination"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/evm"
	"github.com/aliarbak/ethereum-connector/model"
	schemamodel "github.com/aliarbak/ethereum-connector/service/schema/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"math/big"
)

type Service interface {
	CreateSchema(ctx context.Context, input schemamodel.CreateSchemaInput) (schemamodel.SchemaOutput, error)
	GetSchemas(ctx context.Context) (output []schemamodel.SchemaOutput, err error)
	GetSchemaById(ctx context.Context, id string) (schemamodel.SchemaOutput, error)
	CreateSchemaFilter(ctx context.Context, schemaId string, input schemamodel.CreateSchemaFilterInput) (filter model.SchemaFilter, err error)
	GetSchemaFilters(ctx context.Context, schemaId string) ([]model.SchemaFilter, error)
	GetSchemaFilterById(ctx context.Context, schemaId string, schemaFilterId string) (model.SchemaFilter, error)
	Close(ctx context.Context) error
}

type service struct {
	dataSource         datasource.Source
	dataSourceFactory  datasource.Factory
	destinationFactory destination.Factory
}

func NewService(ctx context.Context, dataSourceFactory datasource.Factory, destinationFactory destination.Factory) (s Service, err error) {
	dataSource, err := dataSourceFactory.CreateDataSource(ctx)
	if err != nil {
		return
	}

	return &service{
		dataSource,
		dataSourceFactory,
		destinationFactory,
	}, err
}

func (s service) CreateSchema(ctx context.Context, input schemamodel.CreateSchemaInput) (output schemamodel.SchemaOutput, err error) {
	id := uuid.New().String()
	err = s.dataSource.CreateSchema(ctx, model.Schema{
		Id:     id,
		Name:   input.Name,
		Abi:    input.Abi,
		Events: input.Events,
	})

	if err != nil {
		return output, err
	}

	schema, err := s.dataSource.GetSchema(ctx, id)
	if err != nil {
		return output, err
	}

	return schemamodel.SchemaOutput{
		Id:     schema.Id,
		Name:   schema.Name,
		Abi:    schema.Abi,
		Events: schema.Events,
	}, nil
}

func (s service) GetSchemaById(ctx context.Context, id string) (output schemamodel.SchemaOutput, err error) {
	schema, err := s.dataSource.GetSchema(ctx, id)
	if err != nil {
		return output, err
	}

	return schemamodel.SchemaOutput{
		Id:     schema.Id,
		Name:   schema.Name,
		Abi:    schema.Abi,
		Events: schema.Events,
	}, nil
}

func (s service) GetSchemas(ctx context.Context) (output []schemamodel.SchemaOutput, err error) {
	schemas, err := s.dataSource.GetSchemas(ctx)
	if err != nil {
		return output, err
	}

	for _, schema := range schemas {
		output = append(output, schemamodel.SchemaOutput{
			Id:     schema.Id,
			Name:   schema.Name,
			Abi:    schema.Abi,
			Events: schema.Events,
		})
	}

	return output, err
}

func (s service) CreateSchemaFilter(ctx context.Context, schemaId string, input schemamodel.CreateSchemaFilterInput) (filter model.SchemaFilter, err error) {
	chain, err := s.dataSource.GetChain(ctx, input.ChainId)
	if err != nil {
		return filter, err
	}

	schema, err := s.dataSource.GetSchema(ctx, schemaId)
	if err != nil {
		return filter, err
	}

	if input.Sync && !s.dataSource.RawTransactionLogsPersisted() {
		return filter, errors.InvalidInput("you can not sync contract events since raw transaction logs are not persisted in the datasource")
	}

	var sync *model.SchemaFilterSync
	if input.Sync {
		sync = &model.SchemaFilterSync{
			FromBlockNumber:    0,
			ToBlockNumber:      new(big.Int).Add(chain.BlockNumber, big.NewInt(1)).Int64(),
			CurrentBlockNumber: 0,
			Completed:          false,
		}
	}

	id := uuid.New().String()
	err = s.dataSource.CreateSchemaFilter(ctx, model.SchemaFilter{
		Id:              id,
		ChainId:         chain.Id,
		SchemaId:        schema.Id,
		ContractAddress: common.HexToAddress(input.ContractAddress).String(),
		Sync:            sync,
	})

	if err != nil {
		return filter, err
	}

	filter, err = s.dataSource.GetSchemaFilter(ctx, id)
	if err != nil {
		return filter, err
	}

	if input.Sync {
		err = evm.RunContractSyncWorker(context.Background(), id, s.dataSourceFactory, s.destinationFactory)
		if err != nil {
			return filter, errors.From(err, "schema filter has been created but worker initialization failed")
		}
	}

	return filter, err
}

func (s service) GetSchemaFilters(ctx context.Context, schemaId string) ([]model.SchemaFilter, error) {
	return s.dataSource.GetSchemaFiltersBySchemaId(ctx, schemaId)
}

func (s service) GetSchemaFilterById(ctx context.Context, schemaId string, schemaFilterId string) (model.SchemaFilter, error) {
	filter, err := s.dataSource.GetSchemaFilter(ctx, schemaFilterId)
	if err != nil {
		return filter, err
	}

	if filter.SchemaId != schemaId {
		return filter, errors.InvalidInput("schema filter does not belong to given schema id")
	}

	return filter, nil
}

func (s service) Close(ctx context.Context) error {
	return s.dataSource.Close(ctx)
}
