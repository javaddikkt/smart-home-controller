package main

import (
	"context"
	"errors"
	httpGatewayModels "homework/internal/gateways/http/types"
	"homework/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	httpGateway "homework/internal/gateways/http"
	eventRepository "homework/internal/repository/event/inmemory"
	sensorRepository "homework/internal/repository/sensor/inmemory"
	userRepository "homework/internal/repository/user/inmemory"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	er := eventRepository.NewEventRepository()
	sr := sensorRepository.NewSensorRepository()
	ur := userRepository.NewUserRepository()
	sor := userRepository.NewSensorOwnerRepository()

	useCases := httpGatewayModels.UseCases{
		Event:  usecase.NewEvent(er, sr),
		Sensor: usecase.NewSensor(sr),
		User:   usecase.NewUser(ur, sor, sr),
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt < 1 || portInt > 65535 {
		log.Fatal("port must be between 1 and 65535")
	}


	r := httpGateway.NewServer(useCases, httpGateway.WithHost(host), httpGateway.WithPort(uint16(portInt)))

	if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}
