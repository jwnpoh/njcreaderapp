package broker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jwnpoh/njcreaderapp/backend/services/hasher"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

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

	user.LastLogin = time.Now().Format("02 Jan 2006")
	user.Hash = userInput.Password

	b.Users.UpdateUser(user)

	s = serializer.NewSerializer(false, fmt.Sprintf("successfully generated token for user %s", user.Email), token)
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
