package sse

type Connection struct {
	ID        int64
	Messages  chan Message
	Heartbeat chan bool
}

func NewConnection() *Connection {
	return &Connection{
		Messages:  make(chan Message),
		Heartbeat: make(chan bool),
	}
}

func (c *Connection) QueueMessage(message Message) {
	c.Messages <- message
}

func (c *Connection) Keepalive() {
	c.Heartbeat <- true
}
