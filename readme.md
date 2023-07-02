<!--
  Title: Ethereum Connector
  Description: Connects Ethereum (or other evm chains) to Kafka, RabbitMQ, BigQuery or Postgres
  Author: aliarbak
  -->

# Ethereum Connector

It connects the Ethereum(or other evm chains) to your databases or message-streaming tools by exporting blocks, transactions, and transaction logs to the desired destination. You can also export transaction log data(ERC-721 Transfer events, etc.) as meaningful payloads by providing ABIs(schemas).

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

**Attention: This app is in Alpha stage and should be used with caution. It has only been manually tested and has not proven itself in the production yet.**

------------

### Table of Contents
- [How It Works?](#how-it-works)
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
  - [Filtering Contract Addresses for Schemas (`POST /schemas/{schemaId}/filters`)](#filtering-contract-addresses-for-schemas-post-schemasschemaidfilters)
    - [Synchronizing Filtered Contracts](#synchronizing-filtered-contracts)
- [Additional Configurations](#additional-configurations)
- [Upcoming Improvements](#upcoming-improvements)



------------

## How It Works?

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

Pull the docker image:
```bash
docker pull aliarbak/ethereum-connector
```

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
| DATA_SOURCE_POSTGRES_CONNECTION_STRING | The connection string of the source database | Yes | - | `postgres://postgresuser:password@localhost:5432/sourcedb` |
| DATA_SOURCE_PERSIST_RAW_TRANSACTION_LOGS | Set it true if you want to persist raw transaction logs on the source database. See the "Schema Filters" section for more details. | No | false | true or false |

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

Please visit the [Kafka Destination configuration details page](https://github.com/aliarbak/ethereum-connector/tree/main/examples/kafka).

### - RabbitMQ
You should set `DESTINATION=rabbitmq`

Please visit the [RabbitMQ Destination configuration details page](https://github.com/aliarbak/ethereum-connector/tree/main/examples/rabbitmq).

### - BigQuery
You should set `DESTINATION=bigquery`

Please visit the [BigQuery Destination configuration details page](https://github.com/aliarbak/ethereum-connector/tree/main/examples/bigquery).

### - PostgreSQL
You should set `DESTINATION=postgres`

Please visit the [Postgres Destination configuration details page](https://github.com/aliarbak/ethereum-connector/tree/main/examples/postgres).

------------

## Chains and Schemas

### Registering a Chain (`POST /chains`)
You can register EVM chains with this endpoint. **Once you registered a chain, you should never update the chain's blockNumber manually.**

Request Payload
|  Field Name  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| id | chainId, you can find the list of chain ids from: https://chainlist.org | **Yes** | - | 1 |
| rpcUrl  | EVM JSON RPC url. You can use node providers like alchemy or infura.  | **Yes**  | -  | https://eth-mainnet.g.alchemy.com/v2/<token> |
| blockNumber  | The blockNumber you want to start from listening. -1 = genesis block  | No  | -1  | -1, 17605777, etc.  |
| paused |  Set it true if you want to stop reading blocks. You can update it later with the PATCH endpoint. | No  | false  | true or false  |
| concurrentBlockReadLimit | All blocks are read sequentially and sent to the destination. To increase the processing speed of blocks, you can enable concurrent reading(but it will send them sequentially). However, during concurrent reading, multiple requests are sent to your RPC node simultaneously, so you need to pay attention to the rate limits of your node provider.  | No | 1  | 20 |
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
| name | Name of the schema | **Yes** | - | ERC-721, ERC-1155, etc. |
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

### Filtering Contract Addresses for Schemas (`POST /schemas/{schemaId}/filters`)
When you register a schema, it will be applied to all contracts that match the schema's event ABIs. You can create schema filters if you want to apply schemas(the transaction_logs, not the raw_transaction_logs) only for predefined contract addresses.

And let's say, you are building an NFT marketplace, and you only want to receive ERC-721 Transfer events from specific contracts. Then, you can create schema filters for the ERC-721 schema.

Request Payload
|  Field Name  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| chainId | chainId where the contract is located | **Yes** | - | 1, 1337, etc. |
| contractAddress  | contract address to which the schema will be applied  | **Yes**  | -  | 0x34eebee6942d8def3c125458d1a86e0a897fd6f9 |
| sync  | Set it true if you want to synchronize that contract's events (see the sync descriptions below) | No  | false  | true or false  |

The schema will only be applied to filtered contracts when it has at least one filtered contract.

#### Synchronizing Filtered Contracts
Let's continue with the previous sample. You're building an NFT marketplace. You registered the ERC-721 schema, applied filtering for specific contracts, and you run your chain listener from the genesis block.
And let's say the chain's blockNumber offset is 17600000 now. And you received tens of thousands of ERC-721 Transfer events from filtered contracts and processed them.

And you want to filter a new ERC-721 contract address to that schema, and you need to receive Transfer events belonging to that contract since the blockNumber 13000000. But the blockNumber offset is 17600000, and you haven't received those events because the schema was filtered. What will happen?

If you persist raw transaction logs on the datasource (by setting `DATA_SOURCE_PERSIST_RAW_TRANSACTION_LOGS=true`), you can synchronize your newly filtered contract addresses. However, persisting raw transaction logs in the datasource will increase the disk space used. You can decide whether to persist raw transactions in the datasource, depending on whether you need synchronization.

------------

## Additional Configurations

| Environment Variable  | Description            | Required | Default Value               | Sample Value                                                                                                                                                                                                                                                                                                                                                                                                                       |
|-----------------------|------------------------|----------|-----------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| HTTP_PORT             | HTTP port              | Yes      | 80                          | `80`, `8080`                                                                                                                                                                                                                                                                                                                                                                                                                       |
| HTTP_SWAGGER_ENABLED  | HTTP swagger enabled   | Yes      | true                        | true or false                                                                                                                                                                                                                                                                                                                                                                                                                      |
| LOG_LEVEL             | Minimum log level      | Yes      | info                        | `trace`, `debug`, `info`, `warn`, `error`, `fatal`                                                                                                                                                                                                                                                                                                                                                                                 |
| LOG_TIME_FIELD_FORMAT | Log time field format  | Yes      | `2006-01-02T15:04:05Z07:00` | `01/02 03:04:05PM '06 -0700`, `Mon Jan _2 15:04:05 2006`, `Mon Jan _2 15:04:05 MST 2006`, `Mon Jan 02 15:04:05 -0700 2006`, `02 Jan 06 15:04 MST`, `02 Jan 06 15:04 -0700`, `Monday, 02-Jan-06 15:04:05 MST`, `Mon, 02 Jan 2006 15:04:05 MST`, `Mon, 02 Jan 2006 15:04:05 -0700`, `2006-01-02T15:04:05.999999999Z07:00`, `3:04PM`, `Jan _2 15:04:05`, `Jan _2 15:04:05.000`, `Jan _2 15:04:05.000000`, `Jan _2 15:04:05.000000000` |
| LOG_LEVEL_FIELD_NAME  | Log level field name   | Yes      | level                       | level                                                                                                                                                                                                                                                                                                                                                                                                                              |
| LOG_MESSAGE_FIELD_NAME | Log message field name | Yes      | message                     | message                                                                                                                                                                                                                                                                                                                                                                                                                            |

------------

## Upcoming Improvements
- Add unit and integration tests
- Add MongoDB datasource support
- Add MongoDB destination support
- Add Apache Cassandra destination support