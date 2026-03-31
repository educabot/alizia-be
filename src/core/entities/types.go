package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON is a raw JSON type for GORM JSONB columns.
type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSON", value)
	}
	*j = make(JSON, len(bytes))
	copy(*j, bytes)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	*j = make(JSON, len(data))
	copy(*j, data)
	return nil
}
