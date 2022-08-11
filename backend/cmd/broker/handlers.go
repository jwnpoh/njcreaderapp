package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/articles"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/users"
)

// Get makes a call to the articles service to retrieve articles from the db for a given page.
func (b *broker) Get(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	n, err := strconv.Atoi(page)
	if err != nil {
		s := serializer.NewSerializer(true, "something went wrong with the page number", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, page, s)
		return
	}

	service, err := articles.NewArticlesService()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start articles service", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, page, s)
		return
	}

	data, err := service.Get(n)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, page, err)
		return
	}

	data.Encode(w, http.StatusAccepted)
	b.Logger.Success(r.Method, r.URL.Path, page, nil)
}

// Store parses the new article input in the request body and sends it to the db via articles service.
func (b *broker) Store(w http.ResponseWriter, r *http.Request) {
	data := make(core.ArticleSeries, 0)

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &data)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "", s)
		return
	}

	service, err := articles.NewArticlesService()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start articles service to store input", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "", s)
		return
	}

	err = service.Store(data)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "", s)
		return
	}
	b.Logger.Success(r.Method, r.URL.Path, "", data)
}

// Find makes a call to the database via the articles service to search for a match of the given search term specified in the url params.
func (b *broker) Find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("term")
	service, err := articles.NewArticlesService()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start articles service for search", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "search", q)
		return
	}

	data, err := service.Find(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to find results for search term", nil)
		s.ErrorJson(w, err)
		b.Logger.Info(r.Method, r.URL.Path, q, s)
		return
	}

	data.Encode(w, http.StatusAccepted)
	b.Logger.Success(r.Method, r.URL.Path, q, nil)
}

func (b *broker) InsertUserTest(w http.ResponseWriter, r *http.Request) {
	var u = core.User{
		Email:     "jwn.poh@gmail.com",
		Hash:      "testing",
		Role:      "admin",
		LastLogin: time.Now().Format("02 Jan 2006"),
	}

	um, err := users.NewUserManager()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start user manager service", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "insertusertest", u)
		fmt.Println(err)
		return
	}

	err = um.InsertUser(&u)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add new user", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "insertusertest", u)
		fmt.Println(err)
		return
	}
}

func (b *broker) GetUserTest(w http.ResponseWriter, r *http.Request) {
	um, err := users.NewUserManager()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to start user manager service", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "getusertest", nil)
		fmt.Println(err)
		return
	}

	user, err := um.GetUser("jwn.poh@gmail.com")
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get user", nil)
		s.ErrorJson(w, err)
		b.Logger.Error(r.Method, r.URL.Path, "getusertest", nil)
		fmt.Println(err)
		return
	}

	fmt.Println(*user)

}
