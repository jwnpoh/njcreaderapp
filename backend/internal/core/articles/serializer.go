package articles

type ArticlesSerializer interface {
	Decode(input []byte) (*ArticleSeries, error)
	Encode(input *ArticleSeries) ([]byte, error)
}
