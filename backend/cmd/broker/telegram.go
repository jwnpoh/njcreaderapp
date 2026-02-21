package broker

import (
	"net/http"

	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
	"github.com/jwnpoh/njcreaderapp/backend/services/telegram"
)

func (b *broker) SendTelegramDigest(w http.ResponseWriter, r *http.Request) {
	var articles []telegram.TelegramPayload

	s := serializer.NewSerializer(false, "successfully sent telegram digest", nil)
	err := s.Decode(w, r, &articles)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to decode input data", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	if len(articles) == 0 {
		s := serializer.NewSerializer(true, "no articles selected to send", nil)
		s.ErrorJson(w, nil)
		b.Logger.Error(s, r)
		return
	}

	err = b.Telegram.SendDigest(articles)
	if err != nil {
		s := serializer.NewSerializer(true, "failed to send telegram message", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = s.Encode(w, http.StatusAccepted)
	if err != nil {
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
	}
}
