package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// NewSyncProducer creates synchronous kafka producer - offering consistent semantics
func NewSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start Sarama producer: %w", err)
	}

	return producer, nil
}

// NewAsyncProducer creates an asynchronous kafka producer - offering high throughput
func NewAsyncProducer(brokerList string) (sarama.AsyncProducer, error) {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(strings.Split(brokerList, ","), config)
	if err != nil {
		return nil, fmt.Errorf("failed to start Sarama producer: %w", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			fmt.Printf("Failed to write access log entry: %v\n", err)
		}
	}()

	return producer, nil
}
