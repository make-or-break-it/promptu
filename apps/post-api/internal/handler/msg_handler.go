package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/sarama"

	"promptu/apps/post-api/internal/model"
)

type Handler struct {
	producer sarama.SyncProducer
	topic    string
}

func NewHandler(producer sarama.SyncProducer, topic string) *Handler {
	return &Handler{producer, topic}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	// Needed to support CORS: https://flaviocopes.com/golang-enable-cors/#:~:text=Handling%20pre%2Dflight%20OPTIONS%20requests
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
	_, _, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     h.topic,
		Key:       nil, // if left empty the messages will be distributed equally between partitions
		Value:     &post,
		Timestamp: post.UtcCreatedAt,
	})

	if err != nil {
		writeError(w, fmt.Errorf("couldn't save message! %s", err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
