package broker

import (
	"fmt"
	"net/http"

	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := b.Authenticator.AuthenticateToken(r)
		if err != nil {
			s := serializer.NewSerializer(true, fmt.Sprintf("%v", err), nil)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			return
		}

		// user, err := b.Users.GetUser("id", userID)
		// if err != nil {
		// 	s := serializer.NewSerializer(true, fmt.Sprintf("%v", err), nil)
		// 	s.ErrorJson(w, err)
		// 	b.Logger.Error(s, r)
		// 	return
		// }

		// s := serializer.NewSerializer(false, "user authenticated", user.Email)
		// s.Encode(w, http.StatusAccepted)
		// b.Logger.Success(s, r)
		next.ServeHTTP(w, r)
	})
}
