package models

import "gorm.io/gorm"

const (
	SENT = 0
	REMOVED = 1
)

type Message struct {
	ID 				int    `json:"id" gorm:"primary_key:true;column:id;auto_increment;not null"`
	Content			string `json:"content" gorm:"column:content;not null"`
	UserName 		string `json:"user_name" gorm:"column:user_name;not null"`
	Room            string `json:"room" gorm:"column:room;not null;index:idx_message_room;"`
	CreatedTime     int64  `json:"created_time" gorm:"column:created_time;not null;index:idx_message_created_time"`
	Status          int    `json:"status" gorm:"column:status;not null"`
}

func (m Message) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (m Message) TableName() string {
	return "message"
}