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
	mux.Get("/api/long", b.GetLongTopics)
	mux.Get("/api/long/{topic}", b.GetLong)

	mux.Post("/api/auth", b.Authenticate)
	mux.Post("/api/users/reset-password", b.ResetPassword)

	mux.Route("/api/users", func(mux chi.Router) {
		mux.Use(b.Auth)
		mux.Post("/logout", b.Logout)
		mux.Post("/update-user", b.UpdateUser)
		mux.Post("/get-user", b.GetUser)
		mux.Get("/friends", b.GetFriends)
		mux.Post("/follow", b.Follow)
		mux.Get("/{user}", b.ViewUser)
	})

	mux.Route("/api/posts", func(mux chi.Router) {
		mux.Use(b.Auth)
		mux.Get("/get-article", b.GetArticle)
		mux.Get("/get-post", b.GetPost)
		mux.Get("/public", b.GetPublicFeed)
		mux.Get("/following", b.GetFollowing)
		mux.Get("/notebook", b.GetNotebook)
		mux.Get("/get-liked-posts", b.GetLikedPosts)
		mux.Post("/insert", b.InsertPost)
		mux.Post("/delete", b.DeletePost)
		mux.Post("/update", b.UpdatePost)
		mux.Post("/like", b.Like)
	})

	mux.Route("/api/admin/users", func(mux chi.Router) {
		mux.Use(b.Auth)
		mux.Post("/insert-users", b.InsertUsers)
		mux.Post("/delete-user", b.DeleteUser)
	})

	mux.Route("/api/admin/articles", func(mux chi.Router) {
		mux.Use(b.Auth)
		mux.Post("/insert", b.Store)
		mux.Post("/delete", b.Delete)
		mux.Get("/update", b.Get100)
		mux.Put("/update", b.Update)
		mux.Post("/get-title", b.GetTitle)
	})

	mux.Route("/api/admin/long", func(mux chi.Router) {
		mux.Use(b.Auth)
		mux.Post("/insert", b.StoreLong)
		mux.Get("/update", b.UpdateLong)
		mux.Put("/update", b.UpdateLong)
		mux.Post("/delete", b.DeleteLong)
	})
}
