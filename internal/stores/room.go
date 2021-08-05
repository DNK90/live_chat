package stores

import (
	"github.com/dnk90/chat/internal/models"
	"gorm.io/gorm"
)

type RoomStore struct {
	*gorm.DB
}

func NewRoomStore(db *gorm.DB) *RoomStore {
	return &RoomStore{db}
}

func (m *RoomStore) Save(userName string) (string, error) {
	room := &models.Room{
		CreatedBy: userName,
	}
	if err := room.BeforeCreate(m.DB); err != nil {
		return "", err
	}
	return room.ID, m.Create(room).Error
}

func (m *RoomStore) UpdateStatus(messageId, status int) error {
	return m.Model(&models.Message{}).UpdateColumn("status", status).Error
}

func (m * RoomStore) GetMessages(fromId, limit int, room string) ([]*models.Message, error) {
	var messages []*models.Message
	err := m.Model(&models.Message{}).Where("id < ? AND room = ?",fromId, room).Order("id ASC").Limit(limit).Find(&messages).Error
	return messages, err
}
