package datasource

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type redisDataSourceFactory struct {
	config configs.RedisDataSourceConfig
}

type redisDataSource struct {
	client *redis.Client
}

func newRedisFactory(config configs.RedisDataSourceConfig) Factory {
	return &redisDataSourceFactory{
		config: config,
	}
}

func (f redisDataSourceFactory) CreateDataSource(context.Context) (source Source, err error) {
	if f.config.PersistRawTransactionLogs {
		return source, errors.New("redis data source does not support persisting raw transaction logs")
	}

	hosts := strings.Split(f.config.Hosts, ",")
	var client *redis.Client
	if len(f.config.ServiceName) > 0 {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: hosts,
			MasterName:    f.config.ServiceName,
			Username:      f.config.Username,
			Password:      f.config.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     hosts[0],
			Username: f.config.Username,
			Password: f.config.Password,
		})
	}

	return &redisDataSource{client: client}, err
}

func (r redisDataSource) CreateChain(ctx context.Context, chain model.Chain) (err error) {
	_, err = r.getChain(ctx, chain.Id)
	if err == nil {
		return errors.AlreadyExists("chain already exists with given id: %d", chain.Id)
	}

	if err != redis.Nil {
		return err
	}

	result := r.client.Set(ctx, getChainKey(chain.Id), &ChainDocument{
		Id:                        chain.Id,
		BlockNumber:               chain.BlockNumber,
		Paused:                    chain.Paused,
		RpcUrl:                    chain.RpcUrl,
		BlockNumberWaitPeriodInMS: chain.BlockNumberWaitPeriodInMS,
		BlockNumberReadPeriodInMS: chain.BlockNumberReadPeriodInMS,
		ConcurrentBlockReadLimit:  chain.ConcurrentBlockReadLimit,
		UpdatedAt:                 chain.UpdatedAt,
	}, time.Duration(0))

	err = result.Err()
	if err != nil {
		return
	}

	chainIds, err := r.getChainIds(ctx)
	if err != nil {
		return err
	}

	chainIds = append(chainIds, chain.Id)
	result = r.client.Set(ctx, getChainsKey(), chainIds, time.Duration(0))
	return result.Err()
}

func (r *redisDataSource) UpdateBlock(ctx context.Context, block model.Block) error {
	chain, err := r.getChain(ctx, block.ChainId)
	if err != nil {
		return err
	}

	result := r.client.Set(ctx, getChainKey(block.ChainId), ChainDocument{
		Id:                        chain.Id,
		BlockNumber:               block.Number,
		Paused:                    chain.Paused,
		RpcUrl:                    chain.RpcUrl,
		BlockNumberWaitPeriodInMS: chain.BlockNumberWaitPeriodInMS,
		BlockNumberReadPeriodInMS: chain.BlockNumberReadPeriodInMS,
		ConcurrentBlockReadLimit:  chain.ConcurrentBlockReadLimit,
		UpdatedAt:                 time.Now().UTC().UnixMilli(),
	}, time.Duration(0))

	return result.Err()
}

func (r *redisDataSource) UpdateChain(ctx context.Context, chain model.Chain) error {
	chainDoc, err := r.getChain(ctx, chain.Id)
	if err != nil {
		return err
	}

	result := r.client.Set(ctx, getChainKey(chainDoc.Id), ChainDocument{
		Id:                        chainDoc.Id,
		BlockNumber:               chainDoc.BlockNumber,
		Paused:                    chain.Paused,
		RpcUrl:                    chain.RpcUrl,
		BlockNumberWaitPeriodInMS: chain.BlockNumberWaitPeriodInMS,
		BlockNumberReadPeriodInMS: chain.BlockNumberReadPeriodInMS,
		ConcurrentBlockReadLimit:  chain.ConcurrentBlockReadLimit,
		UpdatedAt:                 chain.UpdatedAt,
	}, time.Duration(0))

	return result.Err()
}

func (r *redisDataSource) GetChain(ctx context.Context, chainId int64) (chain model.Chain, err error) {
	chainDoc, err := r.getChain(ctx, chainId)
	if err != nil {
		if err == redis.Nil {
			return chain, errors.NotFound("chain not found by id: %d", chainId)
		}
		return
	}

	return model.Chain{
		Id:                        chainDoc.Id,
		BlockNumber:               chainDoc.BlockNumber,
		Paused:                    chainDoc.Paused,
		UpdatedAt:                 chainDoc.UpdatedAt,
		RpcUrl:                    chainDoc.RpcUrl,
		BlockNumberWaitPeriodInMS: chainDoc.BlockNumberWaitPeriodInMS,
		BlockNumberReadPeriodInMS: chainDoc.BlockNumberReadPeriodInMS,
		ConcurrentBlockReadLimit:  chainDoc.ConcurrentBlockReadLimit,
	}, err
}

func (r *redisDataSource) GetChains(ctx context.Context) (chains []model.Chain, err error) {
	chainIds, err := r.getChainIds(ctx)
	if err != nil {
		return
	}

	for _, chainId := range chainIds {
		chain, err := r.getChain(ctx, chainId)
		if err != nil {
			return chains, err
		}

		chains = append(chains, model.Chain{
			Id:                        chain.Id,
			BlockNumber:               chain.BlockNumber,
			Paused:                    chain.Paused,
			UpdatedAt:                 chain.UpdatedAt,
			RpcUrl:                    chain.RpcUrl,
			BlockNumberWaitPeriodInMS: chain.BlockNumberWaitPeriodInMS,
			BlockNumberReadPeriodInMS: chain.BlockNumberReadPeriodInMS,
			ConcurrentBlockReadLimit:  chain.ConcurrentBlockReadLimit,
		})
	}
	return
}

