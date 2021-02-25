package pg

import (
	"github.com/openshift/osin"
)

func NewPgStorage() osin.Storage {
	r := &pgStorage{
		clients:   make(map[string]osin.Client),
		authorize: make(map[string]*osin.AuthorizeData),
		access:    make(map[string]*osin.AccessData),
		refresh:   make(map[string]string),
	}

	r.clients["1234"] = &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
	}

	return r
}

type pgStorage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

func (s *pgStorage) Clone() osin.Storage {
	return s
}

func (s *pgStorage) Close() {
}

func (s *pgStorage) GetClient(id string) (osin.Client, error) {
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, osin.ErrNotFound
}

func (s *pgStorage) SetClient(id string, client osin.Client) error {
	s.clients[id] = client
	return nil
}

func (s *pgStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	s.authorize[data.Code] = data
	return nil
}

func (s *pgStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *pgStorage) RemoveAuthorize(code string) error {
	delete(s.authorize, code)
	return nil
}

func (s *pgStorage) SaveAccess(data *osin.AccessData) error {
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *pgStorage) LoadAccess(code string) (*osin.AccessData, error) {
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *pgStorage) RemoveAccess(code string) error {
	delete(s.access, code)
	return nil
}

func (s *pgStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, osin.ErrNotFound
}

func (s *pgStorage) RemoveRefresh(code string) error {
	delete(s.refresh, code)
	return nil
}
