package internal

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
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
	Secret string `json:"secret"`
}

type WriteTokenRequest struct {
	Secret string `json:"secret"`
}

type WriteTokenResponse struct {
	Token string `json:"token"`
}

func (a WriteTokenRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Secret, validation.Required),
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
func (t *TokenService) ReadTokens(req *ReadTokenRequest) (*ReadTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	var ts []TokenSecret

	for _, tok := range req.Tokens {
		ret, err := t.storage.ReadToken(tok)
		if err != nil {
			return nil, err
		}

		d, err := Decrypter([]byte(tok), string(ret))
		if err != nil {
			return nil, err
		}

		ts = append(ts, TokenSecret{
			Token:  tok,
			Secret: string(d),
		})
	}
	out := &ReadTokenResponse{
		tokenSecrets: ts,
	}

	return out, nil
}

func (t *TokenService) WriteToken(req *WriteTokenRequest) (*WriteTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	token := uuid.NewString()

	encrypt, err := Encrypter([]byte(token), []byte(req.Secret))
	if err != nil {
		return nil, err
	}

	err = t.storage.WriteToken(token, encrypt)
	if err != nil {
		return nil, err
	}

	out := &WriteTokenResponse{
		Token: token,
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
