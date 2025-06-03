package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSON)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id" gorm:"index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action    string    `json:"action" gorm:"size:50;not null"`
	TableName string    `json:"table_name" gorm:"size:50;not null"`
	RecordID  uint      `json:"record_id" gorm:"not null"`
	OldValues JSON      `json:"old_values,omitempty" gorm:"type:json"`
	NewValues JSON      `json:"new_values,omitempty" gorm:"type:json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
