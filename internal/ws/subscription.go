package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type subscription struct {
	h *Hub
	conn *connection
	room string
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		s.h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				ll.S.Errorw("[readPump] IsUnexpectedCloseError", "err", err)
			}
			break
		}
		ll.S.Infow("receive message", "msg", msg)
		m := &message{s.conn.ID, msg, s.room}
		if err := s.h.handler.handle(m); err != nil {
			ll.S.Errorw("[readPump] error while handling message", "err", err, "msg", string(m.data), "room", m.room)
			payload := s.h.handler.getReplyMessage(ErrorMessage, "", s.conn.ID, &MessageData{Content: err.Error()})
			if m.data, err = json.Marshal(payload); err != nil {
				ll.S.Errorw("[readPump] error while marshaling error message", "err", err)
				continue
			}
			s.conn.send <- m.data
		} else {
			s.h.broadcast <- m
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
