package main 

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)

	httpHandler := h.HttpHandler{Service: svc}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}
	serverError := make(chan error, 1)
	go func() {
		log.Printf("server started at %s", server.Addr)
		serverError <- server.ListenAndServe()
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, os.Interrupt)
	
	select {
	case err := <-serverError:
		log.Printf("something went wrong when starting the server : %v", err)
	case sig := <-shutdown:
		log.Printf("server is shutting down due to %s signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("could not shutdown gracefully :%v", err)
			server.Close()
		}
	}
}
