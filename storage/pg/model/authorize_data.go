package model

import (
	"encoding/json"
	"time"

	"github.com/wangsying/osin"
	"github.com/wangsying/osin-storage/storage/pg/dbtype"
)

type AuthorizeData struct {
	// Authorization code
	Code string `gorm:"primaryKey;->;<-:create;"`

	ClientId string

	// Client information
	Client *Client `gorm:"foreignKey:ClientId;references:Id;"`

	// Token expiration in seconds
	ExpiresIn int32

	// Requested scope
	Scope string

	// Redirect Uri from request
	RedirectUri string

	// State data from request
	State string

	// Date created
	CreatedAt time.Time

	// Data to be passed to storage. Not used by the library.
	UserData dbtype.DBJson `sql:"type:jsonb"`

	// Optional code_challenge as described in rfc7636
	CodeChallenge string

	// Optional code_challenge_method as described in rfc7636
	CodeChallengeMethod string
}

func (d *AuthorizeData) Copy(authorize *osin.AuthorizeData) *AuthorizeData {
	c := &Client{}
	d.Client = c.Copy(authorize.Client)
	d.Code = authorize.Code
	d.ExpiresIn = authorize.ExpiresIn
	d.Scope = authorize.Scope
	d.RedirectUri = authorize.RedirectUri
	d.State = authorize.State
	d.CreatedAt = authorize.CreatedAt
	d.CodeChallenge = authorize.CodeChallenge
	d.CodeChallengeMethod = authorize.CodeChallengeMethod

	u, err := json.Marshal(authorize.UserData)
	if err != nil {
		d.UserData = u
	}

	return d
}

func (d *AuthorizeData) Osin() *osin.AuthorizeData {
	return &osin.AuthorizeData{
		Client:              d.Client.Osin(),
		Code:                d.Code,
		ExpiresIn:           d.ExpiresIn,
		Scope:               d.Scope,
		RedirectUri:         d.RedirectUri,
		State:               d.State,
		CreatedAt:           d.CreatedAt,
		UserData:            d.UserData,
		CodeChallenge:       d.CodeChallenge,
		CodeChallengeMethod: d.CodeChallengeMethod,
	}
}
