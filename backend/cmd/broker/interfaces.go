package broker

import (
    "net/http"
    "time"

    "github.com/google/uuid"
    "github.com/jwnpoh/njcreaderapp/backend/internal/core"
    "github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

type ArticleService interface {
    Get(page, limit int) (serializer.Serializer, error)
    GetArticle(id uuid.UUID) (serializer.Serializer, error)
    Find(q string) (serializer.Serializer, error)
    Store(input core.ArticlePayload) error
    Update(input core.ArticlePayload) error
    Delete(input []string) error
}

type AuthService interface {
    AuthenticateToken(r *http.Request) (*core.Token, error)
    RefreshToken(token *core.Token) error
    CreateToken(userID uuid.UUID, ttl time.Duration) (*core.Token, error)
    DeleteToken(token string) error
}
