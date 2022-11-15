package broker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

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
	// s := serializer.NewSerializer(false, "successfully added new user", u)
	// s.Encode(w, http.StatusAccepted)
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
	var userInput struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		DisplayName string `json:"display_name"`
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
