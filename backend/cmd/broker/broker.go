package broker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jwnpoh/njcreaderapp/backend/cmd/config"
	"github.com/jwnpoh/njcreaderapp/backend/external/pscale"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
	"github.com/jwnpoh/njcreaderapp/backend/services/auth"
	"github.com/jwnpoh/njcreaderapp/backend/services/logger"
	"github.com/jwnpoh/njcreaderapp/backend/services/posts"
	"github.com/jwnpoh/njcreaderapp/backend/services/socials"
	"github.com/jwnpoh/njcreaderapp/backend/services/users"
)

// BrokerService provides an interface for cmd/main to instantiate the app.
type BrokerService interface {
	Start() error
}

type broker struct {
	Port          string
	Logger        logger.Logger
	Articles      *articles.Articles
	Authenticator *auth.Authenticator
	Users         *users.UserManager
	Posts         *posts.Posts
	Socials       *socials.Socials
}

// NewBrokerService creates a new BrokerService.
func NewBrokerService(config config.Config) BrokerService {
	articlesDB, err := pscale.NewArticlesDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	authDB, err := pscale.NewAuthDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	usersDB, err := pscale.NewUsersDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	postsDB, err := pscale.NewPostsDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	socialsDB, err := pscale.NewSocialsDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	broker := broker{
		Port:          config.Port,
		Logger:        logger.NewAppLogger(),
		Articles:      articles.NewArticlesService(articlesDB),
		Authenticator: auth.NewAuthenticator(authDB),
		Users:         users.NewUserManager(usersDB),
		Posts:         posts.NewPostsDB(postsDB, articlesDB, socialsDB, usersDB),
		Socials:       socials.NewSocialsDB(socialsDB, usersDB),
	}

	return &broker
}

// Start sets up a server with routes and handlers that call the various backend services.
func (b *broker) Start() error {
	log.Printf("Starting broker service on port %s\n", b.Port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", b.Port),
		Handler: b.SetupRouter(),
	}

	log.Fatal(srv.ListenAndServe())
	return nil
}
