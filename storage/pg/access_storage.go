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
		UserId:       data.GetUser().GetId(),
		ExpiresIn:    data.GetExpiresIn(),
		Scope:        data.GetScope(),
		CreatedAt:    data.GetCreatedAt(),
		DeletedAt:    data.GetDeletedAt(),
	}

	return s.db.Model(d).Create(d).Error
}

func (s *accessStorage) BindUser(code string, userId string) error {
	d := &model.Access{}
	return s.db.Model(d).Where("access_token", code).Update("user_id", userId).Error
}

// GetId method define
func (s *accessStorage) Get(code string) (simple_face.Access, error) {
	d := &model.Access{}
	err := s.db.Model(d).Where("access_token", code).Find(d).Error
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

func (s *accessStorage) GetByRefreshToken(code string) (simple_face.Access, error) {
	d := &model.Access{}
	if err := s.db.Model(d).Where("refresh_token", code).First(d).Error; err != nil {
		return nil, err
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
