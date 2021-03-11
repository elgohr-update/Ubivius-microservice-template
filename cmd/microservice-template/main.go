package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Ubivius/microservice-template/pkg/database"
	"github.com/Ubivius/microservice-template/pkg/handlers"
	"github.com/Ubivius/microservice-template/pkg/router"
)

func main() {
	// Logger
	logger := log.New(os.Stdout, "Template", log.LstdFlags)

	// Database init
	db := database.NewMongoProducts()

	// Creating handlers
	productHandler := handlers.NewProductsHandler(logger, db)

	// Router setup
	r := router.New(productHandler, logger)

	// Server setup
	server := &http.Server{
		Addr:        ":9090",
		Handler:     r,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Starting server on port ", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			logger.Println("Server error : ", err)
			logger.Fatal(err)
		}
	}()

	// Handle shutdown signals from operating system
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	receivedSignal := <-signalChannel

	logger.Println("Received terminate, beginning graceful shutdown", receivedSignal)

	// DB connection shutdown
	db.CloseDB()

	// Server shutdown
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = server.Shutdown(timeoutContext)
}
