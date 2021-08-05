package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Room struct {
	ID 				string    	`json:"id" gorm:"primary_key:true;column:id;not null"`
	CreatedBy 		string 		`json:"created_by" gorm:"column:created_by;not null"`
	CreatedTime     int64  		`json:"created_time" gorm:"column:created_time;not null;index:idx_message_created_time"`
	Status          int    		`json:"status" gorm:"column:status;not null"`
}

func (m *Room) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	m.ID = id.String()
	m.CreatedTime = time.Now().Unix()
	return nil
}

func (m Room) TableName() string {
	return "room"
}