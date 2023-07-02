# RabbitMQ Destination

You should set `DESTINATION=rabbitmq`

In cases where you do not enter the exchange name, the relevant messages will not be forwarded to the RabbitMQ.

|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_RABBITMQ_CONNECTION_STRING | The connection string of the RabbitMQ | **Yes** | - | `amqp://guest:guest@localhost:5672` |
| DESTINATION_RABBITMQ_DELIVERY_GUARANTEE | The delivery guarantee option can only be AT_LEAST_ONCE or AT_MOST_ONCE for RabbitMQ. | No | `AT_LEAST_ONCE` | `AT_LEAST_ONCE` or `AT_MOST_ONCE` |
| DESTINATION_RABBITMQ_BLOCKS_EXCHANGE_NAME | The exchange name for the block records. If you do not set the value, block messages will not be transmitted. | No | - | blocks |
| DESTINATION_RABBITMQ_TRANSACTIONS_EXCHANGE_NAME | The exchange name for the transaction records. If you do not set the value, transaction messages will not be transmitted. | No | - | transactions |
| DESTINATION_RABBITMQ_TRANSACTION_LOGS_EXCHANGE_NAME | The exchange name for the transaction log records. If you do not set the value, transaction log messages will not be transmitted. | No | - | transaction_logs |
| DESTINATION_RABBITMQ_RAW_TRANSACTION_LOGS_EXCHANGE_NAME | The exchange name for the raw transaction log records. If you do not set the value, raw transaction log messages will not be transmitted. | No | - | raw_transaction_logs |
| DESTINATION_RABBITMQ_SEND_TRANSACTION_LOGS_TO_ALIAS_EXCHANGES | If you want transaction logs to be transmitted to event alias exchanges, set this value to true. See the Event Alias section for more details. | No | true | true or false |
