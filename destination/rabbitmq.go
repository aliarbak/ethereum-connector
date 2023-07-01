package destination

import (
	"context"
	"encoding/json"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type rabbitmqDestinationFactory struct {
	config     configs.RabbitMQDestinationConfig
	connection *amqp.Connection
}

type rabbitmqDestination struct {
	channel                             *amqp.Channel
	deliveryGuarantee                   DeliveryGuarantee
	blocksExchangeName                  string
	transactionsExchangeName            string
	transactionLogsExchangeName         string
	rawTransactionLogsExchangeName      string
	sendTransactionLogsToAliasExchanges bool
	declaredExchanges                   map[string]bool
}

func newRabbitMQFactory(config configs.RabbitMQDestinationConfig) Factory {
	connection, err := amqp.Dial(config.ConnectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("rabbitmq connection failed")
		return nil
	}

	return &rabbitmqDestinationFactory{
		config:     config,
		connection: connection,
	}
}

func (f *rabbitmqDestinationFactory) CreateDestination(context.Context) (dest Destination, err error) {
	if f.config.DeliveryGuarantee != string(AtLeastOnceDeliveryGuarantee) && f.config.DeliveryGuarantee != string(AtMostOnceDeliveryGuarantee) {
		return nil, errors.InvalidInput("invalid rabbitmq destination delivery guarantee: %s", f.config.DeliveryGuarantee)
	}

	channel, err := f.connection.Channel()
	if err != nil {
		return dest, err
	}

	return &rabbitmqDestination{
		channel:                             channel,
		deliveryGuarantee:                   DeliveryGuarantee(f.config.DeliveryGuarantee),
		blocksExchangeName:                  f.config.BlocksExchangeName,
		transactionsExchangeName:            f.config.TransactionsExchangeName,
		transactionLogsExchangeName:         f.config.TransactionLogsExchangeName,
		rawTransactionLogsExchangeName:      f.config.RawTransactionLogsExchangeName,
		sendTransactionLogsToAliasExchanges: f.config.SendTransactionLogsToAliasExchanges,
		declaredExchanges:                   map[string]bool{},
	}, err
}

func (r rabbitmqDestination) DeliveryGuarantee() DeliveryGuarantee {
	return r.deliveryGuarantee
}

func (r rabbitmqDestination) SendBlock(_ context.Context, block model.Block) (err error) {
	if err = r.channel.Tx(); err != nil {
		return errors.From(err, "rabbitmq begin transaction failed")
	}

	if len(r.blocksExchangeName) > 0 {
		if err = r.publishMessage(r.blocksExchangeName, newBlockMessage(block)); err != nil {
			return r.rollbackTx(err)
		}
	}

	for _, transaction := range block.Transactions {
		if len(r.transactionsExchangeName) > 0 {
			if err = r.publishMessage(r.transactionsExchangeName, newTransactionMessage(block, transaction)); err != nil {
				return r.rollbackTx(err)
			}
		}

		for _, transactionLog := range transaction.Logs {
			if len(r.rawTransactionLogsExchangeName) > 0 {
				if err = r.publishMessage(r.rawTransactionLogsExchangeName, transactionLog[model.RawDataLogField]); err != nil {
					return r.rollbackTx(err)
				}
			}

			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogMessageValue := newTransactionLogMessage(block, transaction, transactionLog)
			if len(r.transactionLogsExchangeName) > 0 {
				if err = r.publishMessage(r.transactionLogsExchangeName, transactionLogMessageValue); err != nil {
					return r.rollbackTx(err)
				}
			}

			if r.sendTransactionLogsToAliasExchanges {
				if err = r.publishMessage(transactionLog[model.EventAliasLogField].(string), transactionLogMessageValue); err != nil {
					return r.rollbackTx(err)
				}
			}
		}
	}

	err = r.channel.TxCommit()
	if err != nil {
		return errors.From(err, "rabbitmq transaction commit failed")
	}

	return
}

func (r rabbitmqDestination) SendSyncLogs(_ context.Context, block model.Block) (err error) {
	if err = r.channel.Tx(); err != nil {
		return errors.From(err, "rabbitmq begin transaction failed")
	}

	for _, transaction := range block.Transactions {
		for _, transactionLog := range transaction.Logs {
			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogMessageValue := newTransactionLogMessage(block, transaction, transactionLog)
			if len(r.transactionLogsExchangeName) > 0 {
				if err = r.publishMessage(r.transactionLogsExchangeName, transactionLogMessageValue); err != nil {
					return r.rollbackTx(err)
				}
			}

			if r.sendTransactionLogsToAliasExchanges {
				if err = r.publishMessage(transactionLog[model.EventAliasLogField].(string), transactionLogMessageValue); err != nil {
					return r.rollbackTx(err)
				}
			}
		}
	}

	err = r.channel.TxCommit()
	if err != nil {
		return errors.From(err, "rabbitmq transaction commit failed")
	}

	return
}

func (r rabbitmqDestination) rollbackTx(err error) error {
	rollbackErr := r.channel.TxRollback()
	if rollbackErr != nil {
		return errors.From(rollbackErr, "rabbitmq rollback failed, original error: %s", err.Error())
	}

	return err
}

func (r rabbitmqDestination) Close(context.Context) error {
	return r.channel.Close()
}

func (r rabbitmqDestination) publishMessage(exchangeName string, msg interface{}) error {
	body, err := r.newMessage(msg)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		Body: body,
	}

	if err = r.declareExchangeIfNotExists(exchangeName); err != nil {
		return err
	}

	return r.channel.Publish(exchangeName, "", false, false, message)
}

func (r rabbitmqDestination) newMessage(value interface{}) (valueJson []byte, err error) {
	switch value.(type) {
	case []byte:
		return value.([]byte), nil
	case string:
		return []byte(value.(string)), nil
	default:
		valueJson, err = json.Marshal(value)
		if err != nil {
			return nil, errors.From(err, "rabbitmq message marshall failed for: %+v", value)
		}
		return valueJson, nil
	}
}

func (r rabbitmqDestination) init(context.Context) (err error) {
	if len(r.blocksExchangeName) > 0 {
		if err = r.declareExchangeIfNotExists(r.blocksExchangeName); err != nil {
			return errors.From(err, "block exchange declare failed for exchange name: %s", r.blocksExchangeName)
		}
	}

	if len(r.transactionsExchangeName) > 0 {
		if err = r.declareExchangeIfNotExists(r.transactionsExchangeName); err != nil {
			return errors.From(err, "transactions exchange declare failed for exchange name: %s", r.transactionsExchangeName)
		}
	}

	if len(r.transactionLogsExchangeName) > 0 {
		if err = r.declareExchangeIfNotExists(r.transactionLogsExchangeName); err != nil {
			return errors.From(err, "transaction logs exchange declare failed for exchange name: %s", r.transactionLogsExchangeName)
		}
	}

	if len(r.rawTransactionLogsExchangeName) > 0 {
		if err = r.declareExchangeIfNotExists(r.rawTransactionLogsExchangeName); err != nil {
			return err
		}
	}

	return nil
}

func (r *rabbitmqDestination) declareExchangeIfNotExists(exchangeName string) error {
	if r.declaredExchanges[exchangeName] == true {
		return nil
	}

	err := r.channel.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		return errors.From(err, "rabbitmq exchange declare failed for: %s", exchangeName)
	}

	r.declaredExchanges[exchangeName] = true
	return nil
}
