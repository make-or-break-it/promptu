package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"promptu/api/internal/model"
	"promptu/api/internal/storage"
)

type Handler struct {
	store storage.Store
}

func NewHandler(store storage.Store) *Handler {
	return &Handler{store}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	feed, err := h.store.GetFeed(context.Background())

	switch {
	case err != nil:
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		jsonResp, err := json.MarshalIndent(feed, "", " ")

		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
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
		writeError(w, "message is missing", http.StatusBadRequest)
		return
	}

	createdAt := time.Now()
	if err := h.store.PostMessage(context.Background(), post, createdAt); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
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
