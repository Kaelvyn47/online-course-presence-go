package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.CreateChannel("course-presence"); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.DeleteChannel("course-presence"); err != nil {
			log.Printf("delete course-presence: %v", err)
		}
	}()
	s := &service{client: client, learners: []Learner{{ID: "lina", Name: "Lina", Deadline: time.Now().Add(24 * time.Hour)}}}
	server := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(s.handler)}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.ListenAndServe() }()
	log.Println("listening on http://localhost:8080/course/status")
	select {
	case err := <-serveErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown: %v", err)
		}
	}
}
