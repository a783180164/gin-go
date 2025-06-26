package collection

import (
	"time"
)

type Collection struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Prompt      string    `gorm:"column:prompt;type:varchar(500)" json:"prompt"`
	Description string    `gorm:"column:description;type:varchar(500)" json:"description"`
	UUID        string    `gorm:"column:uuid;type:varchar(255);uniqueIndex" json:"uuid"`
	Status      int32     `gorm:"column:status;type:tinyint(1);not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdtime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedtime"`
	IsDeleted   bool      `gorm:"column:is_deleted;type:tinyint(1);not null;default:0" json:"is_deleted"`
}
