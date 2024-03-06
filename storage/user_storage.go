package storage

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	simple_face "github.com/osins/osin-simple/simple/model/face"
	simple_storage "github.com/osins/osin-simple/simple/storage"
	"github.com/osins/osin-storage/storage/model"
	"github.com/sirupsen/logrus"
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
	var id uuid.UUID
	id, err = uuid.Parse(data.GetId())
	if err != nil {
		return err
	}

	d := &model.User{
		Id:       id,
		Username: data.GetUsername(),
		Password: data.GetPassword(),
		Salt:     data.GetSalt(),
		EMail:    data.GetEmail(),
		Mobile:   data.GetMobile(),
	}

	if f, e := s.ExistsByCode(data.GetId(), data.GetUsername(), data.GetMobile(), data.GetEmail()); f {
		logrus.WithFields(logrus.Fields{"result": f, "error": e}).Error("user is exists.")
		return fmt.Errorf("user is exists.")
	}

	return s.db.Model(d).Create(d).Error
}

// GetId method define
func (s *userStorage) GetId(code string) (string, error) {
	if u, err := s.GetByCode(code); err != nil && u != nil {
		return "", err
	} else {
		return u.GetId(), err
	}
}

func existsByCode(db *gorm.DB, id string, username string, mobile string, email string) (*gorm.DB, error) {
	if len(id) == 0 {
		return nil, errors.New("id is null")
	}

	query := db.Model(&model.User{}).Where(map[string]interface{}{
		"id": id,
	})

	if username != "" {
		query.Or(map[string]interface{}{
			"username": username,
		})
	}

	if mobile != "" {
		query.Or(map[string]interface{}{
			"mobile": mobile,
		})
	}

	if email != "" {
		query.Or(map[string]interface{}{
			"email": email,
		})
	}

	return query, nil
}

// GetByPassword method define
func (s *userStorage) ExistsByCode(id string, username string, mobile string, email string) (bool, error) {

	count := int64(0)
	zero := int64(0)

	query, err := existsByCode(s.db, id, username, mobile, email)
	if err != nil {
		return true, err
	}

	if err := query.Count(&count).Error; err != nil || count > 0 {
		sql := s.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			q, err := existsByCode(tx, id, username, mobile, email)
			if err != nil {
				return tx
			}

			return q.Count(&count)
		})

		logrus.WithFields(logrus.Fields{
			"sql":      sql,
			"id":       id,
			"username": username,
			"mobile":   mobile,
			"email":    email,
		}).Debug("user is exists, sql out")

		return true, err
	}

	return count > zero, nil
}

func (s *userStorage) GetByCode(code string) (simple_face.User, error) {
	fmt.Printf("\nstorage user, get by code: %s\n", code)

	u := &model.User{}

	if err := s.db.Model(&model.User{}).Debug().Where(map[string]interface{}{
		"mobile": code,
	}).Or(map[string]interface{}{
		"username": code,
	}).Or(map[string]interface{}{
		"email": code,
	}).First(u).Error; err != nil {
		return nil, err
	}

	return u, nil
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
