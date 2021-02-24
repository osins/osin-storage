package model

import (
	"encoding/json"
	"time"

	"github.com/openshift/osin"
	"github.com/wangsying/osin-storage/storage/pg/dbtype"
)

type AccessData struct {
	// Access token
	AccessToken string `gorm:"primaryKey;->;<-:create;"`

	ClientId string

	UserId string

	Code string

	// Refresh Token. Can be blank
	RefreshToken string

	// Token expiration in seconds
	ExpiresIn int32

	// Requested scope
	Scope string

	// Redirect Uri from request
	RedirectUri string

	// Date created
	CreatedAt time.Time

	// Data to be passed to storage. Not used by the library.
	UserData dbtype.DBJson `sql:"type:jsonb"`
}

func (s *AccessData) Copy(access *osin.AccessData) *AccessData {
	s.AccessToken = access.AccessToken
	s.ClientId = access.Client.GetId()
	s.CreatedAt = access.CreatedAt
	s.RedirectUri = access.RedirectUri
	s.RefreshToken = access.RefreshToken
	s.Scope = access.Scope
	s.ExpiresIn = access.ExpiresIn

	u, err := json.Marshal(access.UserData)
	if err != nil {
		s.UserData = u
	}

	return s
}

func (s *AccessData) Osin() *osin.AccessData {
	return &osin.AccessData{
		AccessToken:  s.AccessToken,
		CreatedAt:    s.CreatedAt,
		RedirectUri:  s.RedirectUri,
		RefreshToken: s.RefreshToken,
		ExpiresIn:    s.ExpiresIn,
		Scope:        s.Scope,
		UserData:     s.UserData,
	}
}
