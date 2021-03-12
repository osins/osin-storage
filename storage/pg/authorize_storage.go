package pg

import (
	"fmt"
	"time"

	simple_config "github.com/osins/osin-simple/simple/config"
	simple_face "github.com/osins/osin-simple/simple/model/face"
	simple_storage "github.com/osins/osin-simple/simple/storage"
	"github.com/osins/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

// NewUserStorage func define
func NewAuthorizeStorage() simple_storage.AuthorizeStorage {
	r := &authorizeStorage{
		db: DB(),
	}

	return r
}

// userStorage define

type authorizeStorage struct {
	db *gorm.DB
}

// GetId method define
func (s *authorizeStorage) Get(code string) (simple_face.Authorize, error) {
	d := &model.Authorize{}
	err := s.db.Model(d).Where("code", code).Find(d).Error
	if err != nil {
		return nil, err
	}

	if len(d.ClientId) == 0 {
		return nil, simple_config.ERROR_CLIENT_NOT_EXISTS
	}

	d.Client = &model.Client{}
	err = s.db.Model(d.Client).Where("client_id", d.ClientId).First(d.Client).Error
	if err != nil {
		return nil, simple_config.ERROR_CLIENT_NOT_EXISTS
	}

	if d.ExpireAt().Before(time.Now()) {
		return nil, fmt.Errorf("Token expired at %s.", d.ExpireAt().String())
	}

	return d, nil
}

func (s *authorizeStorage) Create(authorize simple_face.Authorize) (err error) {
	d := &model.Authorize{
		Code:                authorize.GetCode(),
		ClientId:            authorize.GetClient().GetId(),
		ExpiresIn:           authorize.GetExpiresIn(),
		Scope:               authorize.GetScope(),
		RedirectUri:         authorize.GetRedirectUri(),
		State:               authorize.GetState(),
		CreatedAt:           authorize.GetCreatedAt(),
		DeletedAt:           authorize.GetDeletedAt(),
		CodeChallenge:       authorize.GetCodeChallenge(),
		CodeChallengeMethod: authorize.GetCodeChallengeMethod(),
	}

	if authorize.GetClient().GetNeedLogin() && authorize.GetUser() != nil {
		d.UserId = authorize.GetUser().GetId()
	}

	return s.db.Model(d).Create(d).Error
}

func (s *authorizeStorage) BindUser(code string, userId string) error {
	d := &model.Authorize{}
	return s.db.Model(d).Where("code", code).Update("user_id", userId).Error
}
