package broker

import (
	"fmt"
	"net/http"

	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok, err := b.Authenticator.AuthenticateToken(r)
		if err != nil {
			s := serializer.NewSerializer(true, fmt.Sprintf("%v", err), nil)
			s.ErrorJson(w, err)
			b.Logger.Error(s, r)
			return
		}

		err = b.Authenticator.RefreshToken(tok)

		next.ServeHTTP(w, r)
	})
}
