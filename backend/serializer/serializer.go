package serializer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Serializer provides an interface to read and write JSON data passed between server and client.
type Serializer interface {
	Encode(w http.ResponseWriter, status int, headers ...http.Header) error
	// Decode(w http.ResponseWriter, r *http.Request, data any) error
	ErrorJson(w http.ResponseWriter, err error, status ...int) error
}

type serializer struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewSerializer instantiates a new serializer struct to carry out encoding/decoding.
func NewSerializer(err bool, msg string, data ...any) Serializer {
	return &serializer{
		Error:   err,
		Message: msg,
		Data:    data,
	}
}

// Encode writes a JSON payload to response writer to return to the client.
func (s *serializer) Encode(w http.ResponseWriter, status int, headers ...http.Header) error {
	out, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("unable to encode json data - %w", err)
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("unable to write json response - %w", err)
	}

	return nil
}

// Decode reads a JSON payload from the request into data interface that is specified by the application service.
func Decode(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return fmt.Errorf("unable to decode json data - %w", err)
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body contained more than one json value - %w", err)
	}

	return nil
}

// ErrorJson writes an error payload back to the client.
func (s *serializer) ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	s.Error = true
	s.Message = err.Error()

	return s.Encode(w, statusCode)
}
