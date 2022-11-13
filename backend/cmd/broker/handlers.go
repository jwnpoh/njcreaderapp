package broker

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
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

	data, err := b.Articles.Get(n, 10)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data.Encode(w, http.StatusAccepted)
	// b.Logger.Success(data, r)
}

// Get100 makes a call to the articles service to retrieve 100 articles from the db for admin editing/deleting.
func (b *broker) Get100(w http.ResponseWriter, r *http.Request) {
	data, err := b.Articles.Get(1, 100)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data.Encode(w, http.StatusAccepted)
	// b.Logger.Success(data, r)
}

// Find makes a call to the database via the articles service to search for a match of the given search term specified in the url params.
func (b *broker) Find(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("term")

	data, err := b.Articles.Find(q)
	if err != nil {
		s := serializer.NewSerializer(true, fmt.Sprintf("unable to find results for search term %s", q), err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	data.Encode(w, http.StatusAccepted)
	// b.Logger.Success(data, r)
}

// Store parses the new article input in the request body and sends it to the db via articles service.
func (b *broker) Store(w http.ResponseWriter, r *http.Request) {
	input := make(core.ArticlePayload, 0)

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	for i, item := range input {
		date := formatDate(item.Date)
		input[i].Date = date
	}

	err = b.Articles.Store(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

// Update parses the updated articles and sends the update to the db via the articles service.
func (b *broker) Update(w http.ResponseWriter, r *http.Request) {
	input := make(core.ArticlePayload, 0)

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	for i, item := range input {
		date := formatDate(item.Date)
		input[i].Date = date
	}

	err = b.Articles.Update(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

// Delete takes a range of article ids and deletes theme from the database.
func (b *broker) Delete(w http.ResponseWriter, r *http.Request) {
	input := make([]string, 0)

	s := serializer.NewSerializer(false, "", nil)
	err := s.Decode(w, r, &input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = b.Articles.Delete(input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to store input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func (b *broker) Authenticate(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userInput)

	user, err := b.Users.GetUser("email", userInput.Email)
	if err != nil {
		s := serializer.NewSerializer(true, "invalid credentials", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = hasher.CheckHash(user.Hash, userInput.Password)
	if err != nil {
		s := serializer.NewSerializer(true, "invalid credentials", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	token, err := b.Authenticator.CreateToken(user.ID, 12*time.Hour)
	if err != nil {
		s := serializer.NewSerializer(true, fmt.Sprintf("unable to generate token for user %s, %v", user.Email, err), nil)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}

	s = serializer.NewSerializer(false, fmt.Sprintf("successfully generated token for user %s", user.Email), token)
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func (b *broker) InsertUser(w http.ResponseWriter, r *http.Request) {
	// var u = core.User{
	// 	Email:     "jwn.poh@gmail.com",
	// 	Hash:      "testing",
	// 	Role:      "admin",
	// 	LastLogin: time.Now().Format("02 Jan 2006"),
	// }

	// err := b.Users.InsertUser(&u)
	// if err != nil {
	// 	s := serializer.NewSerializer(true, "unable to add new user", err)
	// 	s.ErrorJson(w, err)
	// 	b.Logger.Error(s, r)
	// 	fmt.Println(err)
	// 	return
	// }

	// var u = core.User{
	// 	Email:     "joel_poh_weinan@moe.edu.sg",
	// 	Hash:      "toonpitchgoaltyneHow3",
	// 	Role:      "admin",
	// 	LastLogin: time.Now().Format("02 Jan 2006"),
	// }

	// err := b.Users.InsertUser(&u)
	// if err != nil {
	// 	s := serializer.NewSerializer(true, "unable to add new user", err)
	// 	s.ErrorJson(w, err)
	// 	b.Logger.Error(s, r)
	// 	fmt.Println(err)
	// 	return
	// }
	// s := serializer.NewSerializer(false, "successfully added new user", u)
	// b.Logger.Success(s, r)
}

func (b *broker) GetUser(w http.ResponseWriter, r *http.Request) {
	token, err := b.Authenticator.AuthenticateToken(r)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to authenticate user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	user, err := b.Users.GetUser("id", token.UserID)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s := serializer.NewSerializer(false, "successfully retrieved user", user)
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func (b *broker) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func (b *broker) DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func (b *broker) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := b.Authenticator.AuthenticateToken(r)

	if err != nil {
		s := serializer.NewSerializer(true, "unable to authenticate user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = b.Authenticator.DeleteToken(token.PlainToken)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to log user out", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s := serializer.NewSerializer(false, "successfully logged user out", nil)
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func (b *broker) GetTitle(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Url string `json:"url"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userInput)

	title := regexp.MustCompile(`\s?<meta[^p]*property=\"og:title\"\s?content=\"[^\"]*\"\s?\/?>`)

	resp, err := http.Get(userInput.Url)
	if err != nil {
		s := serializer.NewSerializer(true, "url seems to be invalid", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}
	defer resp.Body.Close()

	b2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse article content", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	t := title.Find(b2)

	regexHead := regexp.MustCompile(`\s?<meta[^p]*property=\"og:title\"\s?content=\"`)
	regexTail := regexp.MustCompile(`\"\s?\/?>`)

	head := regexHead.Find(t)
	tail := regexTail.Find(t)

	output := bytes.TrimPrefix(t, head)
	output = bytes.TrimSuffix(output, tail)
	output = bytes.TrimSpace(output)

	titleText := html.UnescapeString(string(output))
	if titleText == "" {
		s := serializer.NewSerializer(true, "could not find title, please enter manually", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s = serializer.NewSerializer(false, "successfully retrieved title", titleText)
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func formatDate(date string) string {
	day := regexp.MustCompile(`^\S+\s`)
	dateString := regexp.MustCompile(`^\S+\s\S+\s\S+`)

	res := day.ReplaceAllString(date, "")
	res = dateString.FindString(res)

	return res
}
