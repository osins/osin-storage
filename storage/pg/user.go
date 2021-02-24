package pg

import (
	"github.com/wangsying/osin-simple/simple"
	"github.com/wangsying/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

func NewValidateUser() simple.ValidateUser {
	r := &userStorage{
		db: DB(),
	}

	return r
}

type userStorage struct {
	db *gorm.DB
}

func (s *userStorage) Vaildate(code string, password string) error {
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

	return err
}
