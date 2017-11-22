package main

import (
	"fmt"
	"time"
	"math/rand"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/proto"
	"github.com/johnlhamilton/pubsubtester/query"
	"github.com/sirupsen/logrus"
	//"github.com/zenoss/zing-proto/go/query"
	"golang.org/x/net/context"
	"cloud.google.com/go/bigtable"
	"encoding/binary"
	"math"
	"github.com/pkg/errors"
)

var (
	projectID string = "zenoss-zing"
	btInstanceID string = "zenoss-zing-bt1"
	tableName string = "metrics"
	cfName string = "METRIC"
	inputTopic string = "query-in"
	inputTopicSubscription string = "query-in-sub-query"
	outputTopic string = "query-out"
	subscription string = "query-out-gosub"
	metricID string = "mymetricid"
	start int64 = time.Now().Add(-1*time.Hour).Unix()
	end int64 = time.Now().Add(-10*time.Minute).Unix()
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logrus.Info("Context created")
	err := initPubSub(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize pubsub")
	}
	logrus.Info("Pubsub initialized")
	err = initBigTable(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize bigtable")
	}
	logrus.Info("Big Table Initialized")

	err = readBigTable(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read bigtable")
	}

	go readResults(ctx)
	time.Sleep(2 * time.Second)
	logrus.Info("Result listener started")
	pushMessages(ctx)
	logrus.Info("Query messages pushed")

	time.Sleep(30 * time.Second)
}

func initPubSub(ctx context.Context) error{
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "Failed to create pubsub client")
	}

	// input topic
	topic := client.Topic(inputTopic)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if input topic exists")
	}

	if !exists {
		topic, err = client.CreateTopic(ctx, inputTopic)
		if err != nil {
			return errors.Wrap(err, "Failed to create input topic")
		}
	}

	sub := client.Subscription(inputTopicSubscription)
	exists, err = sub.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if subscription to input topic exists")
	}
	if !exists {
		sub, err = client.CreateSubscription(ctx, inputTopicSubscription, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 30 * time.Second,
		})
	}

	if err != nil {
		return errors.Wrap(err, "Failed to create input topic subscription")
	}

	// output topic
	topic = client.Topic(outputTopic)
	exists, err = topic.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if output topic exists")
	}

	if !exists {
		topic, err = client.CreateTopic(ctx, outputTopic)
		if err != nil {
			return errors.Wrap(err, "Failed to create output topic")
		}
	}

	sub = client.Subscription(subscription)
	exists, err = sub.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to check if subscription to output topic exists")
	}
	if !exists {
		sub, err = client.CreateSubscription(ctx, subscription, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 30 * time.Second,
		})
	}

	if err != nil {
		return errors.Wrap(err, "Failed to create output topic subscription")
	}

	return nil
}

func initBigTable(ctx context.Context) error{
	aclient, err := bigtable.NewAdminClient(ctx, projectID, btInstanceID)
	if err != nil {
		return err
	}
	defer aclient.Close()

	err = aclient.CreateTable(ctx, tableName)
	if err != nil {
		logrus.WithError(err).Error("Failed to create table, it probably already exists or something")
		return nil
	}

	err = aclient.CreateColumnFamily(ctx, tableName, cfName)
	if err != nil {
		return err
	}

	client, err := bigtable.NewClient(ctx, projectID, btInstanceID)
	if err != nil {
		return err
	}

	table := client.Open(tableName)
	defer client.Close()

	keys, mutations := generateMetrics()
	errs, err := table.ApplyBulk(ctx, keys, mutations)
	if err != nil {
		return err
	}

	if errs != nil && len(errs) > 0 {
		for i, e := range errs {
			logrus.WithError(e).WithField("rowkey", keys[i]).Error("Individual row insert failed")
		}
		return errors.New("Error inserting individual row")
	}

	return nil
}

func generateMetrics() ([]string, []*bigtable.Mutation) {
	rand.Seed(time.Now().Unix())
	rowkeys := []string{}
	mutations := []*bigtable.Mutation{}
	for ts := start - 900; ts <= end + 900; ts += 300 {
		rowkey := fmt.Sprintf("%s#%d", metricID, ts)
		rowkeys = append(rowkeys, rowkey)
		mut := bigtable.NewMutation()
		mut.Set(cfName, "NAME", bigtable.ServerTime, []byte("cpuUtil_cpuUtil"))
		value := rand.Float64()*100.0
		valueBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(valueBytes[:], math.Float64bits(value))
		mut.Set(cfName, "VALUE", bigtable.ServerTime, valueBytes)
		mutations = append(mutations, mut)
	}

	return rowkeys, mutations
}

func readBigTable(ctx context.Context) error {
	client, err := bigtable.NewClient(ctx, projectID, btInstanceID)
	if err != nil {
		return err
	}

	table := client.Open(tableName)
	defer client.Close()

	err = table.ReadRows(ctx, bigtable.InfiniteRange(metricID), func(r bigtable.Row) bool {
		fmt.Println(r)
		rowkey := r.Key()
		var name string
		var value float64

		for _, readItem := range r["METRIC"] {
			if readItem.Column == "METRIC:VALUE" {
				value = math.Float64frombits(binary.LittleEndian.Uint64(readItem.Value))
			}else if readItem.Column == "METRIC:NAME" {
				name = string(readItem.Value)
			}
		}

		fmt.Println(rowkey)
		fmt.Printf("  NAME:  %s\n", name)
		fmt.Printf("  VALUE:  %v\n\n", value)
		return true
	})

	return err
}

func pushMessages(ctx context.Context) {
	query1 := query.Query{
		Id: "query-metric-1",
		QueryType: &query.Query_MetricQuery_{
			MetricQuery: &query.Query_MetricQuery{
				MetricId: metricID,
				Start:    start,
				End:      end,
			},
		},
	}

	logrus.WithField("query", query1.String()).Info("Created query")

	data, err := proto.Marshal(&query1)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to marshal query protobuf")
	}

	logrus.WithField("data", string(data)).Info("Marshaled query")

	var query2 query.Query
	proto.Unmarshal(data, &query2)
	logrus.WithField("unmarshalledquery", query2.String()).Info("Unmarshaled marshaled query")

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to create client for publishing")
	}

	t := client.Topic(inputTopic)

	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to publish message")
	}
	logrus.WithField("messageid", id).Info("Published a message")
}

func readResults(ctx context.Context) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to create client for reading")
	}

	sub := client.Subscription(subscription)

	if err != nil {
		logrus.WithError(err).Fatalf("Failed to create subscription to output topic")
	}

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var resultReceived query.Result
		proto.Unmarshal(msg.Data, &resultReceived)
		logrus.WithField("messageraw", string(msg.Data)).WithField("messagresult", resultReceived.String()).Info("Got result!")
		msg.Ack()
	})

	if err != nil {
		logrus.WithError(err).Error("Got error on sub.Receive")
	}
}