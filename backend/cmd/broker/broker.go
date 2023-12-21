package broker

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jwnpoh/njcreaderapp/backend/cmd/config"
	"github.com/jwnpoh/njcreaderapp/backend/external/pscale"
	"github.com/jwnpoh/njcreaderapp/backend/external/sheets"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
	"github.com/jwnpoh/njcreaderapp/backend/services/auth"
	"github.com/jwnpoh/njcreaderapp/backend/services/logger"
	"github.com/jwnpoh/njcreaderapp/backend/services/long"
	"github.com/jwnpoh/njcreaderapp/backend/services/mail"
	"github.com/jwnpoh/njcreaderapp/backend/services/posts"
	"github.com/jwnpoh/njcreaderapp/backend/services/socials"
	"github.com/jwnpoh/njcreaderapp/backend/services/stats"
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
	Longs         *long.Longs
	Authenticator *auth.Authenticator
	Users         *users.UserManager
	Posts         *posts.Posts
	Socials       *socials.Socials
	Stats         *stats.Stats
	Mailer        *mail.MailService
}

// NewBrokerService creates a new BrokerService.
func NewBrokerService(config config.Config) BrokerService {
	articlesDB, err := pscale.NewArticlesDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	sheetsDB, err := sheets.NewSheetsService(context.Background(), config.SheetsConfig)
	if err != nil {
		log.Fatal(err)
	}

	longsDB, err := pscale.NewLongsDB(config.DSN)
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

	statsDB, err := pscale.NewStatsDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	mailService := mail.NewMailService(config.MailServiceConfig)

	broker := broker{
		Port:          config.Port,
		Logger:        logger.NewAppLogger(),
		Articles:      articles.NewArticlesService(articlesDB, sheetsDB),
		Longs:         long.NewLongService(longsDB),
		Authenticator: auth.NewAuthenticator(authDB),
		Users:         users.NewUserManager(usersDB),
		Posts:         posts.NewPostsDB(postsDB, articlesDB, socialsDB, usersDB),
		Socials:       socials.NewSocialsDB(socialsDB, usersDB),
		Stats:         stats.NewStatsService(statsDB),
		Mailer:        mailService,
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
