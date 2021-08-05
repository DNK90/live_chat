package services

import (
	"github.com/dnk90/chat/internal/config"
	"github.com/dnk90/chat/internal/models"
	"github.com/dnk90/chat/internal/stores"
)

type Service struct {
	cfg *config.Config
	MessageStore *stores.MessageStore
	RoomStore *stores.RoomStore
}

func NewService() *Service {
	cfg := config.Load()
	return &Service{
		cfg: cfg,
		MessageStore: stores.NewMessageStore(cfg.DB),
		RoomStore: stores.NewRoomStore(cfg.DB),
	}
}

func (s *Service) GetMessages(fromId, limit int, room string) ([]*models.Message, error) {
	return s.MessageStore.GetMessages(fromId, limit, room)
}

func (s *Service) AddNewMessage(room, userName, content string) (int, error) {
	return s.MessageStore.Save(room, userName, content)
}

func (s *Service) NewRoom(userName string) (string, error) {
	return s.RoomStore.Save(userName)
}

func (s *Service) RemoveMessage(messageId int, userName string) error {
	// TODO: check if userName is owner or not
	return s.MessageStore.UpdateStatus(messageId, models.REMOVED)
}
