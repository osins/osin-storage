package pg

import (
	"github.com/google/uuid"
	"github.com/osins/osin-simple/simple/model/face"
	simple_face "github.com/osins/osin-simple/simple/model/face"
	"github.com/osins/osin-simple/simple/storage"
	"github.com/osins/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

// NewClientStorage Client storage
func NewClientStorage() storage.ClientStorage {
	return &clientStorage{
		db: DB(),
	}
}

// clientStorage define
type clientStorage struct {
	db *gorm.DB
}

func (s *clientStorage) Create(data face.Client) error {
	d := &model.Client{
		Id:          uuid.MustParse(data.GetId()),
		Secret:      data.GetSecret(),
		RedirectUri: data.GetRedirectUri(),
		NeedLogin:   data.GetNeedLogin(),
		NeedRefresh: data.GetNeedRefresh(),
	}

	return s.db.Model(d).Where(d).Create(d).Error
}

// GetClient method define
func (s *clientStorage) Get(clientId string) (simple_face.Client, error) {

	d := &model.Client{}

	err := s.db.Model(d).Where("id", clientId).First(d).Error
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Delete method define
func (s *clientStorage) Delete(clientId string) error {

	c := &model.Client{}
	return s.db.Model(c).Where("id", clientId).Delete(c).Error
}
