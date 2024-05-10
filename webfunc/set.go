package webfunc

import "github.com/gorilla/websocket"

type ConnSet map[*websocket.Conn]bool

// Add adds a new element to the Set. Returns a pointer to the Set.
func (s *ConnSet) Add(t *websocket.Conn) *ConnSet {
	_, ok := (*s)[t]
	if !ok {
		(*s)[t] = true
	}
	return s
}

// Clear removes all elements from the Set
func (s *ConnSet) Clear() {
	s = &ConnSet{}
}

// Delete removes the Item from the Set and returns Has(Item)
func (s *ConnSet) Delete(conn *websocket.Conn) bool {
	_, ok := (*s)[conn]
	if ok {
		delete(*s, conn)
	}
	return ok
}

// Has returns true if the Set contains the Item
func (s *ConnSet) Has(item *websocket.Conn) bool {
	_, ok := (*s)[item]
	return ok
}

// Items returns the Item(s) stored
func (s *ConnSet) Connections() []*websocket.Conn {
	connections := []*websocket.Conn{}
	for i := range *s {
		connections = append(connections, i)
	}
	return connections
}

// Size returns the size of the set
func (s *ConnSet) Size() int {
	return len(*s)
}
