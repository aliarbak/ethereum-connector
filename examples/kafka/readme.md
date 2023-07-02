# Kafka Destination

You should set `DESTINATION=kafka`

In cases where you do not enter the topic name, the relevant messages will not be forwarded to the Kafka.

|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_KAFKA_BOOTSTRAP_SERVERS | Comma-seperated kafka bootstrap server addresses | **Yes** | - | `localhost:9092,localhost:9093,localhost:9094` |
| DESTINATION_KAFKA_DELIVERY_GUARANTEE | The delivery guarantee option can only be AT_LEAST_ONCE or AT_MOST_ONCE for Kafka. | No | `AT_LEAST_ONCE` | `AT_LEAST_ONCE` or `AT_MOST_ONCE` |
| DESTINATION_KAFKA_PRODUCER_TRANSACTION_ID | The producer transaction id | No | ethereum-connector | custom-producer-tx-id |
| DESTINATION_KAFKA_BLOCKS_TOPIC_NAME | The topic name for the block records. If you do not set the value, block messages will not be transmitted. | No | - | blocks |
| DESTINATION_KAFKA_TRANSACTIONS_TOPIC_NAME | The topic name for the transaction records. If you do not set the value, transaction messages will not be transmitted. | No | - | transactions |
| DESTINATION_KAFKA_TRANSACTION_LOGS_TOPIC_NAME | The topic name for the transaction log records. If you do not set the value, transaction log messages will not be transmitted. | No | - | transaction_logs |
| DESTINATION_KAFKA_RAW_TRANSACTION_LOGS_TOPIC_NAME | The topic name for the raw transaction log records. If you do not set the value, raw transaction log messages will not be transmitted. | No | - | raw_transaction_logs |
| DESTINATION_KAFKA_SEND_TRANSACTION_LOGS_TO_ALIAS_TOPICS | If you want transaction logs to be transmitted to event alias topics, set this value to true. See the Event Alias section for more details. | No | true | true or false |

