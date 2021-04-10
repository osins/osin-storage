package pg

import (
	"fmt"
	"time"

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

	if d.ExpireAt().Before(time.Now()) {
		return nil, fmt.Errorf("Token expired at %s, create at %s. expires in %v", d.ExpireAt().String(), d.CreatedAt, d.ExpiresIn)
	}

	if len(d.ClientId) == 0 {
		return nil, fmt.Errorf("get authorize data error, client id: %s", d.ClientId)
	}

	d.Client = &model.Client{}
	err = s.db.Model(d.Client).Where("id", d.ClientId).First(d.Client).Error
	if err != nil {
		return nil, fmt.Errorf("get authorize data error, client id: %s\nerror: %s", d.ClientId, err)
	}

	if len(d.UserId) > 0 {
		d.User = &model.User{}
		if err := s.db.Model(d.User).Where("id", d.UserId).First(d.User).Error; err != nil {
			return nil, fmt.Errorf("access client not exists, user id: %s", d.UserId)
		}
	} else if d.Client.NeedLogin {
		return nil, fmt.Errorf("access client need login, user not exists, client id: %s", d.ClientId)
	} else {
		d.User = nil
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
