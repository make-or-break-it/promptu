package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	cli "github.com/urfave/cli/v2"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	hndlr "promptu/apps/post-api/internal/handler"
)

var brokers, topic string

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error while iniatilising logger %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugarLogger := logger.Sugar()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "brokers",
				Usage:       "The address of the kafka brokers",
				Destination: &brokers,
				Value:       os.Getenv("KAFKA_HOST"),
			},
			&cli.StringFlag{
				Name:        "topic",
				Usage:       "The name of the topic that the producer sends messages to",
				Destination: &topic,
				Value:       os.Getenv("KAFKA_TOPIC"),
			},
		},
		Action: func(c *cli.Context) error {
			brokerList := strings.Split(brokers, ",")

			err = ensureTopicExists(sugarLogger, brokerList, topic)
			if err != nil {
				return err
			}

			producer, err := hndlr.NewSyncProducer(brokerList)
			if err != nil {
				return err
			}
			handler := hndlr.NewSyncHandler(producer, topic, sugarLogger)

			// aProducer, err := hndlr.NewAsyncProducer(brokerList)
			// if err != nil {
			// 	return err
			// }
			// handler := hndlr.NewAsyncHandler(aProducer, topic, sugarLogger)
			return runHTTPServer(context.Background(), handler, sugarLogger)
		},
	}

	if err := app.Run(os.Args); err != nil {
		sugarLogger.Fatalf("run failed: %s", err.Error())
	}
}

func runHTTPServer(ctx context.Context, handler handler, logger *zap.SugaredLogger) error {
	r, err := createRouter(handler, logger)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	logger.Info("Starting the server on port 8080")
	srv := &http.Server{
		Addr:         getAddress(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      middleware(r),
	}
	defer srv.Shutdown(ctx)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		return errors.New("shutting down due to received signal")
	case err := <-errCh:
		if err := handler.Close(); err != nil {
			logger.Infof("failed to shut down producer cleanly %v", err)
		}
		return fmt.Errorf("errors received on channel %v", err)
	}
}

func getAddress() string {
	var (
		host, port string
		exists     bool
	)

	host, exists = os.LookupEnv("HOST")
	if !exists {
		host = defaultHost
	}

	port, exists = os.LookupEnv("PORT")
	if !exists {
		port = defaultPort
	}
	return host + ":" + port
}

func createRouter(hndlr handler, logger *zap.SugaredLogger) (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/post", hndlr.PostMessage).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/health", hndlr.Health).Methods(http.MethodGet, http.MethodOptions)

	return r, nil
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, r)
	})
}

type handler interface {
	Health(w http.ResponseWriter, r *http.Request)
	PostMessage(w http.ResponseWriter, r *http.Request)
	Close() error
}

func ensureTopicExists(logger *zap.SugaredLogger, brokers []string, topic string) error {
	// create topic if doesn't exist
	cfg := sarama.NewConfig()
	var err error
	cfg.Version, err = sarama.ParseKafkaVersion("2.8.1")
	cfg.Metadata.Retry.Max = 5
	cfg.Metadata.Retry.Backoff = 10 * time.Second
	cfg.ClientID = "sarama-prepareTestTopics"
	cl, err := sarama.NewClient(brokers, cfg)
	cr, err := cl.Controller()
	defer cr.Close()

	response, err := cr.CreateTopics(createTopic(logger, topic))
	if err != nil {
		return fmt.Errorf("topic creation failed with %w", err)
	}

	logger.Infof("topic creation response: %v \n", response)
	return nil
}

func createTopic(logger *zap.SugaredLogger, topic string) *sarama.CreateTopicsRequest {
	topicCleanUpPolicy := "delete"
	retentionMS := "-1"
	tenMB := "10485760"
	compressionMethod := "snappy"

	req := &sarama.CreateTopicsRequest{
		TopicDetails: map[string]*sarama.TopicDetail{
			topic: {
				NumPartitions:     10,
				ReplicationFactor: 1,
				// ReplicaAssignment: map[int32][]int32{
				// 	0: {0, 1, 2},
				// }, - can't be both ReplicationFactor & ReplicaAssignment specified at the same time

				ConfigEntries: map[string]*string{ // topic properties following kafka doc
					"cleanup.policy": &topicCleanUpPolicy, // default: "delete'
					// "delete": discards old segments when their retention time or size limit has been reached
					// "compact": enables log compaction, which retains the latest value for each key
					// can apply both

					"retention.ms": &retentionMS, // default: 604800000=7days
					// maximum time we will retain a log before discarding old log segments to free up space
					// if we are using the "delete" retention policy.

					//"retention.bytes": ,// default: -1
					// controls the maximum size a partition (which consists of log segments) can grow to before we discard old log segments to free up space
					// if we are using the "delete" retention policy
					// By default there is no size limit only a time limit.
					// Since this limit is enforced at the partition level, multiply it by the number of partitions to compute the topic retention in bytes.

					// "segment.bytes": ,// default: 1 GB - controls the segment file size for the log.
					// Retention and cleaning is always done a file at a time so a larger segment size means fewer files
					// but less granular control over retention.

					"max.message.bytes": &tenMB, // default: 1048588 ~ 1MB

					"compression.type": &compressionMethod, //default: uncompressed

					// min.insync.replicas - creating a topic with a replication factor of 3 & min.insync.replicas to 2 & produce with acks of "all"
					// => ensures that the producer raises an exception if a majority of replicas do not receive a write
				},
			},
		},
		Timeout: 100 * time.Millisecond,
	}
	req.Version = 1
	//req.ValidateOnly = true
	return req
}
