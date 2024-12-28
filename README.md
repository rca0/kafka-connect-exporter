# kafka-connect-exporter

Based on Kafka Connect Confluent REST API, this application get sink and source information and expose metrics
- https://docs.confluent.io/platform/current/connect/references/restapi.html

## Parameters

```bash
Usage of ./kafka_connect_exporter:
  -listen-address string
    	Address on which to expose metrics (default ":8000")
  -scrape-uri string
    	URI on which to scrape Kafka connect. (default "http://127.0.0.1:8000,http://127.0.0.1:8003,http://127.0.0.1:8080")
  -telemetry-path string
    	Path under which to expose metrics (default "/metrics")
  -version
    	show version and exit
```

## USAGE

```bash
./kafka_connect_exporter \
    -listen-address ":7070" \
    -scrape-uri "http://127.0.0.1:8083,http://127.0.0.2:8083,http://127.0.0.3:8083" \
    -telemetry-path "/metrics"
```

## Metrics

```bash
# TYPE kafka_connect_connector_state_running gauge
# sink example
kafka_connect_connector_state_running{bucket_name="bucket_name10010",connector="s3-sink-foo",consumer_group="connect-s3-sink-foo",state="running",worker_id="b-1.domain.com:8083",topic="sink-topic-foo-name"} 1
kafka_connect_connector_state_running{bucket_name="bucket_name_bar-1",connector="s3-sink-bar",consumer_group="connect-s3-sink-bar",state="running",worker_id="b-2.domain.com:8083",topic="sink-topic-bar-name"} 1
# source example
kafka_connect_connector_state_running{connector="source-foo",state="running",worker_id="b-2.domain.com:8083",topic="source-topic_name",source_database_host="db-01.domain.com",source_database_name="db_name_0101"} 1
# HELP kafka_connect_connector_tasks_state the state of tasks. 0-failed, 1-running, 2-unassigned, 3-paused
# TYPE kafka_connect_connector_tasks_state gauge
# sink example
kafka_connect_connector_tasks_state{bucket_name="bucket_name_bar-1",connector="s3-sink-bar",consumer_group="connect-s3-sink-bar",state="running",worker_id="b-2.domain.com:8083",id="0",topic="bar-topic-name"} 1
kafka_connect_connector_tasks_state{bucket_name="bucket_name10010",connector="s3-sink-foo",consumer_group="connect-s3-sink-foo",state="failed",worker_id="b-2.domain.com:8083",id="1",topic="foo-topic-name"} 1
# source example
kafka_connect_connector_tasks_state{connector="source-foo",state="running",worker_id="b-2.domain.com:8083",id="1",topic="foo-topic-name",source_database_host="db-01.domain.com",source_database_name="db_name_0101"} 1
# HELP kafka_connect_connector_count number of deployed connectors
# TYPE kafka_connect_connector_count gauge
kafka_connect_connector_count 2
# HELP kafka_connect_up  was the last scrape of kafka connect successful?
# TYPE kafka_connect_up gauge
kafka_connect_up 1
```