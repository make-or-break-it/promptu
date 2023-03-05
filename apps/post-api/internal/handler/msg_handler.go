package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"promptu/apps/post-api/internal/model"
)

type SyncHandler struct {
	producer sarama.SyncProducer
	topic    string
	logger   *zap.SugaredLogger
}

func NewSyncHandler(producer sarama.SyncProducer, topic string, logger *zap.SugaredLogger) *SyncHandler {
	return &SyncHandler{producer, topic, logger}
}

func (h *SyncHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *SyncHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		writeError(w, "Content type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var post model.Post
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&post); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	} else if post.Message == "" {
		writeError(w, "Post must include a message", http.StatusBadRequest)
		return
	}

	post.UtcCreatedAt = time.Now()
	partition, offset, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: h.topic, // set the topic where all the posts will be stored

		// Key: if left empty the messages will be distributed equally between partitions
		// if topic compaction is enabled, the key has to be set
		// since ordering is only guaranteed within a partition, makes sure that the key guarantees your ordering requirements
		Key: sarama.StringEncoder(post.User),

		Value:     &post,             // the message body
		Timestamp: post.UtcCreatedAt, // Timestamp represents the message creation time
	})

	if err != nil {
		writeError(w, fmt.Errorf("couldn't save message! %s", err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infof("message successfully written to partition %d at offset %d", partition, offset)
	w.WriteHeader(http.StatusCreated)
}

func (h *SyncHandler) Close() error {
	return h.producer.Close()
}

type AsyncHandler struct {
	producer sarama.AsyncProducer
	topic    string
	logger   *zap.SugaredLogger
}

func NewAsyncHandler(producer sarama.AsyncProducer, topic string, logger *zap.SugaredLogger) *AsyncHandler {
	return &AsyncHandler{producer, topic, logger}
}

func (h *AsyncHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *AsyncHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		writeError(w, "Content type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var post model.Post
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&post); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	} else if post.Message == "" {
		writeError(w, "Post must include a message", http.StatusBadRequest)
		return
	}

	post.UtcCreatedAt = time.Now()

	h.producer.Input() <- &sarama.ProducerMessage{
		Topic: h.topic, // set the topic where all the posts will be stored

		// Key: if left empty the messages will be distributed equally between partitions
		// if topic compaction is enabled, the key has to be set
		// since ordering is only guaranteed within a partition, makes sure that the key guarantees your ordering requirements
		Key: sarama.StringEncoder(post.User),

		Value:     &post,             // the message body
		Timestamp: post.UtcCreatedAt, // Timestamp represents the message creation time
	}

	// TODO reminder:
	// there's an error channel for us to read the errors coming back on call-backs
	// there's also a success channel call back
	w.WriteHeader(http.StatusCreated)
}

func (h *AsyncHandler) Close() error {
	return h.producer.Close()
}

func writeError(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)

	// TODO: add logging for when things go bad
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
