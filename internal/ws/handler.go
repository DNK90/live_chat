package ws

import (
	"encoding/json"
	"errors"
	"github.com/dnk90/chat/internal/log"
	"github.com/dnk90/chat/internal/services"
	"github.com/mitchellh/mapstructure"
)

var ll = log.New()

type handler struct {
	service *services.Service
}

func newHandler() *handler {
	return &handler{
		service:   services.NewService(),
	}
}

func (h *handler) handle(msg *message) error {
	var (
		err error
		payload Payload
	)

	if err = json.Unmarshal(msg.data, &payload); err != nil {
		return err
	}
	// TODO: validate if room is exists or not
	switch payload.MessageType {
	case SendMessage:
		err = h.sendMessage(msg.room, &payload)
	case ReceiveMessage:
		err = h.receiveMessage(&payload)
	case DeleteMessage:
		err = h.deleteMessage(&payload)
	//case NewConnection:
	//	payload = h.reply(payload.UserName, )
	default:
		err = errors.New("invalid message")
	}
	if err == nil {
		msg.data, err = json.Marshal(payload)
	}
	return err
}

func (h *handler) sendMessage(room string, payload *Payload) error {
	// add message to database
	if payload.Data != nil {
		id, err := h.service.AddNewMessage(room, payload.Data.Owner, payload.Data.Content)
		if err != nil {
			return err
		}
		payload.Data.MessageId = id
		// change type to receiveMessage
		payload.MessageType = ReceiveMessage
		return nil
	}
	return errors.New("message is empty")
}

func (h *handler) receiveMessage(data *Payload) error {
	// do nothing for now, just forward
	return nil
}

func (h *handler) deleteMessage(payload *Payload) error {
	// update message to delete
	return h.service.RemoveMessage(payload.Data.MessageId, payload.Data.Owner)
}

func (h *handler) getReplyMessage(messageType int, userName string, connectionId uint32, data *MessageData) *Payload {
	return &Payload{
		UserName:     userName,
		ConnectionID: connectionId,
		MessageType:  messageType,
		Data:         data,
	}
}

func decodeMapstructure(input map[string]interface{}, object interface{}) error {
	mConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result: object,
	}
	decoder, err := mapstructure.NewDecoder(mConfig)
	if err != nil {
		ll.S.Errorw("error while creating NewDecoder", "err", err.Error())
		return err
	}
	if err = decoder.Decode(input); err != nil {
		ll.S.Errorw("error while decoding input", "err", err.Error())
		return err
	}
	return nil
}
