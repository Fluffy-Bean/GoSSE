package sse

type Connection struct {
	ID        int64
	Messages  chan string
	Heartbeat chan bool
}

func NewConnection() *Connection {
	return &Connection{
		Messages:  make(chan string),
		Heartbeat: make(chan bool),
	}
}

func (c *Connection) QueueMessage(message string) {
	c.Messages <- message
}

func (c *Connection) Keepalive() {
	c.Heartbeat <- true
}
