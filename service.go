package oniontree

import (
	"github.com/onionltd/go-oniontree/validator"
	"github.com/onionltd/go-oniontree/validator/jsonschema"
	"regexp"
	"strings"
)

type ID string

func (i ID) Validate() error {
	pattern := `^[a-z0-9\-]+$`
	matched, err := regexp.MatchString(pattern, string(i))
	if err != nil {
		return err
	}
	if !matched {
		return &ErrInvalidID{string(i), pattern}
	}
	return nil
}

type Service struct {
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	URLs        []string     `json:"urls" yaml:"urls"`
	PublicKeys  []*PublicKey `json:"public_keys,omitempty" yaml:"public_keys,omitempty"`

	id        ID
	validator *validator.Validator
}

func (s *Service) ID() string {
	return string(s.id)
}

func (s *Service) SetURLs(urls []string) int {
	s.URLs = []string{}
	return s.AddURLs(urls)
}

func (s *Service) AddURLs(urls []string) int {
	added := 0
	for _, url := range urls {
		url = strings.TrimSpace(url)
		_, exists := s.urlExists(url)
		if exists {
			continue
		}
		s.URLs = append(s.URLs, url)
		added++
	}
	return added
}

func (s *Service) SetPublicKeys(publicKeys []*PublicKey) int {
	s.PublicKeys = []*PublicKey{}
	return s.AddPublicKeys(publicKeys)
}

func (s *Service) AddPublicKeys(publicKeys []*PublicKey) int {
	added := 0
	for _, publicKey := range publicKeys {
		idx, exists := s.publicKeyExists(publicKey)
		if exists {
			s.PublicKeys[idx] = publicKey
			continue
		}
		s.PublicKeys = append(s.PublicKeys, publicKey)
		added++
	}
	return added
}

func (s *Service) Validate() error {
	if err := s.id.Validate(); err != nil {
		return err
	}
	if s.validator != nil {
		return s.validator.Validate(s)
	}
	return nil
}

func (s Service) urlExists(url string) (int, bool) {
	for idx, _ := range s.URLs {
		if s.URLs[idx] == url {
			return idx, true
		}
	}
	return -1, false
}

func (s Service) publicKeyExists(publicKey *PublicKey) (int, bool) {
	for idx, _ := range s.PublicKeys {
		if s.PublicKeys[idx].Fingerprint == "" && s.PublicKeys[idx].ID == "" {
			continue
		}
		if s.PublicKeys[idx].Fingerprint == publicKey.Fingerprint || s.PublicKeys[idx].ID == publicKey.ID {
			return idx, true
		}
	}
	return -1, false
}

func NewService(id string) *Service {
	return &Service{
		id:        ID(id),
		validator: validator.NewValidator(jsonschema.V0),
	}
}
