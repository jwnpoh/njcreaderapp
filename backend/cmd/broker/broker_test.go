package broker

// Place this file at: backend/cmd/broker/broker_test.go
// Run with: go test ./cmd/broker/ -v

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/logger"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

// =============================================================================
// Colour helpers
// =============================================================================

const (
	colGreen = "\033[32m"
	colRed   = "\033[31m"
	colBold  = "\033[1m"
	colReset = "\033[0m"
)

func check(t *testing.T, name string, got, want any) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Log(fmt.Sprintf("%s✓%s  %s", colGreen, colReset, name))
	} else {
		t.Errorf(
			"%s✗%s  %s%s%s\n     %sgot:%s  %#v\n     %swant:%s %#v",
			colRed, colReset,
			colBold, name, colReset,
			colRed, colReset, got,
			colGreen, colReset, want,
		)
	}
}

func summary(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(fmt.Sprintf("%s%s  FAIL  %s%s", colBold, colRed, t.Name(), colReset))
		} else {
			t.Log(fmt.Sprintf("%s%s  PASS  %s%s", colBold, colGreen, t.Name(), colReset))
		}
	})
}

// =============================================================================
// Mocks
// =============================================================================

type mockArticleService struct {
	getFn        func(page, limit int) (serializer.Serializer, error)
	getArticleFn func(id uuid.UUID) (serializer.Serializer, error)
	findFn       func(q string) (serializer.Serializer, error)
	storeFn      func(input core.ArticlePayload) error
	updateFn     func(input core.ArticlePayload) error
	deleteFn     func(input []string) error
}

func (m *mockArticleService) Get(page, limit int) (serializer.Serializer, error) {
	return m.getFn(page, limit)
}
func (m *mockArticleService) GetArticle(id uuid.UUID) (serializer.Serializer, error) {
	return m.getArticleFn(id)
}
func (m *mockArticleService) Find(q string) (serializer.Serializer, error) { return m.findFn(q) }
func (m *mockArticleService) Store(input core.ArticlePayload) error        { return m.storeFn(input) }
func (m *mockArticleService) Update(input core.ArticlePayload) error       { return m.updateFn(input) }
func (m *mockArticleService) Delete(input []string) error                  { return m.deleteFn(input) }

type mockAuthService struct {
	authenticateTokenFn func(r *http.Request) (*core.Token, error)
	refreshTokenFn      func(token *core.Token) error
	createTokenFn       func(userID uuid.UUID, ttl time.Duration) (*core.Token, error)
	deleteTokenFn       func(token string) error
}

func (m *mockAuthService) AuthenticateToken(r *http.Request) (*core.Token, error) {
	return m.authenticateTokenFn(r)
}
func (m *mockAuthService) RefreshToken(token *core.Token) error { return m.refreshTokenFn(token) }
func (m *mockAuthService) CreateToken(userID uuid.UUID, ttl time.Duration) (*core.Token, error) {
	return m.createTokenFn(userID, ttl)
}
func (m *mockAuthService) DeleteToken(token string) error { return m.deleteTokenFn(token) }

type mockLogger struct{}

func (l *mockLogger) Info(s serializer.Serializer, r *http.Request)    {}
func (l *mockLogger) Error(s serializer.Serializer, r *http.Request)   {}
func (l *mockLogger) Success(s serializer.Serializer, r *http.Request) {}

// =============================================================================
// Test broker factory + helpers
// =============================================================================

func newTestBroker(articles ArticleService, auth AuthService) *broker {
	return &broker{
		Articles:      articles,
		Authenticator: auth,
		Logger:        &mockLogger{},
	}
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	return bytes.NewBuffer(b)
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) jsonResponse {
	t.Helper()
	var resp jsonResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return resp
}

