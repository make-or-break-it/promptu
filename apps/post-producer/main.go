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

	"promptu/apps/post-producer/internal/handler"
)

var brokers, topic string

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "brokers",
				Usage:       "The address of the kafka brokers",
				Destination: &brokers,
				Value:       "149.248.217.129:9092",
			},
			&cli.StringFlag{
				Name:        "topic",
				Usage:       "The name of the topic that the producer sends messages to",
				Destination: &topic,
				Value:       "posts",
			},
		},
		Action: func(c *cli.Context) error {
			brokerList := strings.Split(brokers, ",")

			// TODO uncomment and create topics with specific config if missing
			// instead of defaults

			// err := ensureTopicExists(brokerList, topic)
			// if err != nil {
			// 	return err
			// }

			producer, err := handler.NewSyncProducer(brokerList)
			if err != nil {
				return err
			}
			return runHTTPServer(context.Background(), producer)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("run failed: %s", err.Error())
	}
}

func runHTTPServer(ctx context.Context, producer sarama.SyncProducer) error {

	r, err := createRouter(producer)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	log.Print("Starting the server on port 8080")
	srv := &http.Server{
		Addr:         getAddress(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      middleware(r),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()
	defer srv.Shutdown(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigChan:
		return errors.New("shutting down due to received signal")
	case err := <-errCh:
		if err := producer.Close(); err != nil {
			log.Println("Failed to shut down access log producer cleanly", err)
		}
		return err
	}
}

func getAddress() string {
	var host, port string
	var exists bool

	host, exists = os.LookupEnv("HOST")
	if !exists {
		host = "0.0.0.0"
	}

	port, exists = os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	return host + ":" + port
}

func createRouter(producer sarama.SyncProducer) (*mux.Router, error) {
	hndlr := handler.NewHandler(producer, topic)
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

func ensureTopicExists(brokers []string, topic string) error {
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
	response, err := cr.CreateTopics(createTopic(topic))
	fmt.Printf("response: %v, error: %v \n", response, err)

	if err != nil {
		return err
	}
	return nil
}

func createTopic(topic string) *sarama.CreateTopicsRequest {
	retention := "-1"
	req := &sarama.CreateTopicsRequest{
		TopicDetails: map[string]*sarama.TopicDetail{
			topic: {
				NumPartitions:     10,
				ReplicationFactor: 1,
				// ReplicaAssignment: map[int32][]int32{
				// 	0: {0, 1, 2},
				// },
				ConfigEntries: map[string]*string{
					"retention.ms": &retention,
				},
			},
		},
		Timeout: 100 * time.Millisecond,
	}
	req.Version = 1
	req.ValidateOnly = true
	return req
}
