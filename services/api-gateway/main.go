package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"syscall"
	"time"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /trip/preview", enableCORS(handleTripPreview))
	mux.HandleFunc("/ws/drivers",enableCORS(handleDriversWebSocket))
	mux.HandleFunc("/ws/riders",enableCORS(handleRidersWebsocket))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	go func() {

		log.Printf("Server listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {

	case err := <-serverErrors:
		log.Printf("error stating the server %v", err)
	case sig := <-shutdown:
		log.Printf("server is shutting down due to %s signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("could not shutdown gracefully :%v", err)
			server.Close()
		}
	}
}
