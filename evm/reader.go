package evm

import (
	"context"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"math/big"
	"sort"
)

type BlockReadsInput struct {
	ChainId                      int64
	BlockNumber                  *big.Int
	ConcurrencyLimit             int
	Schemas                      []model.Schema
	Client                       *ethclient.Client
	IncludeBinaryTransactionData bool
}

type BlockReadInput struct {
	ChainId                      int64
	BlockNumber                  *big.Int
	Schemas                      []model.Schema
	Client                       *ethclient.Client
	IncludeBinaryTransactionData bool
}

func readBlocks(ctx context.Context, input BlockReadsInput) []model.Block {
	readCh := make(chan model.Block, input.ConcurrencyLimit)
	for i := 0; i < input.ConcurrencyLimit; i++ {
		blockNumber := big.NewInt(0)
		blockNumber = blockNumber.Add(input.BlockNumber, big.NewInt(int64(i)))
		go readBlock(ctx, BlockReadInput{
			ChainId:     input.ChainId,
			BlockNumber: blockNumber,
			Schemas:     input.Schemas,
			Client:      input.Client,
		}, readCh)
	}

	var blocks []model.Block
	for i := 0; i < input.ConcurrencyLimit; i++ {
		block := <-readCh
		blocks = append(blocks, block)
		log.Trace().Msgf("block read completed for block number: %s, chain id: %d, transactions count: %d", block.Number.String(), block.ChainId, len(block.Transactions))
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Number.Cmp(blocks[j].Number) == -1
	})

	return blocks
}

func readBlock(ctx context.Context, input BlockReadInput, readCh chan model.Block) {
	log.Trace().Msgf("reading block number: %s for chain id: %d", input.BlockNumber.String(), input.ChainId)

	client := input.Client
	block, _ := client.BlockByNumber(ctx, input.BlockNumber)
	blockTime := int64(block.Time())

	var transactions []model.Transaction
	for _, tx := range block.Transactions() {
		log.Trace().Msgf("reading transaction hash: %s for chain id: %d", tx.Hash().String(), input.ChainId)
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			if err.Error() == "not found" {
				log.Error().Err(err).Msgf("transaction receipt not found for chainId: %d, blockNumber: %s, txHash: %s", input.ChainId, input.BlockNumber.String(), tx.Hash().String())
				continue
			}

			log.Fatal().Err(err).Msgf("transaction receipt retrieval failed for chainId: %d, blockNumber: %s, txHash: %s", input.ChainId, input.BlockNumber.String(), tx.Hash().String())
			return
		}

		var logs []map[string]interface{}
		for _, vlog := range receipt.Logs {
			if len(vlog.Topics) == 0 {
				continue
			}

			logData, err := mapLog(input.ChainId, vlog, input.Schemas)
			if err != nil {
				log.Fatal().Err(err).Msgf("log mapping failed, txHash: %s, logIndex: %d", vlog.TxHash.String(), vlog.Index)
				return
			}

			if logData == nil {
				continue
			}

			log.Trace().Msgf("log alias: %s, contract address: %s for transaction hash: %s", logData[model.EventAliasLogField], logData[model.ContractAddressLogField], tx.Hash().String())
			logs = append(logs, logData)
		}

		to := ""
		if tx.To() != nil {
			to = tx.To().String()
		}

		transactions = append(transactions, model.Transaction{
			ContractAddress: receipt.ContractAddress.String(),
			To:              to,
			Index:           int64(receipt.TransactionIndex),
			Hash:            receipt.TxHash.String(),
			Value:           tx.Value(),
			Logs:            logs,
		})
	}

	readCh <- model.Block{
		ChainId:      input.ChainId,
		Number:       input.BlockNumber,
		Time:         blockTime,
		Hash:         block.Hash().String(),
		Transactions: transactions,
	}
}
