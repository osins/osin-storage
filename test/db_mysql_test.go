package test

import (
	"testing"

	"github.com/osins/osin-storage/storage"
	"gorm.io/driver/mysql"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestDb(t *testing.T) {
	dsn := storage.GetMySQLDSN()

	storage.Init(mysql.Open(dsn))

	storage.DB()
	storage.Migrate()

	t.Logf(`Run Result:  TestDb() = %q`, dsn)
}
