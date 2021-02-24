package pg

import (
	"github.com/wangsying/osin-simple/simple"
	"github.com/wangsying/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

func NewUserStorage() simple.UserStorage {
	r := &userStorage{
		db: DB(),
	}

	return r
}

type userStorage struct {
	db *gorm.DB
}

func (s *userStorage) GetId(code string, password string) (string, error) {
	d := &model.User{}

	err := s.db.Where(map[string]interface{}{
		"username": code,
		"password": password,
	}).Or(map[string]interface{}{
		"mobile":   code,
		"password": password,
	}).Or(map[string]interface{}{
		"id":       code,
		"password": password,
	}).First(d).Error

	return d.Id.String(), err
}

func (s *userStorage) BindToken(token string, userId string) error {
	d := &model.AccessData{}

	return s.db.Model(d).Where("access_token", token).Update("user_id", userId).Error
}

func (s *userStorage) GetUser(id string) (interface{}, error) {
	d := &model.User{}

	err := s.db.Model(d).Where("id", id).First(d).Error
	if err != nil {
		return nil, err
	}

	return d, nil
}
