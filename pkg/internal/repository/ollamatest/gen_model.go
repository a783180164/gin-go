package ollamatest

import (
	"time"
)

type Collection struct {
	ID         int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(255);column:name" json:"name"`
	Prompt     string    `gorm:"type:varchar(500)" json:"prompt"`
	Desc       string    `gorm:"type:varchart(500)" json:"Desc"`
	CreateTime time.Time `gorm:"column:creattime;autoCreateTime" json:"creattime"`
	UpdateTime time.Time `gorm:"column:updatetime;autoUpdateTime" json:"updatetime"`
	IsDeleteD  int32     `gorm:"type:int" json:"isdeleted"`
}

type CollectionPoint struct {
	CreatedAt  string `json:"created_at"`
	Filename   string `json:"filename"`
	ChunkIndex int64  `json:"chunk_index"`
	ChunkSize  int64  `json:"chunk_size"`
	Content    string `json:"content"`
	Type       int    `json:"type"`
}
