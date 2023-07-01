package destination

import (
	"context"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/aliarbak/ethereum-connector/utils"
	"strings"
)

type DeliveryGuarantee string

const (
	AtMostOnceDeliveryGuarantee  DeliveryGuarantee = "AT_MOST_ONCE"
	AtLeastOnceDeliveryGuarantee DeliveryGuarantee = "AT_LEAST_ONCE"
	ExactlyOnceDeliveryGuarantee DeliveryGuarantee = "EXACTLY_ONCE"
)

const (
	destinationPostgres = "postgres"
	destinationKafka    = "kafka"
	destinationRabbitMQ = "rabbitmq"
	destinationBigQuery = "bigquery"
)

var excludedLogFields = []string{model.EventNameLogField, model.EventAliasLogField, model.ContractAddressLogField, model.LogHashLogField, model.LogIndexLogField, model.RawDataLogField}

type Destination interface {
	init(ctx context.Context) error
	SendBlock(ctx context.Context, block model.Block) error
	SendSyncLogs(_ context.Context, block model.Block) error
	DeliveryGuarantee() DeliveryGuarantee
	Close(ctx context.Context) error
}

type Factory interface {
	CreateDestination(ctx context.Context) (Destination, error)
}

func NewFactory(ctx context.Context, config configs.DestinationConfig) (Factory, error) {
	var factory Factory
	switch strings.ToLower(config.Destination) {
	case destinationPostgres:
		factory = newPostgresFactory(config.Postgres)
		break
	case destinationKafka:
		factory = newKafkaFactory(config.Kafka)
		break
	case destinationRabbitMQ:
		factory = newRabbitMQFactory(config.RabbitMQ)
		break
	case destinationBigQuery:
		factory = newBigQueryFactory(config.BigQuery)
		break
	default:
		return nil, errors.New("invalid destination type: %s", config.Destination)
	}

	dest, err := factory.CreateDestination(ctx)
	if err != nil {
		return nil, err
	}

	defer dest.Close(ctx)
	err = dest.init(ctx)
	return factory, err
}

func getLogDataMap(originalMap map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	for k, v := range originalMap {
		if utils.StringArrayContains(excludedLogFields, k) {
			continue
		}

		data[k] = v
	}

	return data
}
