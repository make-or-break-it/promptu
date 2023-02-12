package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"promptu/apps/feeder-api/internal/storage"
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
	queryParams := r.URL.Query()
	dateQuery := queryParams.Get("date")

	var date time.Time

	if dateQuery == "" {
		date = time.Now()
	} else {
		var err error

		date, err = time.Parse("2006-01-02", dateQuery)

		if err != nil {
			writeError(w, "Date format is not recognised - please use the format YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	feed, err := h.store.GetFeed(context.Background(), date)

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

func writeError(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)

	// TODO: add logging for when things go bad
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
