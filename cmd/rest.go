package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/muhfaris/swarm-autoscale/http/v1/router"
)

func HTTPServe() {
	var (
		app  = fiber.New()
		port = 2441
	)

	router.Routers(app)

	// Create a channel to receive OS signals
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	// Start the server
	go func() {
		port := fmt.Sprintf(":%d", port)
		if err := app.Listen(port); err != nil {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for OS signals
	<-osSignal
}
