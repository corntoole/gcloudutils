#!/usr/bin/env bash

echo "dataflow-query initialization begin"

export GCP_PROJECT_ID="${GCP_PROJECT_ID:-zing-dev}"
echo "GCP_PROJECT_ID $GCP_PROJECT_ID"
./gcloudutils create-topic query-in
./gcloudutils create-subscription query-in query-in-sub-query
./gcloudutils create-topic query-out
./gcloudutils create-subscription query-out query-out-gosub
./gcloudutils pubsub topics list
./gcloudutils pubsub subscriptions list
./gcloudutils bigtable -i zing-dev-bt1 -t metrics -c METRIC init

echo "dataflow-query initialization done"

