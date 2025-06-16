package user

import (
	"time"
)

type User struct {
	ID         int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	User       string    `gorm:"type:varchar(255);column:user" json:"user"`
	Password   string    `gorm:"type:varchar(100)" json:"password"`
	UUID       string    `gorm:"type:char(36)" json:"uuid"`
	Nickname   string    `gorm:"type:varchar(50)" json:"nickname"`
	Username   string    `gorm:"type:varchar(50)" json:"username"`
	Avatar     string    `gorm:"type:varchar(500)" json:"avatar"`
	CreateTime time.Time `gorm:"column:creattime;autoCreateTime" json:"creattime"`
	UpdateTime time.Time `gorm:"column:updatetime;autoUpdateTime" json:"updatetime"`
	Isdelete   int32     `gorm:"type:int json:"isdelete"`
}

// func (User) TableName() string {
// 	return "user"
// }
