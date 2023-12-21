package broker

import (
	"net/http"

	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) GetStats(w http.ResponseWriter, r *http.Request) {
	data, err := b.Stats.GetStats()
	if err != nil {
		s := serializer.NewSerializer(true, "unable to get articles from database", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	err = data.Encode(w, http.StatusAccepted)
	if err != nil {
		data.ErrorJson(w, err)
		b.Logger.Error(data, r)
	}
}
