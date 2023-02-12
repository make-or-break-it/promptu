package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"promptu/apps/post-db-updater/internal/model"
	"promptu/apps/post-db-updater/internal/storage"

	"github.com/Shopify/sarama"
)

type MsgHandler struct {
	store storage.Store
}

func NewMsgHandler(store storage.Store) *MsgHandler {
	return &MsgHandler{store}
}

func Consume(ctx context.Context, msgHandler *MsgHandler, topic string, client sarama.ConsumerGroup) {
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := client.Consume(ctx, []string{topic}, msgHandler); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			log.Printf("Error from ctx: %v", ctx.Err())
			return
		}
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (handler *MsgHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (handler *MsgHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler *MsgHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
			var post model.Post
			if err := json.Unmarshal(message.Value, &post); err != nil {
				return fmt.Errorf("error while unmarshalling post: %w", err)
			}
			if err := handler.store.PostMessage(context.Background(), post); err != nil {
				return fmt.Errorf("error while posting message to the db: %w", err)
			}

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
