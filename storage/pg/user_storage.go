package pg

import (
	"github.com/google/uuid"
	simple_face "github.com/osins/osin-simple/simple/model/face"
	simple_storage "github.com/osins/osin-simple/simple/storage"
	"github.com/osins/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

// NewUserStorage func define
func NewUserStorage() simple_storage.UserStorage {
	r := &userStorage{
		db: DB(),
	}

	return r
}

// userStorage define

type userStorage struct {
	db *gorm.DB
}

// Create mothod define
func (s *userStorage) Create(data simple_face.User) (err error) {
	d := &model.User{
		Id:       uuid.MustParse(data.GetId()),
		Username: data.GetUsername(),
		Password: data.GetPassword(),
		EMail:    data.GetEmail(),
		Mobile:   data.GetMobile(),
	}

	return s.db.Model(d).Create(d).Error
}

// GetId method define
func (s *userStorage) GetId(code string, password string) (string, error) {

	if u, err := s.GetByPassword(code, password); err != nil {
		return "", err
	} else {
		return u.GetId(), err
	}
}

// GetByPassword method define
func (s *userStorage) GetByPassword(code string, password string) (simple_face.User, error) {

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

	return d, err
}

// GetUser method define
func (s *userStorage) GetById(id string) (simple_face.User, error) {

	d := &model.User{}

	err := s.db.Model(d).Where("id", id).First(d).Error
	if err != nil {
		return nil, err
	}

	return d, nil
}

// BindToken method define
func (s *userStorage) BindToken(token string, userId string) error {

	d := &model.Access{}

	return s.db.Model(d).Where("access_token", token).Update("user_id", userId).Error
}

// BindToken method define
func (s *userStorage) BindCode(code string, userId string) error {

	d := &model.Authorize{}

	return s.db.Model(d).Where("code", code).Update("user_id", userId).Error
}
