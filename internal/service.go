package internal

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"os"
	"sync"
)

type TokenService struct {
	mu      sync.Mutex
	storage *Storage
	logs    *os.File
}

func NewService(l string) (*TokenService, error) {
	store := NewStorage()
	f, err := os.OpenFile(l, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open logfile: %w", err)
	}
	out := &TokenService{
		storage: store,
		logs:    f,
	}
	return out, nil
}

// WriteLog is a service layer helper function that writes messages to the logfile.
func (t *TokenService) WriteLog(msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fmt.Fprintf(t.logs, "Found value: %v\n", msg)
	return
}

type TokenSecret struct {
	Token  string `json:"token"`
	Secret string `json:"service"`
}

type WriteTokenRequest struct {
	tokenSecret TokenSecret
}

type WriteTokenResponse struct {
	Token string
}

func (a WriteTokenRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.tokenSecret, validation.Required),
	)
}

type ReadTokenRequest struct {
	Tokens []string
}

type ReadTokenResponse struct {
	tokenSecrets []TokenSecret
}

func (g ReadTokenRequest) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Tokens, validation.Required),
	)
}

type UpdateTokenRequest struct {
	tokenSecret TokenSecret
}

func (a UpdateTokenRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.tokenSecret, validation.Required),
	)
}

type DeleteTokenRequest struct {
	Token string
}

func (r DeleteTokenRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
	)
}

// Service layer functions validates the request data and
// calls the appropriate storage layer functions.
func (t *TokenService) WriteToken(req *WriteTokenRequest) (*WriteTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	encrypt, err := Encrypter([]byte(req.tokenSecret.Token), []byte(req.tokenSecret.Secret))
	if err != nil {
		return nil, err
	}

	err = t.storage.WriteToken(req.tokenSecret.Token, encrypt)
	if err != nil {
		return nil, err
	}

	out := &WriteTokenResponse{
		Token: req.tokenSecret.Token,
	}

	return out, nil
}

func (t *TokenService) ReadToken(req *ReadTokenRequest) (*ReadTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	ret, err := t.storage.ReadToken(req.Tokens)
	if err != nil {
		return nil, err
	}

	d, err := Decrypter([]byte(req.Token), string(ret))
	if err != nil {
		return nil, err
	}

	out := &ReadTokenResponse{
		Tokens: req.Tokens,
		Secret: string(d),
	}

	return out, nil
}

func (t *TokenService) UpdateToken(req *UpdateTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	encrypt, err := Encrypter([]byte(req.tokenSecret.Token), []byte(req.tokenSecret.Secret))
	if err != nil {
		return err
	}

	return t.storage.UpdateToken(req.tokenSecret.Token, encrypt)
}

func (t *TokenService) DeleteToken(req *DeleteTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	return t.storage.DeleteToken(req.Token)
}
