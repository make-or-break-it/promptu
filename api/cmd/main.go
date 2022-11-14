package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/gorilla/mux"

	"promptu/api/internal/handler"
	"promptu/api/internal/storage"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			store := storage.NewInMemoryStore()
			return run(context.Background(), store)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("run failed: %s", err.Error())
	}
}

func run(ctx context.Context, store handler.Store) error {

	r := createRouter(store)
	errCh := make(chan error, 1)

	log.Print("Starting the server on port 8080")
	srv := &http.Server{
		Addr:         "localhost:8080",
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

func createRouter(store handler.Store) *mux.Router {
	hndlr := handler.NewHandler(store)

	r := mux.NewRouter()
	r.HandleFunc("/feed", hndlr.GetFeed).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/post", hndlr.PostMessage).Methods(http.MethodPost, http.MethodOptions)

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
