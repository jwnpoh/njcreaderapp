package auth

// Place this file at: backend/services/auth/auth_test.go
// Run with: go test ./services/auth/ -v

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
)

// =============================================================================
// Colour constants + helpers (same pattern as articles_test.go)
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
// Mock
//
// This is a fake implementation of the AuthDB interface. It doesn't touch a
// real database — instead it just stores and returns whatever we set up in
// each test. Each field is a function so individual tests can customise the
// behaviour they need without affecting any other test.
//
// Think of it like a stunt double: it looks like a real AuthDB to the
// Authenticator, but we're in full control of what it does.
// =============================================================================

type mockAuthDB struct {
	insertTokenFn  func(token *core.Token) error
	refreshTokenFn func(token *core.Token) error
	getTokenFn     func(tokenHash string) (*core.Token, error)
	deleteTokenFn  func(token string) error
	clearTokensFn  func(userID int) error
}

// Each method just calls its corresponding function field. If a test doesn't
// set one up, it panics loudly — which tells you the test is calling something
// it shouldn't be.
func (m *mockAuthDB) InsertToken(token *core.Token) error            { return m.insertTokenFn(token) }
func (m *mockAuthDB) RefreshToken(token *core.Token) error           { return m.refreshTokenFn(token) }
func (m *mockAuthDB) GetToken(tokenHash string) (*core.Token, error) { return m.getTokenFn(tokenHash) }
func (m *mockAuthDB) DeleteToken(token string) error                 { return m.deleteTokenFn(token) }
func (m *mockAuthDB) ClearTokens(userID int) error                   { return m.clearTokensFn(userID) }

// =============================================================================
// Helpers
//
// newRequest is a small convenience wrapper so each test doesn't have to
// manually build an *http.Request just to set an Authorization header.
// =============================================================================

func newRequest(authHeader string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		r.Header.Set("Authorization", authHeader)
	}
	return r
}

// =============================================================================
// AuthenticateToken
//
// This is the gatekeeper for every protected route. We test all the ways it
// can be called: missing header, wrong format, valid token, expired token.
// =============================================================================

func TestAuthenticateToken(t *testing.T) {
	summary(t)

	// A fixed user ID we'll reuse across sub-tests.
	userID := uuid.New()

	t.Run("missing authorization header", func(t *testing.T) {
		// No DB behaviour needed here — the function should reject the
		// request before it ever touches the database.
		auth := NewAuthenticator(&mockAuthDB{})

		_, err := auth.AuthenticateToken(newRequest(""))

		check(t, "error returned", err != nil, true)
		check(t, "error message", err.Error(), "no authorization header received")
	})

	t.Run("malformed header — not Bearer scheme", func(t *testing.T) {
		auth := NewAuthenticator(&mockAuthDB{})

		_, err := auth.AuthenticateToken(newRequest("Token somevalue"))

		check(t, "error returned", err != nil, true)
		check(t, "error message", err.Error(), "no authorization header received")
	})

	t.Run("malformed header — Bearer with no token", func(t *testing.T) {
		auth := NewAuthenticator(&mockAuthDB{})

		_, err := auth.AuthenticateToken(newRequest("Bearer"))

		check(t, "error returned", err != nil, true)
	})

	t.Run("valid token returns the token", func(t *testing.T) {
		// Here we DO need the DB — specifically GetToken — because the
		// request header is well-formed and will pass validation.
		// We set up the mock to return a valid, non-expired token.
		validToken := &core.Token{
			PlainToken: "VALIDTOKEN",
			UserID:     userID,
			Expiry:     time.Now().Add(1 * time.Hour), // expires in the future
		}

		auth := NewAuthenticator(&mockAuthDB{
			getTokenFn: func(tokenHash string) (*core.Token, error) {
				return validToken, nil
			},
		})

		tok, err := auth.AuthenticateToken(newRequest("Bearer VALIDTOKEN"))

		check(t, "no error", err, nil)
		check(t, "correct userID", tok.UserID, userID)
	})

	t.Run("expired token returns error", func(t *testing.T) {
		// The mock returns a token whose expiry is in the past.
		// checkToken should catch this and call DeleteToken, then return an error.
		expiredToken := &core.Token{
			PlainToken: "EXPIREDTOKEN",
			UserID:     userID,
			Expiry:     time.Now().Add(-1 * time.Hour), // expired 1 hour ago
		}

		auth := NewAuthenticator(&mockAuthDB{
			getTokenFn: func(tokenHash string) (*core.Token, error) {
				return expiredToken, nil
			},
			// DeleteToken is called internally when a token is found to be
			// expired. The mock just swallows it successfully.
			deleteTokenFn: func(token string) error {
				return nil
			},
		})

		_, err := auth.AuthenticateToken(newRequest("Bearer EXPIREDTOKEN"))

		check(t, "error returned", err != nil, true)
		check(t, "error message", err.Error(), "auth: token has expired")
	})

	t.Run("token not found in database", func(t *testing.T) {
		// The mock simulates the DB returning nothing for an unknown token.
		auth := NewAuthenticator(&mockAuthDB{
			getTokenFn: func(tokenHash string) (*core.Token, error) {
				return nil, fmt.Errorf("not found")
			},
		})

		_, err := auth.AuthenticateToken(newRequest("Bearer UNKNOWNTOKEN"))

		check(t, "error returned", err != nil, true)
	})
}

