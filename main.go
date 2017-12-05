package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"encoding/binary"
	"math"

	"cloud.google.com/go/pubsub"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
	"github.com/zenoss/zing-proto/go/query"

	"cloud.google.com/go/bigtable"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"net/http"
	"encoding/json"
)

var (
	projectID              string = "zing-dev"
	btInstanceID           string = "zing-dev-bt1"
	tableName              string = "metrics"
	cfName                 string = "METRIC"
	inputTopic             string = "query-in"
	inputTopicSubscription string = "query-in-sub-query"
	outputTopic            string = "query-out"
	subscription           string = "query-out-gosub"
	metricID               string = "mymetricid"
	start                  int64  = time.Now().Add(-1 * time.Hour).Unix()
	end                    int64  = time.Now().Add(-10 * time.Minute).Unix()
	newTopicSetByUser      bool   = false
)

func main() {

	app := cli.App("gcloudutils", "Utils for Zing GCP data apps")
	app.Version("v version", "gcloudutils 0.0.1")

	projectIsSetByUser := false

	projectIDOpt := app.String(cli.StringOpt{
		Name:      "p projectid",
		Desc:      "The Id of the GCP project",
		Value:     "zing-dev",
		EnvVar:    "GCP_PROJECT_ID",
		SetByUser: &projectIsSetByUser,
	})

	app.Spec = "[-p|--projectid]"

	app.Command("create-topic", "Create pubsub topic", func(cmd *cli.Cmd) {
		pubsubtopic := cmd.String(cli.StringArg{
			Name:      "TOPIC",
			Desc:      "The name of the topic to be created",
			Value:     "",
			SetByUser: &newTopicSetByUser,
		})

		cmd.Spec = "TOPIC"
		cmd.Action = func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			logrus.Info("Context created")
			newTopic, err := createNewTopic(ctx, *projectIDOpt, *pubsubtopic)
			if err != nil {
				logrus.WithError(err).Fatal("Failed to create pubsub topic")
			}
			logrus.Info("Pubsub topic created: ", newTopic.String())
		}

	})

	app.Command("create-subscription", "Create pubsub subscription", func(cmd *cli.Cmd) {
		pubsubtopic := cmd.String(cli.StringArg{
			Name:  "TOPIC",
			Desc:  "The name of the topic to subscribe to",
			Value: "",
		})

		topicSubscriptionArg := cmd.String(cli.StringArg{
			Name:  "SUBSCRIPTION",
			Desc:  "The ID of the subscription",
			Value: "",
		})

		cmd.Spec = "TOPIC SUBSCRIPTION"
		cmd.Action = func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			logrus.Info("Context created")
			newSubscription, err := createNewSubscription(ctx, *projectIDOpt, *pubsubtopic, *topicSubscriptionArg)
			if err != nil {
				logrus.WithError(err).Fatal("Failed to create pubsub topic")
			}
			logrus.Info("Pubsub subscription created: ", newSubscription.String())
		}

	})

	app.Command("bigtable", "Bigtable actions", bigTableCmd(*projectIDOpt))
	app.Command("pubsub", "Pubsub actions", pubsubCmd(*projectIDOpt))
	app.Command("query", "Domain query actions", domainQueryCmd(*projectIDOpt))
	app.Command("datastack-prereqs", "Data-stack prerequisites actions", serveDatastackInitCheckerCmd(*projectIDOpt))
	app.Run(os.Args)
}

func createNewTopic(ctx context.Context, projectID, topicName string) (*pubsub.Topic, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pubsub client")
	}

	// input topic
	topic := client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check if input topic exists")
	}

	if !exists {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create input topic")
		}
	}
	return topic, nil
}

func getTopics(ctx context.Context, projectID string) (*pubsub.TopicIterator, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pubsub client")
	}

	return client.Topics(ctx), nil
}

func getSubscriptions(ctx context.Context, projectID string) (*pubsub.SubscriptionIterator, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pubsub client")
	}

	return client.Subscriptions(ctx), nil
}

func createNewSubscription(ctx context.Context, projectID, topicName, subscriptionName string) (*pubsub.Subscription, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pubsub client")
	}

	sub := client.Subscription(subscriptionName)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check if subscription to input topic exists")
	}

	topic := client.Topic(topicName)
	topicExists, topicErr := topic.Exists(ctx)
	if topicErr != nil {
		return nil, errors.Wrap(topicErr, "Failed to check if topic exists")
	}
	if !topicExists {
		return nil, fmt.Errorf("Failed to create subscription, topic %s does not exist", topic.String())
	}

	if !exists {
		sub, err = client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 30 * time.Second,
		})
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create input topic subscription")
	}
	return sub, nil
}

