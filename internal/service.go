package internal

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TokenService struct {
	mu      sync.Mutex
	storage *Storage
	nonce   []byte
}

// NewService creates and returns a Service
func NewService() (*TokenService, error) {
	store := NewStorage()

	n := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, n); err != nil {
		return nil, err
	}

	out := &TokenService{
		storage: store,
		nonce:   n,
	}
	return out, nil
}

// Service layer request/response objects and corresponding validation functions.

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

// ReadTokens servers as the service layer GET function of /tokens endpoint,
// validating the request data, calling storage layer and decrypting the results.
func (t *TokenService) ReadTokens(req *ReadTokenRequest) (*ReadTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	var ts []TokenSecret

	// loop through tokens reading, decrypting and creating a list of TokenSecrets.
	for _, tok := range req.Tokens {
		ret, err := t.storage.ReadToken(tok)
		if err != nil {
			return nil, err
		}

		key, err := extractKey(tok)
		if err != nil {
			return nil, err
		}

		d, err := Decrypter(key, t.nonce, string(ret))
		if err != nil {
			return nil, err
		}

		ts = append(ts, TokenSecret{
			Token:  tok,
			Secret: string(d),
		})
	}

	// instantiate response object
	out := &ReadTokenResponse{
		tokenSecrets: ts,
	}

	return out, nil
}

// WriteToken servers as the service layer POST function of /tokens endpoint,
// validating, encrypting and formatting the request data, then calling the Storage layer.
func (t *TokenService) WriteToken(req *WriteTokenRequest) (*WriteTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	key := keyGen()

	encrypt, err := Encrypter(key, t.nonce, []byte(req.Secret))
	if err != nil {
		return nil, err
	}

	token := formatToken(string(key))

	err = t.storage.WriteToken(token, encrypt)
	if err != nil {
		return nil, err
	}

	out := &WriteTokenResponse{
		Token: token,
	}

	return out, nil
}

// UpdateToken servers as the service layer PUT function of /tokens endpoint,
// validating and encrypting the request data, extracts the key, then calling the Storage layer.
func (t *TokenService) UpdateToken(req *UpdateTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	key, err := extractKey(req.tokenSecret.Token)
	if err != nil {
		return err
	}

	encrypt, err := Encrypter(key, t.nonce, []byte(req.tokenSecret.Secret))
	if err != nil {
		return err
	}

	return t.storage.UpdateToken(req.tokenSecret.Token, encrypt)
}

// DeleteToken servers as the service layer DELETE function of /tokens endpoint,
// validating the request data, then calling the Storage layer.
func (t *TokenService) DeleteToken(req *DeleteTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	return t.storage.DeleteToken(req.Token)
}

func formatToken(key string) string {
	ret := "dp.token." + key
	return ret
}

func extractKey(token string) ([]byte, error) {
	formatting := token[:9]

	if formatting != "dp.token." {
		return nil, errors.New("improperly formatted token")
	}

	return []byte(token[9:]), nil
}