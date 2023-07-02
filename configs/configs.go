package configs

import (
	"github.com/Netflix/go-env"
	"log"
	"strings"
)

type LoggingConfig struct {
	Level              string `env:"LOG_LEVEL,default=info"`
	TimeFieldFormat    string `env:"LOG_TIME_FIELD_FORMAT,default=2006-01-02T15:04:05Z07:00"`
	TimestampFieldName string `env:"LOG_TIMESTAMP_FIELD_NAME,default=time"`
	LevelFieldName     string `env:"LOG_LEVEL_FIELD_NAME,default=level"`
	MessageFieldName   string `env:"LOG_MESSAGE_FIELD_NAME,default=message"`
}

type HttpConfig struct {
	Port           string `env:"HTTP_PORT,default=80"`
	SwaggerEnabled bool   `env:"HTTP_SWAGGER_ENABLED,default=true"`
}

type PostgresDataSourceConfig struct {
	ConnectionString          string `env:"DATA_SOURCE_POSTGRES_CONNECTION_STRING"`
	PersistRawTransactionLogs bool   `env:"DATA_SOURCE_PERSIST_RAW_TRANSACTION_LOGS,default=false"`
}

type DataSourceConfig struct {
	Source   string `env:"DATA_SOURCE,required=true"`
	Postgres PostgresDataSourceConfig
}

type DestinationConfig struct {
	Destination string `env:"DESTINATION,required=true"`
	Postgres    PostgresDestinationConfig
	Kafka       KafkaDestinationConfig
	RabbitMQ    RabbitMQDestinationConfig
	BigQuery    BigQueryDestinationConfig
}

type PostgresDestinationConfig struct {
	ConnectionString          string `env:"DESTINATION_POSTGRES_CONNECTION_STRING"`
	PersistRawTransactionLogs bool   `env:"DESTINATION_POSTGRES_PERSIST_RAW_TRANSACTION_LOGS,default=false"`
}

type BigQueryDestinationConfig struct {
	ProjectId                 string `env:"DESTINATION_BIGQUERY_PROJECT_ID"`
	Dataset                   string `env:"DESTINATION_BIGQUERY_DATASET"`
	DeliveryGuarantee         string `env:"DESTINATION_BIGQUERY_DELIVERY_GUARANTEE,default=AT_LEAST_ONCE"`
	PersistRawTransactionLogs bool   `env:"DESTINATION_BIGQUERY_PERSIST_RAW_TRANSACTION_LOGS,default=false"`
}

type RabbitMQDestinationConfig struct {
	ConnectionString                    string `env:"DESTINATION_RABBITMQ_CONNECTION_STRING"`
	DeliveryGuarantee                   string `env:"DESTINATION_RABBITMQ_DELIVERY_GUARANTEE,default=AT_LEAST_ONCE"`
	BlocksExchangeName                  string `env:"DESTINATION_RABBITMQ_BLOCKS_EXCHANGE_NAME"`
	TransactionsExchangeName            string `env:"DESTINATION_RABBITMQ_TRANSACTIONS_EXCHANGE_NAME"`
	TransactionLogsExchangeName         string `env:"DESTINATION_RABBITMQ_TRANSACTION_LOGS_EXCHANGE_NAME"`
	RawTransactionLogsExchangeName      string `env:"DESTINATION_RABBITMQ_RAW_TRANSACTION_LOGS_EXCHANGE_NAME"`
	SendTransactionLogsToAliasExchanges bool   `env:"DESTINATION_RABBITMQ_SEND_TRANSACTION_LOGS_TO_ALIAS_EXCHANGES,default=true"`
}

type KafkaDestinationConfig struct {
	BootstrapServers                 string `env:"DESTINATION_KAFKA_BOOTSTRAP_SERVERS"`
	DeliveryGuarantee                string `env:"DESTINATION_KAFKA_DELIVERY_GUARANTEE,default=AT_LEAST_ONCE"`
	ProducerTransactionId            string `env:"DESTINATION_KAFKA_PRODUCER_TRANSACTION_ID,default=ethereum-connector"`
	BlocksTopicName                  string `env:"DESTINATION_KAFKA_BLOCKS_TOPIC_NAME"`
	TransactionsTopicName            string `env:"DESTINATION_KAFKA_TRANSACTIONS_TOPIC_NAME"`
	TransactionLogsTopicName         string `env:"DESTINATION_KAFKA_TRANSACTION_LOGS_TOPIC_NAME"`
	RawTransactionLogsTopicName      string `env:"DESTINATION_KAFKA_RAW_TRANSACTION_LOGS_TOPIC_NAME"`
	SendTransactionLogsToAliasTopics bool   `env:"DESTINATION_KAFKA_SEND_TRANSACTION_LOGS_TO_ALIAS_TOPICS,default=true"`
}

func (c KafkaDestinationConfig) GetBootstrapServers() []string {
	return strings.Split(c.BootstrapServers, ",")
}

type Config struct {
	Logging     LoggingConfig
	Http        HttpConfig
	DataSource  DataSourceConfig
	Destination DestinationConfig
}

func ReadConfigs() Config {
	var envConfig Config
	_, err := env.UnmarshalFromEnviron(&envConfig)
	if err != nil {
		log.Fatalf("config err: %s", err.Error())
	}

	if envConfig.DataSource.Source == "postgres" &&
		envConfig.Destination.Destination == "postgres" &&
		envConfig.DataSource.Postgres.ConnectionString == envConfig.Destination.Postgres.ConnectionString {
		log.Fatalf("Source and destination databases must not be the same")
	}

	return envConfig
}
