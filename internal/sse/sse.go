package sse

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type SSE struct {
	counter     atomic.Int64
	mut         sync.Mutex
	cancel      context.CancelFunc
	connections map[int64]*Connection
	Heartbeat   time.Duration
}

func NewSSE(ctx context.Context) *SSE {
	ctx, cancel := context.WithCancel(ctx)
	heartbeat := 5 * time.Second

	gateway := &SSE{
		connections: make(map[int64]*Connection),
		cancel:      cancel,
		Heartbeat:   heartbeat,
	}

	go func() {
		time.Sleep(time.Duration(rand.Int63n(int64(heartbeat / 2))))

		ticker := time.NewTicker(heartbeat)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				gateway.Keepalive()

			case <-ctx.Done():
				return
			}
		}
	}()

	return gateway
}

func (s *SSE) Subscribe(connection *Connection) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if connection.ID == 0 {
		connection.ID = s.counter.Add(1)
	}

	s.connections[connection.ID] = connection
}

func (s *SSE) Unsubscribe(connection *Connection) {
	s.mut.Lock()
	defer s.mut.Unlock()

	delete(s.connections, connection.ID)
}

func (s *SSE) Keepalive() {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, conn := range s.connections {
		conn.Keepalive()
	}
}

func (s *SSE) Broadcast(message Message) {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, conn := range s.connections {
		conn.QueueMessage(message)
	}
}
