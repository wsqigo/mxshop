package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	// 手动指定类型，以防外键创建失败
	ID        int32          `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, g)
}
