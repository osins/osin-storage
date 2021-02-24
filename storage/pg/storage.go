package pg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/openshift/osin"
	"github.com/wangsying/osin-storage/storage/pg/model"
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
	fmt.Printf("GetClient: %s\n", id)

	c := &model.Client{}
	s.db.Model(c).Where("id", id).Find(c)
	if c != nil {
		fmt.Printf("query client, find client id: %s.\n", c.Id)
		return c, nil
	}

	fmt.Printf("query client, not exists.\n")
	return nil, osin.ErrNotFound
}

func (s *storage) SetClient(id string, client osin.Client) error {
	fmt.Printf("SetClient: %s\n", id)

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
	fmt.Printf("SaveAuthorize: %s\n", data.Code)

	d := &model.AuthorizeData{}
	con := &model.AuthorizeData{Client: &model.Client{Id: data.Client.GetId()}, Code: data.Code}

	s.db.Where(con).Find(d)
	if d.Code == "" {
		d.Copy(data)

		b, err := json.Marshal(d)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(string(b))

		s.db.Model(d).Create(d)
		return nil
	}

	return fmt.Errorf("code cannot duplicate.")
}

func (s *storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	fmt.Printf("LoadAuthorize: %s\n", code)
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
	fmt.Printf("RemoveAuthorize: %s\n", code)
	d := &model.AuthorizeData{}
	s.db.Model(d).Where("code", code).Delete(d)
	return nil
}

func (s *storage) SaveAccess(data *osin.AccessData) error {
	fmt.Printf("SaveAccess: %s\n", data.AccessToken)

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
	fmt.Printf("LoadAccess: %s\n", accessToken)

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
		err = s.db.Model(u).Where("id", d.UserId).First(u).Error
		if err == nil {
			ret.UserData = u
		}
	}

	return ret, nil
}

func (s *storage) RemoveAccess(accessToken string) error {
	fmt.Printf("RemoveAccess: %s\n", accessToken)
	d := &model.AccessData{}
	s.db.Model(d).Where("access_token", accessToken).Delete(d)

	return nil
}

func (s *storage) LoadRefresh(refreshToken string) (*osin.AccessData, error) {
	fmt.Printf("LoadRefresh: %s\n", refreshToken)

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
	fmt.Printf("RemoveRefresh: %s\n", refreshToken)

	d := &model.AccessData{}
	s.db.Model(d).Where("refresh_token", refreshToken).Delete(d)

	return nil
}
