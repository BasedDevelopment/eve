package util

import (
	"encoding/json"
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Validateable[T any] interface {
	Validate() error
	*T
}

type CreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Disabled bool   `json:"disabled"`
	IsAdmin  bool   `json:"is_admin"`
	Remarks  string `json:"remarks"`
}

func (s CreateRequest) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required, validation.Length(2, 20)),
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(10, 0), is.PrintableASCII),
		validation.Field(&s.Disabled, validation.Required),
		validation.Field(&s.IsAdmin, validation.Required),
		validation.Field(&s.Remarks, validation.Required),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s LoginRequest) Validate() error {
	if len(s.Email) != len("luke@ericz.me") {
		return errors.New("email is wrong length")
	}

	return nil
}

func ParseRequest[R LoginRequest | CreateRequest, T Validateable[R]](r *http.Request, rq T) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rq); err != nil {
		return err
	}

	return rq.Validate()
}
