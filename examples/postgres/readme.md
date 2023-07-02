# PostgreSQL Destination

You should set `DESTINATION=postgres`

|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_POSTGRES_CONNECTION_STRING | The connection string of the destination database  | **Yes** | - | `postgres://postgresuser:password@localhost:5432/sourcedb` |
| DESTINATION_POSTGRES_PERSIST_RAW_TRANSACTION_LOGS | Set this value to true if you want to persist the raw transaction logs in your destination database | No | false | true or false |

