package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/osins/osin-storage/storage"
	"github.com/osins/osin-storage/storage/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Debug("test init.")
}

func TestDb(t *testing.T) {
	dsn := storage.GetMySQLDSN()

	storage.Init(mysql.Open(dsn))

	storage.DB()
	storage.Migrate()

	NewClient()

	t.Logf(`Run Result:  TestDb() = %q`, dsn)
}

func NewClient() {
	u := &model.User{
		Id:       uuid.UUID(uuid.New()),
		Username: "richard",
		EMail:    "296907@qq.com",
		Password: []byte("123456"),
	}

	if err := storage.NewUserStorage().Create(u); err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error(err)
	}

	client := &model.Client{
		Id:          uuid.MustParse("12465071-a4fa-4a45-b3f1-f0a972fb6875"),
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
		NeedLogin:   true,
		NeedRefresh: true,
	}

	// pg.NewClientManage().Delete(client.Id)
	if err := storage.NewClientStorage().Create(client); err != nil {
		fmt.Println(err)
	}

	f, err := storage.NewClientStorage().Get(client.Id.String())
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nclient manage first: %v\n", f)
	}

	b, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\nclient manage first: %s\n", b)
	}
}