func initPubSub(ctx context.Context) error {
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

func pubsubCmd(projectId string) func(*cli.Cmd) {

	return func(cmd *cli.Cmd) {

		cmd.Command("topics", "Pubsub topics commands", func(cmd *cli.Cmd) {
			cmd.Command("list", "List PubSub topics", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					logrus.Info("Context created")
					topics, err := getTopics(ctx, projectId)
					if err != nil {
						logrus.Error(err)
					}
					for {
						t, err := topics.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							logrus.Error(err)
						}
						fmt.Println(t)
					}

				}
			})
		})

		cmd.Command("subscriptions", "Pubsub subscriptions commands", func(cmd *cli.Cmd) {
			cmd.Command("list", "List Pubsub subscriptions", func(cmd *cli.Cmd) {
				cmd.Action = func() {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					logrus.Info("Context created")
					subscriptions, err := getSubscriptions(ctx, projectId)
					if err != nil {
						logrus.Error(err)
					}
					for {
						t, err := subscriptions.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							logrus.Error(err)
						}
						fmt.Println(t)
					}
				}
			})
		})
	}
}

func bigTableCmd(projectId string) func(*cli.Cmd) {

	return func(cmd *cli.Cmd) {
		var (
			btInstanceID = cmd.StringOpt("i bigtable-instanceid", "", "Id of the Bigtable Instance")
			btTableName  = cmd.StringOpt("t bigtable-table", "", "Name of the BigTable")
			cfName       = cmd.StringOpt("c bigtable-columnfamily", "", "Name of the Bigtable Column Family")
		)

		cmd.Spec = "-i -t -c"
		cmd.Command("init", "Initialize Bigtable schema", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				logrus.Info("Context created")
				logrus.Infof("Initializing BigTable for project: %s with bigtableInstance: %s bigtableTable: %s columnFamily: %s",
					projectId, *btInstanceID, *btTableName, *cfName)
				initBigTable(ctx, projectId, *btInstanceID, *btTableName, *cfName)
			}
		})
	}

}

func domainQueryCmd(projectId string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("test", "Publish Domain Query request and Consume results", func(cmd *cli.Cmd) {
			subscriptionOpt := cmd.StringOpt("s subscription", "", "PubSub subscription to consume from")
			cmd.Action = func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				go readResults(ctx, projectId, *subscriptionOpt)
				time.Sleep(2 * time.Second)
				logrus.Info("Result listener started")
				pushMessages(ctx, projectId)
				logrus.Info("Query messages pushed")
				time.Sleep(30 * time.Second)
			}
		})
	}
}

func initBigTable(ctx context.Context, projectID, btInstanceID, tableName, cfName string) error {
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
	for ts := start - 900; ts <= end+900; ts += 300 {
		rowkey := fmt.Sprintf("%s#%d", metricID, ts)
		rowkeys = append(rowkeys, rowkey)
		mut := bigtable.NewMutation()
		mut.Set(cfName, "NAME", bigtable.ServerTime, []byte("cpuUtil_cpuUtil"))
		value := rand.Float64() * 100.0
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
			} else if readItem.Column == "METRIC:NAME" {
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

func pushMessages(ctx context.Context, projectID string) {
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

func readResults(ctx context.Context, projectID, subscription string) {
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

type DataStackReqs struct {
	Topics         []string `json:topics`
	Subscriptions  []string `json:subscriptions`
	BigTableTables []string `json:tables`
}

func zingDataStackInitCheckHandler(w http.ResponseWriter, r *http.Request) {
	var appReqs DataStackReqs
	json.NewDecoder(r.Body).Decode(&appReqs)

	fmt.Fprintf(w, "nope")
}

func serveDatastackInitCheckerCmd(projectId string) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Command("server", "A web server for checking if data stack is initialized", func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				http.HandleFunc("/isstackup", zingDataStackInitCheckHandler)
				logrus.Info("Starting server")
				http.ListenAndServe(":8199", nil)
				logrus.Info("Stopping server")
			}
		})
	}
}
