package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxBodyBytes = 1 << 20 // 1MB

// HTTPError carrega status + mensagem para responder padronizado em JSON.
type HTTPError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func BadRequest(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusBadRequest, Message: msg}
}
func UnsupportedMediaType() *HTTPError {
	return &HTTPError{Status: http.StatusUnsupportedMediaType, Message: "Content-Type deve ser application/json"}
}
func RequestTooLarge() *HTTPError {
	return &HTTPError{Status: http.StatusRequestEntityTooLarge, Message: "corpo da requisição muito grande"}
}

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, dst T) error {
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.HasPrefix(ct, "application/json") {
		return UnsupportedMediaType()
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syn *json.SyntaxError
		var ute *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syn):
			return BadRequest(fmt.Sprintf("JSON inválido (pos %d)", syn.Offset))
		case errors.Is(err, io.ErrUnexpectedEOF):
			return BadRequest("JSON inválido (truncado)")
		case errors.As(err, &ute):
			return BadRequest(fmt.Sprintf("tipo inválido para o campo %q (pos %d)", ute.Field, ute.Offset))
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return BadRequest(fmt.Sprintf("campo desconhecido %s", field))
		case errors.Is(err, io.EOF):
			return BadRequest("corpo vazio")
		case errors.Is(err, http.ErrBodyReadAfterClose):
			return BadRequest("não foi possível ler o corpo")
		default:
			return BadRequest("JSON inválido")
		}
	}

	if dec.More() {
		return BadRequest("apenas um objeto JSON é permitido")
	}

	return nil
}

func WriteJSON[T any](w http.ResponseWriter, status int, payload T) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, err error) {
	var he *HTTPError
	if errors.As(err, &he) && he != nil {
		WriteJSON(w, he.Status, he)
		return
	}
	WriteJSON(w, http.StatusInternalServerError, &HTTPError{
		Status:  http.StatusInternalServerError,
		Message: "erro interno",
	})
}
