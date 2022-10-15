package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

// Get makes a call to the articles service to retrieve articles from the db for a given page.
func (b *broker) Get(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	n, err := strconv.Atoi(page)
	if err != nil {
		s := serializer.NewSerializer(true, "something went wrong with the page number", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	// service, err := articles.NewArticlesService()
	// if err != nil {
	// 	s := serializer.NewSerializer(true, "unable to start articles service", err)
	// 	s.ErrorJson(w, err)
	// 	b.Logger.Error(s, r)
	// 	return
	// }

	data, err := b.Articles.Get(n)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data.Encode(w, http.StatusAccepted)
	b.Logger.Success(data, r)
}

// Store parses the new article input in the request body and sends it to the db via articles service.
func (b *broker) Store(w http.ResponseWriter, r *http.Request) {
	data := make(core.ArticleSeries, 0)

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &data)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = b.Articles.Store(data)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	b.Logger.Success(s, r)
}

// Find makes a call to the database via the articles service to search for a match of the given search term specified in the url params.
func (b *broker) Find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("term")
	// service, err := articles.NewArticlesService()
	// if err != nil {
	// 	s := serializer.NewSerializer(true, "unable to start articles service for search", err)
	// 	s.ErrorJson(w, err)
	// 	b.Logger.Error(s, r)
	// 	return
	// }

	data, err := b.Articles.Find(q)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to find results for search term", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data.Encode(w, http.StatusAccepted)
	b.Logger.Success(data, r)
}

func (b *broker) Authenticate(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"username"`
		Password string `json:"password"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userInput)

	user, err := b.Users.GetUser("email", userInput.Email)
	if err != nil {
		s := serializer.NewSerializer(true, "invalid credentials", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}

	err = hasher.CheckHash(user.Hash, userInput.Password)
	if err != nil {
		s := serializer.NewSerializer(true, "invalid credentials", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	token, err := b.Authenticator.CreateToken(user.ID, 24*time.Hour)
	if err != nil {
		s := serializer.NewSerializer(true, fmt.Sprintf("unable to generate token for user %s, %v", user.Email, err), nil)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}

	s = serializer.NewSerializer(false, fmt.Sprintf("successfully generated token for user %s", user.Email), token)
	s.Encode(w, http.StatusAccepted)
	b.Logger.Success(s, r)
}

func (b *broker) InsertUserTest(w http.ResponseWriter, r *http.Request) {
	var u = core.User{
		Email:     "jwn.poh@gmail.com",
		Hash:      "testing",
		Role:      "admin",
		LastLogin: time.Now().Format("02 Jan 2006"),
	}

	err := b.Users.InsertUser(&u)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add new user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}
	s := serializer.NewSerializer(false, "successfully added new user", u)
	b.Logger.Success(s, r)
}

func (b *broker) GetUserTest(w http.ResponseWriter, r *http.Request) {
	user, err := b.Users.GetUser("email", "jwn.poh@gmail.com")
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}

	s := serializer.NewSerializer(false, "successfully retrieved user", user)
	b.Logger.Success(s, r)
}

func (b *broker) UpdateUserTest(w http.ResponseWriter, r *http.Request) {
	newPassword := "mypassword"

	u, _ := b.Users.GetUser("email", "jwn.poh@gmail.com")

	err := b.Users.UpdateUserPassword(u.ID, newPassword)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to update user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}
	s := serializer.NewSerializer(false, "successfully updated user", u.ID)
	b.Logger.Success(s, r)
}

func (b *broker) DeleteUserTest(w http.ResponseWriter, r *http.Request) {
	user, _ := b.Users.GetUser("email", "jwn.poh@gmail.com")

	err := b.Users.DeleteUser(user.ID)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to delete user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}

	s := serializer.NewSerializer(false, "successfully deleted user", user.ID)
	b.Logger.Success(s, r)
}
