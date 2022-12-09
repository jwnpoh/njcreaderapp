package broker

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
	"github.com/jwnpoh/njcreaderapp/backend/services/profanity"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) InsertUser(w http.ResponseWriter, r *http.Request) {
	var u = []core.User{
		// {
		// 	Email:       "tom@test.com",
		// 	Hash:        "testing",
		// 	Role:        "student",
		// 	DisplayName: "Tom",
		// 	Class:       "22SH01",
		// 	LastLogin:   time.Now().Format("02 Jan 2006"),
		// },
		// {
		// 	Email:       "dick@test.com",
		// 	Hash:        "testing",
		// 	Role:        "student",
		// 	DisplayName: "Dick",
		// 	Class:       "22SH02",
		// 	LastLogin:   time.Now().Format("02 Jan 2006"),
		// },
		// {
		// 	Email:       "harry@test.com",
		// 	Hash:        "testing",
		// 	Role:        "student",
		// 	DisplayName: "Harry",
		// 	Class:       "22SH03",
		// 	LastLogin:   time.Now().Format("02 Jan 2006"),
		// },
		// {
		// 	Email:       "jane@test.com",
		// 	Hash:        "testing",
		// 	Role:        "student",
		// 	DisplayName: "Jane",
		// 	Class:       "22SH04",
		// 	LastLogin:   time.Now().Format("02 Jan 2006"),
		// },
	}

	for _, v := range u {
		err := b.Users.InsertUser(&v)
		if err != nil {
			s := serializer.NewSerializer(true, "unable to add new user", err)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			fmt.Println(err)
			return
		}
	}

	s := serializer.NewSerializer(false, "successfully added new users", u)
	s.Encode(w, http.StatusAccepted)
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
	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}

func (b *broker) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		UserID      int    `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		DisplayName string `json:"display_name"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userInput)

	profanityCheck := profanity.CheckProfanity(userInput.DisplayName)
	if profanityCheck.IsProfane {
		s := serializer.NewSerializer(true, fmt.Sprintf("Please use clean language on this platform.\nThe system auto-detected the use of the profanity: '%s'.\nIf you think this is a mistake, please report the false positive via the feedback form.", profanityCheck.Profanity), userInput.UserID)
		s.Encode(w, http.StatusBadRequest)
		b.Logger.Info(s, r)
		return
	}

	user, err := b.Users.GetUser("id", userInput.UserID)
	if err != nil {
		s := serializer.NewSerializer(true, "invalid credentials", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = hasher.CheckHash(user.Hash, userInput.OldPassword)
	if err != nil {
		s := serializer.NewSerializer(true, "invalid credentials", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	if userInput.NewPassword == "" {
		userInput.NewPassword = userInput.OldPassword
	}

	newUser := &core.User{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Class:       user.Class,
		LastLogin:   time.Now().Format("02 Jan 2006"),
		DisplayName: userInput.DisplayName,
		Hash:        userInput.NewPassword,
	}

	err = b.Users.UpdateUser(newUser)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to update user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}

	s = serializer.NewSerializer(false, "successfully updated user", newUser.Email)
	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
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
	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}

func (b *broker) ViewUser(w http.ResponseWriter, r *http.Request) {
	q := chi.URLParam(r, "user")

	userID, err := strconv.Atoi(q)
	if err != nil {
		s := serializer.NewSerializer(true, "something went wrong with the page number", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	user, err := b.Users.GetUser("id", userID)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s := serializer.NewSerializer(false, "successfully retrieved user", user)
	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}
