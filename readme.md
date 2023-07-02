<!--
  Title: Ethereum Connector
  Description: Export ethereum blocks, transaction and transaction logs to Kafka, RabbitMQ, BigQuery or Postgres
  Author: aliarbak
  -->

# Ethereum Connector

The Ethereum Connector lets you listen to the Ethereum network (or other evm chains) and export blocks, transactions, and transaction logs to a desired destination. You can also export meaningful transaction log data (ERC-721 Transfer events, etc.) by providing ABIs.

Unlike applications designed to listen specific contract addresses, this project is created to listen to the entire network and export its data. It is suitable for use cases where tracking all activities on the network is essential, such as analytics tools or NFT marketplaces.

Supported Destinations:
- Kafka
- RabbitMQ
- BigQuery
- PostgreSQL
- *(coming soon) MongoDB*
- *(coming soon) Apache Cassandra*

Supported Data Sources:
- PostgreSQL
- *(coming soon) MongoDB*

**This app is in Alpha and should be used with caution. It has only been manually tested and has not proven itself in the production yet.**

------------

### Table of Contents
- [How Does It Work?](#how-does-it-work)
- [Quickstart](#quickstart)
- [Data Sources](#data-sources)
  - [PostgreSQL](#--postgresql)
- [Destinations](#destinations)
  - [Delivery Guarantee](#delivery-guarantee)
  - [Kafka](#--kafka)
  - [RabbitMQ](#--rabbitmq)
  - [BigQuery](#--bigquery)
  - [PostgreSQL](#--postgresql-1)
- [Chains and Schemas](#chains-and-schemas)
  - [Registering a Chain (`POST /chains`)](#registering-a-chain-post-chains)
  - [Registering a Schema (`POST /schemas`)](#registering-a-schema-post-schemas)
    - [Event Aliases](#event-aliases)


------------

## How Does It Work?

- You register the schemas and the chains you want to listen to.
- The EVM listener starts processing blocks from the specified block.
- For each block:
  - It retrieves the receipts of transactions in the block.
  - It processes the logs within transactions:
    - It creates a raw transaction log entry for the log.
    - If there is a matching event schema for the log, it creates a transaction log entry using the event ABI.
- It sends the block record to the destination.
- It sends the transaction records to the destination.
- It sends the transaction log records matched with event schemas to the destination.
- (If enabled) It sends the raw transaction log records to the destination.
- It proceeds to the next block.
  - If it has reached the last block, it waits for the next block to be added to the chain.

## Quickstart

Let's start the connector with Postgres datasource and Kafka destination:

```bash
docker run -d -p 8080:80 \
  -e DATASOURCE=postgres \
  -e DATASOURCE_POSTGRES_CONNECTION_STRING=postgres://postgresuser:password@localhost:5432/sourcedb \
  -e DATA_SOURCE_PERSIST_RAW_TRANSACTION_LOGS=true \
  -e DESTINATION=kafka \
  -e DESTINATION_KAFKA_BOOTSTRAP_SERVERS=localhost:9092,localhost:9093,localhost:9094 \
  -e DESTINATION_KAFKA_BLOCKS_TOPIC_NAME=blocks \
  -e DESTINATION_KAFKA_TRANSACTIONS_TOPIC_NAME=transactions \
  -e DESTINATION_KAFKA_TRANSACTION_LOGS_TOPIC_NAME=transaction_logs \
  -e DESTINATION_KAFKA_RAW_TRANSACTION_LOGS_TOPIC_NAME=raw_transaction_logs \
  -e DESTINATION_KAFKA_SEND_TRANSACTION_LOGS_TO_ALIAS_TOPICS=true \
  aliarbak/ethereum-connector
```

Once the container is up and running, you can access the swagger via: `http://localhost:8080/swagger/index.html`

To start listening to the Ethereum MAINNET (id: 1):
```curl
curl --location --request POST 'http://localhost/chains' \
--header 'accept: application/json' \
--header 'Content-Type: application/json' \
--data-raw '{ "id": 1, "rpcUrl": "https://eth-mainnet.g.alchemy.com/v2/<token>", "concurrentBlockReadLimit": 1 }'
```

**That's it. The connector has started reading blocks from the Ethereum network and block, transaction, and raw transaction log data are being sent to Kafka.**

Now, let's dive into more details.

------------

## Data Sources
The Ethereum Connector needs a data source to store the states of the chains it is processing and schema information. It only supports PostgreSQL as a data source for now, but additional data sources (initially MongoDB) will be supported soon.

### - PostgreSQL
You need to create a source database on postgres and set the related environment variables, that's it.

You should also set `DATASOURCE=postgres`
|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_POSTGRES_CONNECTION_STRING | The connection string of the source database | Yes | - | `postgres://postgresuser:password@localhost:5432/sourcedb` |
| DESTINATION_POSTGRES_PERSIST_RAW_TRANSACTION_LOGS | Set it true if you want to persist raw transaction logs on the source database. See "Schema Filters" section for more details. | No | false | true or false |

------------

## Destinations

With the ethereum connector, you can connect EVM chains to your Kafka, RabbitMQ, BigQuery or PostgreSQL destinations.

### Delivery Guarantee
Delivery Guarantee options vary depending on the destination type used. There are three types of delivery guarantees:
- **EXACTLY_ONCE**: Block, transaction and transaction log records are sent to the destination exactly once. This is only possible for PostgreSQL destination type.
- **AT_LEAST_ONCE**: Records are sent to the destination at least once. If a problem occurs in the connector, it may be possible to send the same records more than once. Therefore, your message streaming consumers should be idempotent, and you need to deduplicate the records if your destination is BigQuery.
- **AT_MOST_ONCE**: Records are sent to the destination only once. If a problem occurs in the connector, some records may not reach the destination at all and data loss occurs. Unless you have a specific reason, this option is **not recommended**.

### - Kafka
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


### - RabbitMQ
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

### - BigQuery
You should set `DESTINATION=bigquery`

|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_BIGQUERY_PROJECT_ID | BigQuery Project ID | **Yes** | - | xyz-123 |
| DESTINATION_BIGQUERY_DATASET | BigQuery Data Set | **Yes** | - | chain |
| GOOGLE_APPLICATION_CREDENTIALS | The path of your Google Cloud Service Account Key json file  | **Yes** | - | /Users/username/service_account.json |
| DESTINATION_BIGQUERY_DELIVERY_GUARANTEE | The delivery guarantee option can only be AT_LEAST_ONCE or AT_MOST_ONCE for BigQuery. | No | `AT_LEAST_ONCE` | `AT_LEAST_ONCE` or `AT_MOST_ONCE` |
| DESTINATION_BIGQUERY_PERSIST_RAW_TRANSACTION_LOGS | Set this value to true if you want to persist the raw transaction logs in the BigQuery | No | false | true or false |

### - PostgreSQL
You should set `DESTINATION=postgres`

|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_POSTGRES_CONNECTION_STRING | The connection string of the destination database  | **Yes** | - | `postgres://postgresuser:password@localhost:5432/sourcedb` |
| DESTINATION_POSTGRES_PERSIST_RAW_TRANSACTION_LOGS | Set this value to true if you want to persist the raw transaction logs in your destination database | No | false | true or false |


------------

## Chains and Schemas

### Registering a Chain (`POST /chains`)
You can register EVM chains with this endpoint. **Once you registered a chain, you should never update the chain's blockNumber manually.**

Request Payload
|  Field Name  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| id | chainId, you can find the list of chain ids from: https://chainlist.org | **Yes** | - | 1 |
| rpcUrl  | EVM JSON RPC url. You can use node providers like alchemy or infura.  | **Yes**  | -  | https://eth-mainnet.g.alchemy.com/v2/<token> |
| blockNumber  | The blockNumber you want to start from listening. -1 = genesis block  | No  | -1  | -1 or 17605777  |
| paused |  Set it true if you want to stop reading blocks. You can update it later with the PATCH endpoint. | No  | false  | true or false  |
| concurrentBlockReadLimit | All blocks are read sequentially and sent to the destination. To increase the processing speed of blocks, you can enable concurrent reading (while still sending them sequentially). However, during concurrent reading, multiple requests are sent to your RPC node simultaneously, so you need to pay attention to the rate limits of your node provider.  | No | 1  | 5 |
| blockNumberReadPeriodInMs  |  If you want to apply a delay (for throttling purposes) after processing each block before sending another request to your RPC node, you can specify the waiting time in milliseconds.  | No  | 500 (ms)  | 200  |
| blockNumberWaitPeriodInMs  |  How long it takes (in milliseconds) to check if a new block has been added to the chain after it has been read to the last block  | No  |  2000 (ms) |  500  |

Sample Request Payload:
```json
{
  "id": 1,  
  "blockNumber": -1,
  "blockNumberReadPeriodInMs": 200,
  "blockNumberWaitPeriodInMs": 2000,
  "concurrentBlockReadLimit": 2,
  "paused": false,
  "rpcUrl": "https://eth-mainnet.g.alchemy.com/v2/<token>"
}
```

------------

### Registering a Schema (`POST /schemas`)

The connector does not know the details of log topics and log data when reading transaction logs. Therefore, it can extract these event logs only as raw transaction logs.

You should provide ABI schemas if you want to extract specific event logs(like ERC-20 Transfer, ERC-721 Transfer, ERC-1155 TransferSingle, etc.). For example, if you provide the ERC-721 schema, the connector can process ERC-721 Transfer logs and send them to the destination as ERC-721 Transfers.

Request Payload
|  Field Name  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| name | The name of the schema | **Yes** | - | ERC-721, ERC-1155, etc. |
| abi  | The ABI of the contract. You only need to include event ABIs.  | **Yes**  | -  | (see the sample request below) |
| events  | You need to specify the names and the aliases of the events you want to extract. For example, for ERC-721, you may want only to include Transfer events and exclude others (Approval events, etc.) You also need to specify the aliases for the events (see the Event Alias details section below for more details) | Yes  | -  | (see the sample request below)  |

Sample Request Payload:
```json
{
  "name": "ERC-721",
  "abi": [
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "tokenId",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    }
  ],
  "events": [
    {
      "alias": "ERC721Transfer",
      "name": "Transfer"
    }
  ]
}
```

#### Event Aliases
Event aliases serve two main purposes:
- Preventing the mixing of events with the same name but belonging to different contract schemas. For example, the "Transfer" event in ERC-721 and the "Transfer" event in ERC-20 share the same name.
- Facilitating the separating and processing of records transferred to the destination.

For example, if you have registered the ERC-721 contract schema and assigned the "ERC721Transfer" alias to the Transfer event, the following will occur:
- In the ERC-721 transaction_log records saved to the destination, the "eventAlias" field will appear as "ERC721Transfer." This allows for easier separating and processing of transaction logs specific to ERC-721.
- If your destination is a message streaming tool like RabbitMQ or Kafka, in addition to the transaction_log topics and exchanges, you can also choose to send events to alias-specific topics or exchanges. For instance, if your destination is Kafka, ERC-721 Transfer events will be directly sent to the "ERC721Transfer" topic. This means that consuming this topic will be sufficient for handling ERC-721 transfer events.

------------