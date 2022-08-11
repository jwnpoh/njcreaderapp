package broker

import (
	"fmt"
	"github.com/jwnpoh/njcreaderapp/backend/services/logger"
	"net/http"
)

// BrokerService provides an interface for cmd/main to instantiate the app.
type BrokerService interface {
	Start() error
}

type broker struct {
	Port   string
	Logger logger.Logger
}

// NewBrokerService creates a new BrokerService.
func NewBrokerService(port string) BrokerService {
	return &broker{Port: port, Logger: logger.NewAppLogger()}
}

// Start sets up a server with routes and handlers that call the various backend services.
func (b *broker) Start() error {
	// log.Printf("Starting broker service on port %s\n", b.Port)
	b.Logger.Info("Start", "broker", fmt.Sprintf("port %s", b.Port), nil)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", b.Port),
		Handler: b.SetupRouter(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		b.Logger.Error("ListenAndServe", "broker", "", nil)
	}
	return nil
}
