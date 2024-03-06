package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/osins/osin-storage/storage/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormDialector *gorm.Dialector
var gormInstance *gorm.DB
var gormOnce sync.Once
var migrated bool

func Init(dialector gorm.Dialector) {
	gormDialector = &dialector
}

// DB func define
func DB() *gorm.DB {
	gormOnce.Do(func() {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // Disable color
			},
		)

		if os.Getenv("APP_DEBUG") == "true" {
			logrus.WithFields(logrus.Fields{
				"debug":           true,
				"db table prefix": os.Getenv("DB_TABLE_PREFIX"),
				"dsn":             *gormDialector,
			}).Debug("osin-storage db dsn")
		}

		db, err := gorm.Open(*gormDialector, &gorm.Config{
			Logger:                                   newLogger,
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   os.Getenv("DB_TABLE_PREFIX"),      // table name prefix, table for `User` would be `t_users`
				SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
				NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
			},
		})

		if err != nil {
			newLogger.Error(context.Background(), err.Error())
		}

		gormInstance = db

		if !migrated {
			Migrate()
			migrated = false
		}
	})

	return gormInstance
}

func GetDSN() *gorm.Dialector {
	return gormDialector
}

func GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
}

func GetPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)
}

func Migrate() {
	var tables []interface{}
	tables = append(tables,
		&model.Client{},
		&model.Authorize{},
		&model.Access{},
		&model.User{},
	)

	for _, t := range tables {
		err := gormInstance.AutoMigrate(t)
		if err != nil {
			return
		}
	}
}
