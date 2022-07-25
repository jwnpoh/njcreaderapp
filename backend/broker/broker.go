package broker

import (
	"fmt"
	"log"
	"net/http"
)

type BrokerService interface {
	Start() error
}

type broker struct {
	Port string
}

func NewBrokerService(port string) BrokerService {
	return &broker{Port: port}
}

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