// chiRequest injects a chi URL param into the request context so chi.URLParam
// works outside of a real router.
func chiRequest(r *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// =============================================================================
// Auth middleware
// =============================================================================

func TestAuthMiddleware(t *testing.T) {
	summary(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("missing token blocks request", func(t *testing.T) {
		b := newTestBroker(nil, &mockAuthService{
			authenticateTokenFn: func(r *http.Request) (*core.Token, error) {
				return nil, fmt.Errorf("no authorization header received")
			},
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		b.Auth(next).ServeHTTP(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})

	t.Run("valid token passes through to next handler", func(t *testing.T) {
		b := newTestBroker(nil, &mockAuthService{
			authenticateTokenFn: func(r *http.Request) (*core.Token, error) {
				return &core.Token{PlainToken: "VALID", UserID: uuid.New(), Expiry: time.Now().Add(time.Hour)}, nil
			},
			refreshTokenFn: func(token *core.Token) error { return nil },
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer VALID")

		b.Auth(next).ServeHTTP(w, r)

		check(t, "status code", w.Code, http.StatusOK)
	})

	t.Run("refresh failure blocks request", func(t *testing.T) {
		b := newTestBroker(nil, &mockAuthService{
			authenticateTokenFn: func(r *http.Request) (*core.Token, error) {
				return &core.Token{PlainToken: "VALID", UserID: uuid.New(), Expiry: time.Now().Add(time.Hour)}, nil
			},
			refreshTokenFn: func(token *core.Token) error { return fmt.Errorf("db is down") },
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer VALID")

		b.Auth(next).ServeHTTP(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})
}

// =============================================================================
// Get handler — GET /api/articles/{page}
// =============================================================================

func TestGet(t *testing.T) {
	summary(t)

	t.Run("valid page returns 202", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			getFn: func(page, limit int) (serializer.Serializer, error) {
				return serializer.NewSerializer(false, "got articles", nil), nil
			},
		}, nil)

		w := httptest.NewRecorder()
		r := chiRequest(httptest.NewRequest(http.MethodGet, "/api/articles/1", nil), "page", "1")

		b.Get(w, r)

		check(t, "status code", w.Code, http.StatusAccepted)
		check(t, "no error", decodeResponse(t, w).Error, false)
	})

	t.Run("non-numeric page returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{}, nil)

		w := httptest.NewRecorder()
		r := chiRequest(httptest.NewRequest(http.MethodGet, "/api/articles/abc", nil), "page", "abc")

		b.Get(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})

	t.Run("service error returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			getFn: func(page, limit int) (serializer.Serializer, error) {
				return nil, fmt.Errorf("db is down")
			},
		}, nil)

		w := httptest.NewRecorder()
		r := chiRequest(httptest.NewRequest(http.MethodGet, "/api/articles/1", nil), "page", "1")

		b.Get(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})
}

// =============================================================================
// Find handler — GET /api/articles/find?term=...
// =============================================================================

func TestFind(t *testing.T) {
	summary(t)

	t.Run("valid search term returns 202", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			findFn: func(q string) (serializer.Serializer, error) {
				return serializer.NewSerializer(false, "found articles", nil), nil
			},
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/articles/find?term=geopolitics", nil)

		b.Find(w, r)

		check(t, "status code", w.Code, http.StatusAccepted)
		check(t, "no error", decodeResponse(t, w).Error, false)
	})

	t.Run("service error returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			findFn: func(q string) (serializer.Serializer, error) {
				return nil, fmt.Errorf("no results")
			},
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/articles/find?term=geopolitics", nil)

		b.Find(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})
}

// =============================================================================
// Store handler — POST /api/admin/articles/insert
// =============================================================================

func TestStore(t *testing.T) {
	summary(t)

	validPayload := core.ArticlePayload{
		{
			Title:    "Test Article",
			URL:      "https://example.com",
			Tags:     "geopolitics;2023 - Q5",
			Date:     "2023-01-01",
			MustRead: "on",
		},
	}

	t.Run("valid payload returns 202", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			storeFn: func(input core.ArticlePayload) error { return nil },
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/admin/articles/insert", jsonBody(t, validPayload))
		r.Header.Set("Content-Type", "application/json")

		b.Store(w, r)

		check(t, "status code", w.Code, http.StatusAccepted)
		check(t, "no error", decodeResponse(t, w).Error, false)
	})

	t.Run("service error returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			storeFn: func(input core.ArticlePayload) error { return fmt.Errorf("db is down") },
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/admin/articles/insert", jsonBody(t, validPayload))
		r.Header.Set("Content-Type", "application/json")

		b.Store(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})

	t.Run("malformed body returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/admin/articles/insert", bytes.NewBufferString("not json"))
		r.Header.Set("Content-Type", "application/json")

		b.Store(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})
}

// =============================================================================
// Update handler — PUT /api/admin/articles/update
// =============================================================================

func TestUpdate(t *testing.T) {
	summary(t)

	validPayload := core.ArticlePayload{
		{
			ID:       uuid.New(),
			Title:    "Updated Title",
			URL:      "https://example.com",
			Tags:     "environment",
			Date:     "2023-06-01",
			MustRead: "",
		},
	}

	t.Run("valid payload returns 202", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			updateFn: func(input core.ArticlePayload) error { return nil },
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/api/admin/articles/update", jsonBody(t, validPayload))
		r.Header.Set("Content-Type", "application/json")

		b.Update(w, r)

		check(t, "status code", w.Code, http.StatusAccepted)
		check(t, "no error", decodeResponse(t, w).Error, false)
	})

	t.Run("service error returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			updateFn: func(input core.ArticlePayload) error { return fmt.Errorf("db is down") },
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/api/admin/articles/update", jsonBody(t, validPayload))
		r.Header.Set("Content-Type", "application/json")

		b.Update(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})

	t.Run("malformed body returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/api/admin/articles/update", bytes.NewBufferString("not json"))
		r.Header.Set("Content-Type", "application/json")

		b.Update(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})
}

// =============================================================================
// Delete handler — POST /api/admin/articles/delete
// =============================================================================

func TestDelete(t *testing.T) {
	summary(t)

	t.Run("valid ids return 202", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			deleteFn: func(input []string) error { return nil },
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/admin/articles/delete", jsonBody(t, []string{"id-1", "id-2"}))
		r.Header.Set("Content-Type", "application/json")

		b.Delete(w, r)

		check(t, "status code", w.Code, http.StatusAccepted)
		check(t, "no error", decodeResponse(t, w).Error, false)
	})

	t.Run("service error returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{
			deleteFn: func(input []string) error { return fmt.Errorf("db is down") },
		}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/admin/articles/delete", jsonBody(t, []string{"id-1"}))
		r.Header.Set("Content-Type", "application/json")

		b.Delete(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})

	t.Run("malformed body returns 400", func(t *testing.T) {
		b := newTestBroker(&mockArticleService{}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/admin/articles/delete", bytes.NewBufferString("not json"))
		r.Header.Set("Content-Type", "application/json")

		b.Delete(w, r)

		check(t, "status code", w.Code, http.StatusBadRequest)
		check(t, "error flag", decodeResponse(t, w).Error, true)
	})
}

// Compile-time interface check.
var _ logger.Logger = (*mockLogger)(nil)
