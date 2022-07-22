package serializer

import (
	"encoding/json"
	"fmt"

	"github.com/jwnpoh/njcreaderapp/backend/internal/articles-service"
)

type ArticleSeries struct{}

func (aSeries *ArticleSeries) Decode(input []byte) (*articles.ArticleSeries, error) {
	series := &articles.ArticleSeries{}
	if err := json.Unmarshal(input, series); err != nil {
		return nil, fmt.Errorf("%w: error occured at serializer.ArticleSeries.Decode", err)
	}
	return series, nil
}

func (aSeries *ArticleSeries) Encode(input *articles.ArticleSeries) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("%w: error occured at serializer.ArticleSeries.Encode", err)
	}
	return rawMsg, nil
}
