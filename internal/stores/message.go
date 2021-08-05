package stores

import (
	"github.com/dnk90/chat/internal/models"
	"gorm.io/gorm"
	"time"
)

type MessageStore struct {
	*gorm.DB
}

func NewMessageStore(db *gorm.DB) *MessageStore {
	return &MessageStore{db}
}

func (m *MessageStore) Save(room, userName, content string) (int, error) {
	message := &models.Message{
		Content:     content,
		UserName:    userName,
		Room:        room,
		CreatedTime: time.Now().Unix(),
		Status:      0,
	}
	err := m.Create(message).Error
	return message.ID, err
}

func (m *MessageStore) UpdateStatus(messageId, status int) error {
	return m.Model(&models.Message{}).Where("ID=?", messageId).UpdateColumn("status", status).Error
}

func (m * MessageStore) GetMessages(fromId, limit int, room string) ([]*models.Message, error) {
	var (
		err error
		messages []*models.Message
	)
	tx := m.Model(&models.Message{}).Where("room = ?", room)
	if fromId > 1 {
		tx = tx.Where("id < ?", fromId)
	}
	err = tx.Order("id DESC").Limit(limit).Find(&messages).Error
	return messages, err
}
