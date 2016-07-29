package sockets

import "gopkg.in/mgo.v2/bson"

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan *sendRequest

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

type sendRequest struct {
	userID bson.ObjectId
	msg    []byte
}

var h = hub{
	broadcast:   make(chan *sendRequest),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				if m.userID != c.userID {
					continue
				}

				select {
				case c.send <- m.msg:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

func StartWS() {
	h.run()
}
