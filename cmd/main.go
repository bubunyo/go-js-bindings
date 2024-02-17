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

	"github.com/bubunyo/go-js-bindings/server"
	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()

	frontend := server.NewServer(server.Config{
		Dir: "./ui",
	})

	mux.PathPrefix("/ui").Handler(frontend)

	srv := &http.Server{
		Addr:    ":5004",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println("server start failed", err)
			}
		}
	}()

	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	frontend.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
