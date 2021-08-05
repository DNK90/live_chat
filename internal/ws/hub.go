package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {

	ID uint32

	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type message struct {
	connectionId uint32
	data []byte
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	mtx *sync.Mutex

	counter uint32

	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan *message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription

	handler *handler
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		rooms:      make(map[string]map[*connection]bool),
		handler:    newHandler(),
		mtx:        &sync.Mutex{},
	}
}

func (h *Hub) getConnectionId() uint32 {
	connectionId := atomic.LoadUint32(&h.counter)
	atomic.AddUint32(&h.counter, 1)
	return connectionId
}

func (h *Hub) NewConnection(room string, conn *connection) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*connection]bool)
	}
	h.rooms[room][conn] = true
}

func (h *Hub) RemoveConnection(room string, conn *connection) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if h.rooms[room] != nil && h.rooms[room][conn] {
		delete(h.rooms[room], conn)
		close(conn.send)
	}
	if len(h.rooms[room]) == 0 {
		delete(h.rooms, room)
	}
}

func (h *Hub) Run() {
	defer func() {
		if err := recover(); err != nil {
			ll.S.Error(err)
		}
	}()
	for {
		select {
		case s := <-h.register:
			h.NewConnection(s.room, s.conn)
			// broadcast back the connectionId
			payload, err := json.Marshal(h.handler.getReplyMessage(NewConnection, "", s.conn.ID, nil))
			if err != nil {
				ll.S.Errorw("[register][callback]error occurs", "err", err)
			} else {
				ll.S.Infow("[register][callback]register successfully", "connectionId", s.conn.ID)
				s.conn.send <- payload
			}
		case s := <-h.unregister:
			h.RemoveConnection(s.room, s.conn)
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				if c.ID != m.connectionId {
					select {
					case c.send <- m.data:
					default:
						h.RemoveConnection(m.room, c)
					}
				}
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request, roomId string) {
	ll.S.Infow("[ServeWs]roomId", "roomId", roomId)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{
		ID:   h.getConnectionId(),
		ws:   ws,
		send: make(chan []byte, 256),
	}
	s := subscription{h, c, roomId}
	h.register <- s
	go s.writePump()
	go s.readPump()
}
