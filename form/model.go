package form

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

type Form struct {
	ID               uint `gorm:"primaryKey;autoIncrement"`
	AuthorID         uint
	Title            string           `gorm:"type:varchar(100);not null"`
	Description      string           `gorm:"type:varchar(300);not null"`
	Quiz_Setting     Quiz_Setting     `gorm:"type:jsonb"`
	Response_Setting Response_Setting `gorm:"type:jsonb"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Quiz_Setting struct {
	IsQuiz        bool `gorm:"default:true"`
	DefaultPoints uint `gorm:"default:1"`
}

func (quizeS *Quiz_Setting) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Quiz_Setting{}
	err := json.Unmarshal(bytes, &result)
	*quizeS = Quiz_Setting(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (quizeS Quiz_Setting) Value() (driver.Value, error) {
	return json.Marshal(quizeS)
}

type Response_Setting struct {
	CollectEmail bool `gorm:"default:true"`
	AllowEditRes bool `gorm:"default:true"`
	SendResCopy  bool `gorm:"default:true"`
}

func (resS *Response_Setting) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := Response_Setting{}
	err := json.Unmarshal(bytes, &result)
	*resS = Response_Setting(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (resS Response_Setting) Value() (driver.Value, error) {
	return json.Marshal(resS)
}

func SetFormTable(gormDB *gorm.DB) {
	db = gormDB
	db.AutoMigrate(&Form{})
}
