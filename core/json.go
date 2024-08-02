package core

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type (
	JSONList []interface{}

	JSONMap map[string]interface{}
)

// Value Marshal
func (a JSONList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a *JSONList) LoadFromStringArray(arr []string) {
	for _, item := range arr {
		*a = append(*a, item)
	}
}

func (a *JSONList) ToStringArray() []string {
	var arr []string
	for _, item := range *a {
		if str, ok := item.(string); ok {
			arr = append(arr, str)
		}
	}
	return arr
}

func NewJSONMap() JSONMap {
	return make(JSONMap)
}

// Value Marshal
func (a JSONMap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a *JSONMap) GetString(key string) string {
	if val, ok := (*a)[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

func (a *JSONMap) GetInt64(key string) int64 {
	if val, ok := (*a)[key]; ok {
		if intVal, ok := val.(int64); ok {
			return intVal
		} else if intVal, ok := val.(float64); ok {
			return int64(intVal)
		}
	}
	return 0
}

func (a *JSONMap) GetBool(key string) bool {
	if val, ok := (*a)[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

func (a *JSONMap) GetMap(key string) JSONMap {
	if val, ok := (*a)[key]; ok {
		if mapVal, ok := val.(JSONMap); ok {
			return mapVal
		}
	}
	return NewJSONMap()
}

func (a *JSONMap) SetValue(key string, value interface{}) {
	(*a)[key] = value
}
