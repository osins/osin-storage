package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wangsying/osin"
	"github.com/wangsying/osin-storage/storage/pg/dbtype"
	"gorm.io/gorm"
)

type Client struct {
	Id          string `gorm:"primaryKey;->;<-:create;"`
	Secret      string
	RedirectUri string
	UserData    dbtype.DBJson
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (d *Client) GetId() string {
	return d.Id
}

func (d *Client) GetSecret() string {
	return d.Secret
}

func (d *Client) GetRedirectUri() string {
	return d.RedirectUri
}

func (d *Client) GetUserData() interface{} {
	return d.UserData
}

// Implement the ClientSecretMatcher interface
func (d *Client) ClientSecretMatches(secret string) bool {
	return d.Secret == secret
}

func (d *Client) Copy(client osin.Client) *Client {
	d.Id = client.GetId()
	d.Secret = client.GetSecret()
	d.RedirectUri = client.GetRedirectUri()

	u, err := json.Marshal(client.GetUserData())
	fmt.Printf("\nosin.Client to model.Client: \n%s\n%s\n", client.GetUserData(), u)
	if err != nil {
		fmt.Printf("\nosin.Client to model.Client error: \n%s\n%s\n%s\n", client.GetUserData(), u, err.Error())
	}

	d.UserData = u
	return d
}

func (d *Client) Osin() osin.Client {
	return d
}
