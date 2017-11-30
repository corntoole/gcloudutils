# gcloudutils

## Usage


To create a topic and subscription
```
gcloudutils create-topic metric-raw-rm
gcloudutils create-subscription metric-raw-rm metric-raw-rm-sub-ingest
```

To initialize and populate BigTable with metrics
```
gcloudutils bigtable init
```

To test the domain-query dataflow-pipeline:
```
gcloudutils query test
```
