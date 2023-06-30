package destination

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/aliarbak/ethereum-connector/utils"
)

type kafkaDestinationFactory struct {
	config        configs.KafkaDestinationConfig
	producerCount int64
}

type kafkaDestination struct {
	producer                      sarama.SyncProducer
	deliveryGuarantee             DeliveryGuarantee
	bootstrapServers              []string
	blocksTopicName               string
	transactionsTopicName         string
	transactionLogsTopicName      string
	rawTransactionLogsTopicName   string
	sendTransferLogsToAliasTopics bool
}

func newKafkaFactory(config configs.KafkaDestinationConfig) Factory {
	sarama.Logger = utils.NewStdLogger()
	return &kafkaDestinationFactory{
		config: config,
	}
}

func (f *kafkaDestinationFactory) CreateDestination(context.Context) (dest Destination, err error) {
	f.producerCount++
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Transaction.ID = fmt.Sprintf("%s-%d", f.config.ProducerTransactionId, f.producerCount)
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Idempotent = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.MaxOpenRequests = 1

	bootstrapServers := f.config.GetBootstrapServers()
	producer, err := sarama.NewSyncProducer(bootstrapServers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	if f.config.DeliveryGuarantee != string(AtLeastOnceDeliveryGuarantee) && f.config.DeliveryGuarantee != string(AtMostOnceDeliveryGuarantee) {
		return nil, errors.InvalidInput("invalid kafka destination delivery guarantee: %s", f.config.DeliveryGuarantee)
	}

	return &kafkaDestination{
		producer:                      producer,
		deliveryGuarantee:             DeliveryGuarantee(f.config.DeliveryGuarantee),
		blocksTopicName:               f.config.BlocksTopicName,
		transactionsTopicName:         f.config.TransactionsTopicName,
		transactionLogsTopicName:      f.config.TransactionLogsTopicName,
		rawTransactionLogsTopicName:   f.config.RawTransactionLogsTopicName,
		sendTransferLogsToAliasTopics: f.config.SendTransferLogsToAliasTopics,
		bootstrapServers:              bootstrapServers,
	}, err
}

func (r kafkaDestination) DeliveryGuarantee() DeliveryGuarantee {
	return r.deliveryGuarantee
}

func (r kafkaDestination) SendBlock(_ context.Context, block model.Block) error {
	err := r.producer.BeginTxn()
	if err != nil {
		return errors.From(err, "kafka begin transaction failed")
	}

	messages, err := r.prepareMessages(block)
	if err != nil {
		r.producer.AbortTxn()
		return err
	}

	err = r.producer.SendMessages(messages)
	if err != nil {
		r.producer.AbortTxn()
		return err
	}

	return r.producer.CommitTxn()
}

func (r kafkaDestination) SendSyncLogs(_ context.Context, block model.Block) error {
	err := r.producer.BeginTxn()
	if err != nil {
		return errors.From(err, "kafka begin transaction failed")
	}

	messages, err := r.prepareSyncMessages(block)
	if err != nil {
		r.producer.AbortTxn()
		return err
	}

	err = r.producer.SendMessages(messages)
	if err != nil {
		r.producer.AbortTxn()
		return err
	}

	return r.producer.CommitTxn()
}

func (r kafkaDestination) prepareMessages(block model.Block) (messages []*sarama.ProducerMessage, err error) {
	if len(r.blocksTopicName) > 0 {
		blockMessage, err := newMessage(r.blocksTopicName, block.Number.String(), newBlockMessage(block))
		if err != nil {
			return messages, errors.From(err, "block message initialization failed, blockNumber: %s", block.Number.String())
		}

		messages = append(messages, blockMessage)
	}

	for _, transaction := range block.Transactions {
		if len(r.transactionsTopicName) > 0 {
			transactionMessage, err := newMessage(r.transactionsTopicName, transaction.Hash, newTransactionMessage(block, transaction))
			if err != nil {
				return messages, errors.From(err, "transaction message initialization failed, txHash: %s", transaction.Hash)
			}

			messages = append(messages, transactionMessage)
		}

		for _, transactionLog := range transaction.Logs {
			if len(r.rawTransactionLogsTopicName) > 0 {
				rawTransactionLogMessage, err := newMessage(r.rawTransactionLogsTopicName, transactionLog[model.ContractAddressLogField].(string), transactionLog[model.RawDataLogField])
				if err != nil {
					return messages, fmt.Errorf("raw transaction log message initialization failed, txHash: %s, log hash: %s, err: %s", transaction.Hash, transactionLog[model.LogHashLogField].(string), err.Error())
				}

				messages = append(messages, rawTransactionLogMessage)
			}

			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogMessageValue := newTransactionLogMessage(block, transaction, transactionLog)
			if len(r.transactionLogsTopicName) > 0 {
				transactionLogMessage, err := newMessage(r.transactionLogsTopicName, transactionLog[model.ContractAddressLogField].(string), transactionLogMessageValue)
				if err != nil {
					return messages, fmt.Errorf("transaction log message initialization failed, txHash: %s, log hash: %s, err: %s", transaction.Hash, transactionLog[model.LogHashLogField].(string), err.Error())
				}

				messages = append(messages, transactionLogMessage)
			}

			if r.sendTransferLogsToAliasTopics {
				transactionLogMessage, err := newMessage(transactionLog[model.EventAliasLogField].(string), transactionLog[model.ContractAddressLogField].(string), transactionLogMessageValue)
				if err != nil {
					return messages, fmt.Errorf("transaction log message for alias topic initialization failed, txHash: %s, log hash: %s, err: %s", transaction.Hash, transactionLog[model.LogHashLogField].(string), err.Error())
				}

				messages = append(messages, transactionLogMessage)
			}
		}
	}

	return messages, err
}

func (r kafkaDestination) Close(context.Context) error {
	return r.producer.Close()
}

func (r kafkaDestination) prepareSyncMessages(block model.Block) (messages []*sarama.ProducerMessage, err error) {
	for _, transaction := range block.Transactions {
		for _, transactionLog := range transaction.Logs {
			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogMessageValue := newTransactionLogMessage(block, transaction, transactionLog)
			if len(r.transactionLogsTopicName) > 0 {
				transactionLogMessage, err := newMessage(r.transactionLogsTopicName, transactionLog[model.ContractAddressLogField].(string), transactionLogMessageValue)
				if err != nil {
					return messages, fmt.Errorf("transaction log message initialization failed, txHash: %s, log hash: %s, err: %s", transaction.Hash, transactionLog[model.LogHashLogField].(string), err.Error())
				}

				messages = append(messages, transactionLogMessage)
			}

			if r.sendTransferLogsToAliasTopics {
				transactionLogMessage, err := newMessage(transactionLog[model.EventAliasLogField].(string), transactionLog[model.ContractAddressLogField].(string), transactionLogMessageValue)
				if err != nil {
					return messages, fmt.Errorf("transaction log message for alias topic initialization failed, txHash: %s, log hash: %s, err: %s", transaction.Hash, transactionLog[model.LogHashLogField].(string), err.Error())
				}

				messages = append(messages, transactionLogMessage)
			}
		}
	}

	return messages, err
}

func (r kafkaDestination) init(context.Context) error {
	client, err := sarama.NewClient(r.bootstrapServers, sarama.NewConfig())
	if err != nil {
		return err
	}

	defer client.Close()
	controller, err := client.Controller()
	if err != nil {
		return err
	}

	defer controller.Close()
	return nil
}

func newMessage(topic string, key string, value interface{}) (message *sarama.ProducerMessage, err error) {
	var valueJson []byte
	switch value.(type) {
	case []byte:
		valueJson = value.([]byte)
		break
	case string:
		valueJson = []byte(value.(string))
		break
	default:
		valueJson, err = json.Marshal(value)
		if err != nil {
			return nil, err
		}
		break
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(valueJson),
	}, nil
}
