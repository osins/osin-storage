package pg

import (
	"fmt"
	"time"

	"github.com/openshift/osin"
	"github.com/osins/osin-storage/storage/pg/model"
	"gorm.io/gorm"
)

func New() osin.Storage {
	r := &storage{
		db: DB(),
	}

	return r
}

type storage struct {
	db *gorm.DB
}

func (s *storage) Clone() osin.Storage {
	return s
}

func (s *storage) Close() {
}

func (s *storage) GetClient(id string) (osin.Client, error) {

	c := &model.Client{}
	s.db.Model(c).Where("id", id).Find(c)
	if c != nil {
		return c, nil
	}

	return nil, osin.ErrNotFound
}

func (s *storage) SetClient(id string, client osin.Client) error {

	c := &model.Client{}

	s.db.Model(c).Where("id", id).Find(c)
	c.Copy(client)
	if c.Id == "" {
		s.db.Model(c).Create(c)
		return nil
	}

	s.db.Model(c).Where("id", id).Save(c)

	return nil
}

func (s *storage) SaveAuthorize(data *osin.AuthorizeData) error {

	d := &model.AuthorizeData{}
	con := &model.AuthorizeData{Client: &model.Client{Id: data.Client.GetId()}, Code: data.Code}

	s.db.Where(con).Find(d)
	if d.Code == "" {
		d.Copy(data)

		s.db.Model(d).Create(d)
		return nil
	}

	return fmt.Errorf("code cannot duplicate.")
}

func (s *storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	d := &model.AuthorizeData{}
	s.db.Model(d).Where("code", code).Find(d)
	if d != nil {
		if d.ClientId == "" {
			return nil, osin.ErrNotFound
		}

		c, err := s.GetClient(d.ClientId)
		if err != nil {
			return nil, osin.ErrNotFound
		}

		data := d.Osin()
		data.Client = c

		if data.ExpireAt().Before(time.Now()) {
			return nil, fmt.Errorf("Token expired at %s.", data.ExpireAt().String())
		}

		return data, nil
	}

	return nil, osin.ErrNotFound
}

func (s *storage) RemoveAuthorize(code string) error {
	d := &model.AuthorizeData{}
	s.db.Model(d).Where("code", code).Delete(d)
	return nil
}

func (s *storage) SaveAccess(data *osin.AccessData) error {

	d := &model.AccessData{}
	s.db.Model(d).Where("access_token", data.AccessToken).Find(d)
	if d.AccessToken == data.AccessToken {
		return fmt.Errorf("access_token cannot duplicate.")
	}

	d.Copy(data)
	s.db.Omit("Client", "AuthorizeData").Create(d)

	return nil
}

func (s *storage) LoadAccess(accessToken string) (*osin.AccessData, error) {

	d := &model.AccessData{}
	err := s.db.Model(d).Where("access_token", accessToken).Find(d).Error
	if err != nil {
		return nil, err
	}

	c := &model.Client{}
	err = s.db.Model(c).Where("id", d.ClientId).Find(c).Error
	if err != nil {
		return nil, err
	}

	ret := d.Osin()
	ret.Client = c.Osin()

	u := &model.User{}

	if len(d.UserId) > 0 {
		if err = s.db.Model(u).Where("id", d.UserId).First(u).Error; err == nil {
			ret.UserData = u
		}
	}

	return ret, nil
}

func (s *storage) RemoveAccess(accessToken string) error {
	d := &model.AccessData{}
	s.db.Model(d).Where("access_token", accessToken).Delete(d)

	return nil
}

func (s *storage) LoadRefresh(refreshToken string) (*osin.AccessData, error) {

	d := &model.AccessData{}
	s.db.Model(d).Where("refresh_token", refreshToken).Find(d)
	if d.RefreshToken == refreshToken {
		c := &model.Client{}
		s.db.Model(c).Where("id", d.ClientId).First(c)

		ret := d.Osin()
		ret.Client = c.Osin()

		return ret, nil
	}

	return nil, osin.ErrNotFound
}

func (s *storage) RemoveRefresh(refreshToken string) error {

	d := &model.AccessData{}
	s.db.Model(d).Where("refresh_token", refreshToken).Delete(d)

	return nil
}
