package broker

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRouter returns a chi.Router to start with the broker server.
func (b *broker) SetupRouter() chi.Router {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	b.serveRoutes(mux)

	return mux
}

func (b *broker) serveRoutes(mux chi.Router) {
	mux.Get("/api/articles/{page}", b.Get)
	mux.Get("/api/articles/find", b.Find)

	//test user functions
	mux.Get("/api/users/insert-test", b.InsertUserTest)
	mux.Get("/api/users/get-test", b.GetUserTest)

	mux.Post("/api/articles", b.Store)
}
