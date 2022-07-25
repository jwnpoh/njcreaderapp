package articles

import (
	"fmt"

	fire "github.com/jwnpoh/njcreaderapp/backend/firestore"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/serializer"
)

type articleService struct {
}

// NewArticlesService returns an articleService object to implement methods to interact with Firestore repo.
func NewArticlesService() *articleService {
	return &articleService{}
}

// Get gets up to 144 documents from Firestore and serves them in pages of 12 articles each.
func (a *articleService) Get(n ...int) (serializer.Serializer, error) {
	// check if page argument is passed
	if len(n) > 1 {
		return serializer.NewSerializer(true, fmt.Sprintf("faulty request url"), nil), fmt.Errorf("too many params in request url")
	}

	articles, err := fire.NewFireStoreRepo().Get()
	if err != nil {
		return serializer.NewSerializer(true, fmt.Sprintf("unable to get articles"), nil), err
	}

	data := make([]core.Article, 0, 12)

	// if no page argument passed
	if len(n) == 0 {
		for i, j := range *articles {
			if i > 11 {
				break
			}
			data = append(data, j)
		}
		return serializer.NewSerializer(false, "hit the broker", data), nil

	}

	// if page is 2 or more
	if n[0] > 1 {
		// implement paging logic
		startAt := (n[0] - 1) * 12
		endAt := (n[0] * 12) - 1

		for i, j := range *articles {
			if i < startAt {
				continue
			}
			if i > endAt {
				break
			}
			data = append(data, j)
		}
	}
	return serializer.NewSerializer(false, "hit the broker", data), nil

	// if page is 1 or 0
}