func (r redisDataSource) CreateSchema(ctx context.Context, schema model.Schema) (err error) {
	_, err = r.getSchema(ctx, schema.Id)
	if err == nil {
		return errors.AlreadyExists("schema already exists with given id: %s", schema.Id)
	}

	if err != redis.Nil {
		return err
	}

	var eventDocs []SchemaEventDocument
	for _, event := range schema.Events {
		eventDocs = append(eventDocs, SchemaEventDocument{
			Name:  event.Name,
			Alias: event.Alias,
		})
	}

	result := r.client.Set(ctx, getSchemaKey(schema.Id), &SchemaDocument{
		Id:     schema.Id,
		Name:   schema.Name,
		Abi:    schema.Abi,
		Events: eventDocs,
	}, time.Duration(0))

	err = result.Err()
	if err != nil {
		return
	}

	schemaIds, err := r.getSchemaIds(ctx)
	if err != nil {
		return err
	}

	schemaIds = append(schemaIds, schema.Id)
	result = r.client.Set(ctx, getSchemasKey(), schemaIds, time.Duration(0))
	return result.Err()
}

func (r redisDataSource) GetSchemas(ctx context.Context) (schemas []model.Schema, err error) {
	schemaIds, err := r.getSchemaIds(ctx)
	if err != nil {
		if err == redis.Nil {
			return schemas, nil
		}

		return schemas, err
	}

	for _, schemaId := range schemaIds {
		schemaDoc, err := r.getSchema(ctx, schemaId)
		if err != nil {
			return schemas, err
		}

		var events []model.SchemaEvent
		for _, eventDoc := range schemaDoc.Events {
			events = append(events, model.SchemaEvent{
				Name:  eventDoc.Name,
				Alias: eventDoc.Alias,
			})
		}

		schemas = append(schemas, model.Schema{
			Id:     schemaDoc.Id,
			Name:   schemaDoc.Name,
			Abi:    schemaDoc.Abi,
			Events: events,
		})
	}

	return schemas, nil
}

func (r redisDataSource) GetSchema(ctx context.Context, id string) (schema model.Schema, err error) {
	schemaDoc, err := r.getSchema(ctx, id)
	if err != nil {
		if err == redis.Nil {
			return schema, errors.NotFound("schema not found by id: %d", id)
		}

		return schema, err
	}

	var events []model.SchemaEvent
	for _, eventDoc := range schemaDoc.Events {
		events = append(events, model.SchemaEvent{
			Name:  eventDoc.Name,
			Alias: eventDoc.Alias,
		})
	}

	return model.Schema{
		Id:     schemaDoc.Id,
		Name:   schemaDoc.Name,
		Abi:    schemaDoc.Abi,
		Events: events,
	}, err

}

func (r *redisDataSource) Close(context.Context) error {
	return r.client.Close()
}

func (r redisDataSource) getChain(ctx context.Context, chainId int64) (chain ChainDocument, err error) {
	chainJson, err := r.client.Get(ctx, getChainKey(chainId)).Bytes()
	if err != nil {
		return chain, err
	}

	err = json.Unmarshal(chainJson, &chain)
	if err != nil {
		return chain, errors.From(err, "getChain unmarshall errors for chainId: %d", chainId)
	}

	return chain, err
}

func (r redisDataSource) getSchema(ctx context.Context, schemaId string) (schema SchemaDocument, err error) {
	schemaJson, err := r.client.Get(ctx, getSchemaKey(schemaId)).Bytes()
	if err != nil {
		return schema, err
	}

	err = json.Unmarshal(schemaJson, &schema)
	if err != nil {
		return schema, errors.From(err, "getSchema unmarshall errors for schemaId: %s", schemaId)
	}

	return schema, err
}

func (r redisDataSource) getChainIds(ctx context.Context) (chainIds []int64, err error) {
	chainsJson, err := r.client.Get(ctx, getChainsKey()).Bytes()
	if err != nil {
		return chainIds, err
	}

	err = json.Unmarshal(chainsJson, &chainIds)
	if err != nil {
		return chainIds, errors.From(err, "getChainIds unmarshall failed")
	}

	return chainIds, err
}

func (r redisDataSource) getSchemaIds(ctx context.Context) (schemas []string, err error) {
	schemasJson, err := r.client.Get(ctx, getSchemasKey()).Bytes()
	if err != nil {
		return schemas, err
	}

	err = json.Unmarshal(schemasJson, &schemas)
	if err != nil {
		return schemas, fmt.Errorf("getSchemas unmarshall errors, err: %s", err.Error())
	}

	return schemas, err
}

func (r redisDataSource) init(ctx context.Context) (err error) {
	_, err = r.getChainIds(ctx)
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		result := r.client.Set(ctx, getChainsKey(), []int64{}, time.Duration(0))
		err = result.Err()
		if err != nil {
			return err
		}
	}

	_, err = r.getSchemaIds(ctx)
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		result := r.client.Set(ctx, getSchemasKey(), []string{}, time.Duration(0))
		err = result.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func getChainKey(chainId int64) string {
	return fmt.Sprintf("ethereum-connector:chain:%d", chainId)
}

func getChainsKey() string {
	return "ethereum-connector:chains"
}

func getSchemaKey(schemaId string) string {
	return fmt.Sprintf("ethereum-connector:schema:%s", schemaId)
}

func getSchemasKey() string {
	return "ethereum-connector:schemas"
}
*/
