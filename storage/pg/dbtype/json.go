package dbtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBJson json.RawMessage

func (j DBJson) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	return json.RawMessage(j).MarshalJSON()
}

func (j *DBJson) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to unmarshal JSONB value: %s", value)
	}

	var i json.RawMessage
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*j = DBJson(i)

	return nil
}

func (j *DBJson) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// use field.Tag, field.TagSettings gets field's tags
	// checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

	fmt.Printf("\ndb type: %s, field name: %s\n", db.Dialector.Name(), field.Name)

	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
