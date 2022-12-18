package service

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mrasnake/TokenService/internal/storage"
	"golang.org/x/mod/sumdb/storage"
	"os"
	"strings"
	"sync"
)

type TokenService struct {
	mu      sync.Mutex
	storage *storage.Storage
	logs    *os.File
}

func NewService(l string) (*TokenService, error) {
	store := datastore.NewStorage()
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

// ProcessMassage works much like a router, reading in the request type
// and routing it to the appropriate service layer function.
func (t *TokenService) ProcessMessage(in []byte) {

	s := strings.Fields(string(in))
	switch s[0] {
	case "ADD_ITEM":
		if err := t.WriteToken(&WriteTokenRequest{
			Value: s[1],
		}); err != nil {
			fmt.Printf("failed to process message %v: %v\n", s[0], err.Error())
		}
	case "GET_ITEM":
		val, err := t.ReadToken(&ReadTokenRequest{
			Value: s[1],
		})
		if err != nil {
			fmt.Printf("failed to process message %v: %v\n", s[0], err.Error())
			return
		}
		t.WriteLog(val)
	case "REMOVE_ITEM":
		if err := t.DeleteToken(&DeleteTokenRequest{
			Value: s[1],
		}); err != nil {
			fmt.Printf("failed to process message %v: %v\n", s[0], err.Error())
		}
	case "GET_ALL_ITEMS":
		vals, err := t.UpdateTokens()
		if err != nil {
			fmt.Printf("failed to process message %v: %v\n", s[0], err.Error())
			return
		}
		for _, v := range vals {
			t.WriteLog(v)
		}
	default:
		fmt.Println("Invalid request")

	}

	return
}

type WriteTokenRequest struct {
	Token  string
	Secret string
}

func (a WriteTokenRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Token, validation.Required),
		validation.Field(&a.Secret, validation.Required),
	)
}

type ReadTokenRequest struct {
	Token string
}

func (g ReadTokenRequest) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Token, validation.Required),
	)
}

type UpdateTokenRequest struct {
	Token  string
	Secret string
}

func (a UpdateTokenRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Token, validation.Required),
		validation.Field(&a.Secret, validation.Required),
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
func (t *TokenService) WriteToken(req *WriteTokenRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	encrypt, err := Encrypter([]byte(req.Token), []byte(req.Secret))
	if err != nil {
		return nil, err
	}

	ret, err := t.storage.WriteToken(req.Token, encrypt)
}

func (t *TokenService) ReadToken(req *ReadTokenRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	return t.storage.ReadToken(req.Token)
}

func (t *TokenService) UpdateTokens() ([]string, error) {
	return t.storage.UpdateTokens()
}

func (t *TokenService) DeleteToken(req *DeleteTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	return t.storage.DeleteToken(req.Token)
}
