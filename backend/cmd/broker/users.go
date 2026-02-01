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

func (b *broker) InsertUsers(w http.ResponseWriter, r *http.Request) {
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

	input := make([]core.User, 0)

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &input)

	users := make([]core.User, 0, len(input))

	for _, v := range input {
		v.LastLogin = time.Now().Format("02 Jan 2006")

		users = append(users, v)
	}

	err = b.Users.InsertUsers(&users)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to add new users", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		s, err = b.Mailer.AdminConfirmation(user.Email, "Unable to add new users to database.")
		if err != nil {
			s := serializer.NewSerializer(true, "unable to send email confirmation to notify error in inserting new users", err)
			b.Logger.Error(s, r)
			fmt.Println(err)
		}
		return
	}

	// TODO: format table to display list of users added.

	s, err = b.Mailer.AdminConfirmation(user.Email, "Successfully inserted new users.")
	if err != nil {
		s := serializer.NewSerializer(true, "unable to send email confirmation to confirm new users insertion", err)
		b.Logger.Error(s, r)
		// We don't return here because the insertion was successful
	}

	if len(users) > 0 {
		s = serializer.NewSerializer(false, "successfully inserted new users", nil)
		s.Encode(w, http.StatusOK)
	}
}

func (b *broker) UpdateClasses(w http.ResponseWriter, r *http.Request) {
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

	input := make([]core.User, 0)

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &input)

	if len(input) > 0 {
		s = serializer.NewSerializer(false, "successfully sent data to backend. a confirmation email will be sent shortly after the update is completed.", nil)
		s.Encode(w, http.StatusAccepted)
	}

	err = b.Users.UpdateClasses(&input)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to update classes", err)
		// s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		s, err = b.Mailer.AdminConfirmation(user.Email, "Unable to update classes in database.")
		if err != nil {
			s := serializer.NewSerializer(true, "unable to send email confirmation to notify error in updating classes", err)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			fmt.Println(err)
			return
		}
		return
	}

	// TODO: format table to display list of users added.

	s, err = b.Mailer.AdminConfirmation(user.Email, "Successfully updated classes.")
	if err != nil {
		s := serializer.NewSerializer(true, "unable to send email confirmation to confirm class update", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}
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
	token, err := b.Authenticator.AuthenticateToken(r)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to authenticate user", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	var userInput struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		DisplayName string `json:"display_name"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userInput)

	profanityCheck := profanity.CheckProfanity(userInput.DisplayName)
	if profanityCheck.IsProfane {
		s := serializer.NewSerializer(true, "Please use clean language on this platform.\nThe system auto-detected the use of profanity in your entry.\nIf you think this is a mistake, please submit a request with a screenshot of your entry via your tutor.", token.UserID)
		s.Encode(w, http.StatusBadRequest)
		b.Logger.Info(s, r)
		return
	}

	user, err := b.Users.GetUser("id", token.UserID)
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
		Hash:        user.Hash,
	}

	err = b.Users.UpdateUser(newUser, userInput.NewPassword)
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

func (b *broker) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var userEmail string

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userEmail)

	user, err := b.Users.GetUser("email", userEmail)
	if err != nil {
		s := serializer.NewSerializer(true, "user not found", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	t, err := time.Parse("02 Jan 2006", user.LastLogin)
	if err != nil {
		s := serializer.NewSerializer(true, "user not found", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	if time.Since(t) < 1*time.Minute {
		// if err != nil {
		s := serializer.NewSerializer(true, "operation not allowed", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
		// }
	}

	newRandPassword := core.GenerateRandomString()
	newHash, err := hasher.GenerateHash(newRandPassword)

	newUser := &core.User{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Class:       user.Class,
		LastLogin:   time.Now().Format("02 Jan 2006"),
		DisplayName: user.DisplayName,
		Hash:        newHash,
	}

	err = b.Users.UpdateUser(newUser, newRandPassword)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to reset password. try again or contact system adminstrator", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}

	s, err = b.Mailer.ResetPassword(newUser, newRandPassword)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to reset password. try again or contact system adminstrator", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		fmt.Println(err)
		return
	}

	err = serializer.NewSerializer(false, "successfully reset password.", nil).Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}
