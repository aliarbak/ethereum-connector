package chain

import (
	"context"
	"github.com/aliarbak/ethereum-connector/datasource"
	"github.com/aliarbak/ethereum-connector/destination"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/evm"
	"github.com/aliarbak/ethereum-connector/model"
	chainmodel "github.com/aliarbak/ethereum-connector/service/chain/model"
	"math/big"
	"time"
)

type Service interface {
	CreateChain(ctx context.Context, input chainmodel.CreateChainInput) (model.Chain, error)
	GetChains(ctx context.Context) ([]model.Chain, error)
	GetChainById(ctx context.Context, id int64) (model.Chain, error)
	PatchChain(ctx context.Context, id int64, input chainmodel.PatchChainInput) (chain model.Chain, err error)
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

func (s service) CreateChain(ctx context.Context, input chainmodel.CreateChainInput) (chain model.Chain, err error) {
	blockNumber := int64(-1)
	if input.BlockNumber != nil {
		blockNumber = *input.BlockNumber
	}

	blockNumberWaitPeriodInMs := int64(2000)
	if input.BlockNumberWaitPeriodInMS != nil {
		blockNumberWaitPeriodInMs = *input.BlockNumberWaitPeriodInMS
	}

	blockNumberReadPeriodInMs := int64(500)
	if input.BlockNumberReadPeriodInMS != nil {
		blockNumberReadPeriodInMs = *input.BlockNumberReadPeriodInMS
	}

	concurrentBlockReadLimit := 1
	if input.ConcurrentBlockReadLimit != nil {
		concurrentBlockReadLimit = *input.ConcurrentBlockReadLimit
	}

	err = s.dataSource.CreateChain(ctx, model.Chain{
		Id:                        input.Id,
		BlockNumber:               big.NewInt(blockNumber),
		Paused:                    input.Paused,
		UpdatedAt:                 time.Now().UTC().UnixMilli(),
		RpcUrl:                    input.RpcUrl,
		BlockNumberWaitPeriodInMS: blockNumberWaitPeriodInMs,
		BlockNumberReadPeriodInMS: blockNumberReadPeriodInMs,
		ConcurrentBlockReadLimit:  concurrentBlockReadLimit,
	})

	if err != nil {
		return
	}

	chain, err = s.dataSource.GetChain(ctx, input.Id)
	if err != nil {
		return
	}

	err = evm.RunListener(context.Background(), chain.Id, s.dataSourceFactory, s.destinationFactory)
	if err != nil {
		return chain, errors.From(err, "chain has been created but listener initialization failed")
	}

	return
}

func (s service) GetChains(ctx context.Context) ([]model.Chain, error) {
	return s.dataSource.GetChains(ctx)
}

func (s service) GetChainById(ctx context.Context, id int64) (model.Chain, error) {
	return s.dataSource.GetChain(ctx, id)
}

func (s service) PatchChain(ctx context.Context, id int64, input chainmodel.PatchChainInput) (chain model.Chain, err error) {
	chain, err = s.dataSource.GetChain(ctx, id)
	if err != nil {
		return
	}

	chain.UpdatedAt = time.Now().UTC().UnixMilli()
	if input.RpcUrl != nil {
		chain.RpcUrl = *input.RpcUrl
	}

	if input.Paused != nil {
		chain.Paused = *input.Paused
	}

	if input.ConcurrentBlockReadLimit != nil {
		chain.ConcurrentBlockReadLimit = *input.ConcurrentBlockReadLimit
	}

	if input.BlockNumberReadPeriodInMS != nil {
		chain.BlockNumberReadPeriodInMS = *input.BlockNumberReadPeriodInMS
	}

	if input.BlockNumberWaitPeriodInMS != nil {
		chain.BlockNumberWaitPeriodInMS = *input.BlockNumberWaitPeriodInMS
	}

	err = s.dataSource.UpdateChain(ctx, chain)
	if err != nil {
		return
	}

	return s.dataSource.GetChain(ctx, id)
}

func (s service) Close(ctx context.Context) error {
	return s.dataSource.Close(ctx)
}
