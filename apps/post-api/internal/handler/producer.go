package handler

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// NewSyncProducer creates synchronous kafka producer - offering consistent semantics
/*
	Notes.
	The trade-off here is performance.
	-> brokers depending on how busy they are can take from 2ms to few seconds to reply
	-> during this time the sending thread is waiting, not even sending new messages

	This method is rarely used in production.
*/
func NewSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	// For data collection, we might be looking for synchrounous producer

	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	// With flushing some batching mechanism will apply before sending the message.

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message

	config.Producer.Retry.Max = 10 // Retry up to 10 times to produce the message -
	// Retryable errors will be retried. Examples for such: "connection errors", "no leader for partition".
	// Non-retryable errors can't be automatically fixed, and won't be retried, for example: "message size too large"

	config.Producer.Return.Successes = true

	// On the broker side, you can confirm the following settings to get
	// stronger consistency guarantees:
	// - `unclean.leader.election.enable` should be false

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start Sarama producer: %w", err)
	}

	return producer, nil
}

// NewAsyncProducer creates an asynchronous kafka producer - offering high throughput
/*
	Notes.
	In most cases we really don't need a reply from the producer (and sending back the partition and offset is not useful, just slows the producer down).
	However, we do want to know if a message failed to be written.
	For these cases, the producer supports adding a callback when sending a record.

	This method is usually preferred in production due to it's high throughput!
*/
func NewAsyncProducer(brokerList []string) (sarama.AsyncProducer, error) {
	// For metrics or processing logs, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start Sarama producer: %w", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			// In a production application we would want to make some changes to this code:
			// -> handle the error in a safe way - send it to a deadletter queue, use structure logging, make it visible in log aggregation system, create alert on its appearance
			// -> this goroutine will "hang in the air", and won't be correctly synchronised with the other goroutines, in case of failure we might loose error messages
			fmt.Printf("Failed to write posts to the queue: %v\n", err)
		}
	}()

	return producer, nil
}
