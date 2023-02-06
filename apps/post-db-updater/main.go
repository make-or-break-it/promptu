package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/gorilla/mux"

	"promptu/apps/post-db-updater/internal/config"
	"promptu/apps/post-db-updater/internal/handler"
	"promptu/apps/post-db-updater/internal/storage"

	"github.com/Shopify/sarama"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {

			ctx := context.Background()
			cfg := config.AppConfig() // todo change

			store := storage.NewMongoDbStore("promptu-db")
			consumer, err := sarama.NewConsumerGroup(strings.Split(cfg.KafkaBrokers, ","), "post-producer", createConsumerGroupConfig(cfg.KafkaVersion))
			if err != nil {
				log.Panicf("Error creating consumer group client: %v", err)
			}

			go handler.Consume(ctx, handler.NewMsgHandler(store), cfg.PostTopic, consumer)
			return runHTTPServer(ctx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("run failed: %s", err.Error())
	}
}

func runHTTPServer(ctx context.Context) error {

	r := createRouter()
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

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handler.Health).Methods(http.MethodGet, http.MethodOptions)

	return r
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, r)
	})
}

func createConsumerGroupConfig(version string) *sarama.Config {
	v, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = v
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	// sticky: []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	// range: []sarama.BalanceStrategy{sarama.BalanceStrategyRange}

	return config
}
