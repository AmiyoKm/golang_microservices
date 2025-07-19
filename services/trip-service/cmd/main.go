package main

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	ctx := context.Background()
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)

	fare := &domain.RideFareModel{
		UserID: "42",
	}
	svc.CreateTrip(ctx, fare)

	for {
		time.Sleep(time.Second)
	}
}
