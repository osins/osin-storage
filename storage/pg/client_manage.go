package pg

import (
	"github.com/openshift/osin"
	"github.com/wangsying/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

func NewClientManage() ClientManage {
	return &clientMange{
		db: DB(),
	}
}

type ClientManage interface {
	Create(client *osin.DefaultClient) error
	Delete(clientId string) error
	First(clientId string) osin.Client
}

type clientMange struct {
	db *gorm.DB
}

func (s *clientMange) Create(client *osin.DefaultClient) error {
	c := &model.Client{}
	c.Copy(client)
	return s.db.Model(c).Create(c).Error
}

func (s *clientMange) Delete(clientId string) error {
	c := &model.Client{}
	return s.db.Model(c).Where("id", clientId).Delete(c).Error
}

func (s *clientMange) First(clientId string) osin.Client {
	c := &model.Client{
		Id: clientId,
	}
	s.db.Model(c).Where(c).First(c)
	return c.Osin()
}
