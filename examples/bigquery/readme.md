# BigQuery Destination

You should set `DESTINATION=bigquery`

|  Environment Variable  | Description   |  Required  |  Default Value  | Sample Value  |
| --- | --- | --- | --- | --- |
| DESTINATION_BIGQUERY_PROJECT_ID | BigQuery Project ID | **Yes** | - | xyz-123 |
| DESTINATION_BIGQUERY_DATASET | BigQuery Data Set | **Yes** | - | chain |
| GOOGLE_APPLICATION_CREDENTIALS | The path of your Google Cloud Service Account Key json file  | **Yes** | - | /Users/username/service_account.json |
| DESTINATION_BIGQUERY_DELIVERY_GUARANTEE | The delivery guarantee option can only be AT_LEAST_ONCE or AT_MOST_ONCE for BigQuery. | No | `AT_LEAST_ONCE` | `AT_LEAST_ONCE` or `AT_MOST_ONCE` |
| DESTINATION_BIGQUERY_PERSIST_RAW_TRANSACTION_LOGS | Set this value to true if you want to persist the raw transaction logs in the BigQuery | No | false | true or false |
