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

	mux.Post("/api/auth", b.Authenticate)

	//test user functions
	mux.Route("/api/users", func(mux chi.Router) {
		// mux.Use(b.Auth)
		mux.Get("/insert-test", b.InsertUserTest)
		mux.Get("/get-test", b.GetUserTest)
		mux.Get("/update-test", b.UpdateUserTest)
		mux.Get("/delete-test", b.DeleteUserTest)

		mux.Post("/api/articles", b.Store)
	})
}
