package broker

import (
	"fmt"
	"log"
	"net/http"
)

// BrokerService provides an interface for cmd/main to instantiate the app.
type BrokerService interface {
	Start() error
}

type broker struct {
	Port string
}

// NewBrokerService creates a new BrokerService.
func NewBrokerService(port string) BrokerService {
	return &broker{Port: port}
}

// Start sets up a server with routes and handlers that call the various backend services.
func (b *broker) Start() error {
	log.Printf("Starting broker service on port %s\n", b.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", b.Port),
		Handler: b.SetupRouter(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("unable to start broker service - %w", err)
	}

	return nil
}
