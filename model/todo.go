package model

import "time"

type Todo struct {
	ID        string     `gorm:"primarykey" json:"id" bson:"id"`
	Title     string     `json:"text" binding:"required"`
	Href      string     `json:"href,omitempty"`
	CreatedAt time.Time  `json:"-" bson:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"-" bson:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"-" bson:"deleted_at,omitempty"`
}

func (Todo) TableName() string {
	return "todos"
}
