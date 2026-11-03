package main

import (
	"fmt"
	"net"
	"slices"
	"sync"
	"time"

	"realtime-chat/internal/message"
)

const writeTimeout = 5 * time.Second

type ChatGroup struct {
	mu          sync.RWMutex
	connections []net.Conn
}

func (g *ChatGroup) Register(conn net.Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.connections = append(g.connections, conn)
}

func (g *ChatGroup) Unregister(conn net.Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i, connection := range g.connections {
		if connection == conn {
			g.connections = slices.Delete(g.connections, i, i+1)
			return
		}
	}
}

func (g *ChatGroup) snapshot() []net.Conn {
	g.mu.RLock()
	defer g.mu.RUnlock()

	connections := make([]net.Conn, len(g.connections))
	copy(connections, g.connections)

	return connections
}

func (g *ChatGroup) Broadcast(sender net.Conn, outgoing message.Message) {
	data := []byte(outgoing.ToJSONString())
	for _, connection := range g.snapshot() {
		if connection == sender {
			continue
		}
		if err := writeMessage(connection, data); err != nil {
			fmt.Printf("Erro ao enviar mensagem ao cliente: %v\n", err)
			g.Unregister(connection)
			connection.Close()
		}
	}
}

func writeMessage(connection net.Conn, data []byte) error {
	if err := connection.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}
	_, err := connection.Write(data)
	return err
}
