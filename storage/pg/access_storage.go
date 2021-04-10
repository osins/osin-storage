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
func NewAccessStorage() simple_storage.AccessStorage {
	r := &accessStorage{
		db: DB(),
	}

	return r
}

// userStorage define

type accessStorage struct {
	db *gorm.DB
}

func (s *accessStorage) Create(data simple_face.Access) (err error) {
	d := &model.Access{
		AccessToken:  data.GetAccessToken(),
		RefreshToken: data.GetRefreshToken(),
		ClientId:     data.GetClient().GetId(),
		ExpiresIn:    data.GetExpiresIn(),
		Scope:        data.GetScope(),
		CreatedAt:    data.GetCreatedAt(),
		DeletedAt:    data.GetDeletedAt(),
	}

	if data.GetUser() != nil && len(data.GetUser().GetId()) > 0 {
		d.UserId = data.GetUser().GetId()
	}

	return s.db.Model(d).Create(d).Error
}

func (s *accessStorage) BindUser(code string, userId string) error {
	d := &model.Access{}

	if len(code) == 0 || len(userId) == 0 {
		return fmt.Errorf("code or user id is null.")
	}

	return s.db.Model(d).Where("access_token", code).Update("user_id", userId).Error
}

// GetId method define
func (s *accessStorage) Get(code string) (simple_face.Access, error) {
	d := &model.Access{}
	err := s.db.Model(d).Where("access_token", code).Find(d).Error
	if err != nil {
		return nil, err
	}

	if d.ExpireAt().Before(time.Now()) {
		return nil, fmt.Errorf("Token expired at %s, create at %s. expires in %v", d.ExpireAt().String(), d.CreatedAt, d.ExpiresIn)
	}

	if len(d.ClientId) == 0 {
		return nil, fmt.Errorf("client id not exists in access data.")
	}

	d.Client = &model.Client{}
	err = s.db.Model(d.Client).Where("id", d.ClientId).First(d.Client).Error
	if err != nil {
		return nil, simple_config.ERROR_CLIENT_NOT_EXISTS
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

func (s *accessStorage) GetByRefreshToken(code string) (simple_face.Access, error) {
	d := &model.Access{}
	if err := s.db.Model(d).Where("refresh_token", code).First(d).Error; err != nil {
		return nil, err
	}

	if len(d.ClientId) == 0 {
		return nil, fmt.Errorf("access client id is null, client id: %s", d.ClientId)
	}

	d.Client = &model.Client{}
	if err := s.db.Model(d.Client).Where("id", d.ClientId).First(d.Client).Error; err != nil {
		return nil, fmt.Errorf("access client not exists, client id: %s", d.ClientId)
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

func (s *accessStorage) RemoveAuthorize(code string) error {
	d := &model.Access{}
	return s.db.Model(d).Where("access_token", code).Delete(d).Error
}

func (s *accessStorage) RemoveRefresh(code string) error {
	d := &model.Access{}
	return s.db.Model(d).Where("refresh_token", code).Delete(d).Error
}
