package evm

import (
	"context"
	"github.com/aliarbak/ethereum-connector/datasource"
	"github.com/aliarbak/ethereum-connector/destination"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"math/big"
	"time"
)

func RunListeners(ctx context.Context, dataSourceFactory datasource.Factory, destinationFactory destination.Factory) error {
	dataSource, err := dataSourceFactory.CreateDataSource(ctx)
	if err != nil {
		return errors.From(err, "datasource creation failed for run listeners")
	}

	defer dataSource.Close(ctx)

	chains, err := dataSource.GetChains(ctx)
	if err != nil {
		return err
	}

	for _, chain := range chains {
		err = RunListener(ctx, chain.Id, dataSourceFactory, destinationFactory)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunListener(ctx context.Context, chainId int64, dataSourceFactory datasource.Factory, destinationFactory destination.Factory) error {
	dataSource, err := dataSourceFactory.CreateDataSource(ctx)
	if err != nil {
		return errors.From(err, "datasource creation failed for listener")
	}

	dest, err := destinationFactory.CreateDestination(ctx)
	if err != nil {
		return errors.From(err, "destination creation failed for listener")
	}

	go func() {
		err := runListener(ctx, chainId, dataSource, dest)
		defer dataSource.Close(ctx)
		defer dest.Close(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("listener crashed")
		}
	}()

	log.Info().Msgf("evm listener is running for chainId: %d", chainId)
	return nil
}

func runListener(ctx context.Context, chainId int64, source datasource.Source, dest destination.Destination) error {
	chain, err := source.GetChain(ctx, chainId)
	if err != nil {
		return errors.From(err, "chain read from datasource failed for chainId: %d", chainId)
	}

	rpcUrl := chain.RpcUrl
	paused := chain.Paused
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return errors.From(err, "ethClient connection failed")
	}
	defer client.Close()

	currentBlockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return errors.From(err, "current block number read from chain failed")
	}

	for {
		chain, err = source.GetChain(ctx, chainId)
		blockNumber := chain.BlockNumber
		if err != nil {
			return errors.From(err, "chain read from datasource failed")
		}

		if chain.RpcUrl != rpcUrl {
			log.Info().Msgf("chain rpc client url changed for chainId: %d, rpcUrl: %s", chain.Id, chain.RpcUrl)
			client.Close()
			client, err = ethclient.Dial(rpcUrl)
			if err != nil {
				return errors.From(err, "ethClient connection failed with new rpc url. old rpc url: %s, new rpc url: %s", rpcUrl, chain.RpcUrl)
			}
			rpcUrl = chain.RpcUrl
		}

		if chain.Paused != paused {
			log.Info().Msgf("chain pause state changed for chainId: %d, paused: %t", chain.Id, chain.Paused)
			paused = chain.Paused
		}

		if paused {
			time.Sleep(5 * time.Second)
			continue
		}

		blockNumber = blockNumber.Add(blockNumber, big.NewInt(1))
		if blockNumber.Int64() > int64(currentBlockNumber) {
			currentBlockNumber, err = client.BlockNumber(ctx)
			if err != nil {
				return errors.From(err, "current block number read from chain failed")
			}

			log.Trace().Msgf("waiting for the block number: %s", blockNumber.String())
			time.Sleep(time.Millisecond * time.Duration(chain.BlockNumberWaitPeriodInMS))
			continue
		}

		schemas, err := source.GetSchemasWithFilters(ctx)
		if err != nil {
			return errors.From(err, "schemas read from datasource failed")
		}

		blocks := readBlocks(ctx, BlockReadsInput{
			ChainId:          chainId,
			BlockNumber:      blockNumber,
			ConcurrencyLimit: chain.ConcurrentBlockReadLimit,
			Schemas:          schemas,
			Client:           client,
		})

		if err = sendBlocks(ctx, blocks, source, dest); err != nil {
			return err
		}

		_, blockNumMod := new(big.Int).DivMod(blockNumber, big.NewInt(100), new(big.Int))
		if blockNumMod.Cmp(big.NewInt(int64(chain.ConcurrentBlockReadLimit))) <= 0 {
			log.Info().Msgf("blockNumber: %s for chainId: %d", blockNumber.String(), chainId)
		}

		time.Sleep(time.Millisecond * time.Duration(chain.BlockNumberReadPeriodInMS))
	}
}

func sendBlocks(ctx context.Context, blocks []model.Block, source datasource.Source, dest destination.Destination) (err error) {
	for _, block := range blocks {
		if dest.DeliveryGuarantee() == destination.AtMostOnceDeliveryGuarantee {
			if err = source.UpdateBlock(ctx, block); err != nil {
				return errors.From(err, "block update failed for chainId: %d, blockNumber: %s", block.ChainId, block.Number.String())
			}
		}

		if err = dest.SendBlock(ctx, block); err != nil {
			return errors.From(err, "block send failed for chainId: %d, blockNumber: %s", block.ChainId, block.Number.String())
		}

		if dest.DeliveryGuarantee() != destination.AtMostOnceDeliveryGuarantee {
			if err = source.UpdateBlock(ctx, block); err != nil {
				return errors.From(err, "block update failed for chainId: %d, blockNumber: %s", block.ChainId, block.Number.String())
			}
		}
	}

	return nil
}
