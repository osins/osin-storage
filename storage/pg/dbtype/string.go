package dbtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DBString []string

func (j DBString) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *DBString) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		return nil

	case string:
		if err := json.Unmarshal([]byte(value), &j); err != nil {
			return err
		}

	case []byte:
		// if an empty UUID comes from a table, we return a null UUID
		if len(value) == 0 {
			return nil
		}

		if err := json.Unmarshal(value, &j); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Scan: unable to scan type %T into JSONB", value)
	}

	return nil
}