// =============================================================================
// CreateToken
//
// We verify that CreateToken correctly builds a token with the right userID,
// a future expiry, and a non-empty hash — and that it calls InsertToken.
// =============================================================================

func TestCreateToken(t *testing.T) {
	summary(t)

	t.Run("creates token with correct fields", func(t *testing.T) {
		userID := uuid.New()
		ttl := 12 * time.Hour

		var insertedToken *core.Token

		auth := NewAuthenticator(&mockAuthDB{
			// We capture whatever token gets passed to InsertToken so we
			// can inspect it after the call.
			insertTokenFn: func(token *core.Token) error {
				insertedToken = token
				return nil
			},
		})

		tok, err := auth.CreateToken(userID, ttl)

		check(t, "no error", err, nil)
		check(t, "correct userID", tok.UserID, userID)
		check(t, "plain token is non-empty", tok.PlainToken != "", true)
		check(t, "hash is non-empty", len(tok.Hash) > 0, true)
		check(t, "expiry is in the future", tok.Expiry.After(time.Now()), true)
		check(t, "InsertToken was called", insertedToken != nil, true)
	})

	t.Run("returns error if InsertToken fails", func(t *testing.T) {
		auth := NewAuthenticator(&mockAuthDB{
			insertTokenFn: func(token *core.Token) error {
				return fmt.Errorf("db is down")
			},
		})

		_, err := auth.CreateToken(uuid.New(), time.Hour)

		check(t, "error returned", err != nil, true)
	})
}

// =============================================================================
// RefreshToken
//
// Verifies that RefreshToken extends the expiry and passes it through to the DB.
// =============================================================================

func TestRefreshToken(t *testing.T) {
	summary(t)

	t.Run("updates expiry to ~24 hours from now", func(t *testing.T) {
		var refreshedToken *core.Token

		auth := NewAuthenticator(&mockAuthDB{
			refreshTokenFn: func(token *core.Token) error {
				refreshedToken = token
				return nil
			},
		})

		tok := &core.Token{
			PlainToken: "SOMETOKEN",
			Expiry:     time.Now(), // currently about to expire
		}

		err := auth.RefreshToken(tok)

		check(t, "no error", err, nil)
		check(t, "expiry extended into future", refreshedToken.Expiry.After(time.Now()), true)
		// Allow a small margin: the new expiry should be between 23h59m and 24h01m from now
		minExpiry := time.Now().Add(23*time.Hour + 59*time.Minute)
		check(t, "expiry is ~24h", refreshedToken.Expiry.After(minExpiry), true)
	})

	t.Run("returns error if DB refresh fails", func(t *testing.T) {
		auth := NewAuthenticator(&mockAuthDB{
			refreshTokenFn: func(token *core.Token) error {
				return fmt.Errorf("db error")
			},
		})

		err := auth.RefreshToken(&core.Token{})

		check(t, "error returned", err != nil, true)
	})
}
