package pg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/osins/osin-storage/storage/pg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormInstance *gorm.DB
var gormOnce sync.Once
var migrated bool

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
			fmt.Printf("osin-storage db dsn: %s", GetPostgresDSN())
			fmt.Printf("\nosin-storage db table prefix: %s\n", os.Getenv("PGDB_TABLE_PREFIX"))
		}

		db, err := gorm.Open(postgres.Open(GetPostgresDSN()), &gorm.Config{
			Logger:                                   newLogger,
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   os.Getenv("PGDB_TABLE_PREFIX"),    // table name prefix, table for `User` would be `t_users`
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

func GetPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("PGDB_HOST"),
		os.Getenv("PGDB_USER"),
		os.Getenv("PGDB_PASSWORD"),
		os.Getenv("PGDB_DATABASE"),
		os.Getenv("PGDB_PORT"),
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
