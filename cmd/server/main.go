package main

import (
	"bufio"
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

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao subir o servidor")
		return
	}

	defer listener.Close()

	chatGroup := &ChatGroup{}

	fmt.Println("Servidor rodando na porta 8080")

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexao")
			return
		}

		chatGroup.Register(connection)
		go handleConnection(connection, chatGroup)
	}
}

func handleConnection(connection net.Conn, chat *ChatGroup) {
	defer func() {
		chat.Unregister(connection)
		connection.Close()
	}()

	fmt.Println("Novo cliente conectado")

	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		clientMsg := scanner.Text()

		message, err := message.FromJSONString(clientMsg)
		if err != nil {
			return
		}

		fmt.Printf("Mensagem recebida do usuário: [%s] - Mensagem [%s]\n", message.SenderName, message.MessageText)

		for _, userConn := range chat.snapshot() {
			if userConn == connection {
				continue
			}

			if err := userConn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				fmt.Printf("Erro ao definir prazo de envio ao cliente: %v\n", err)
				chat.Unregister(userConn)
				userConn.Close()
				continue
			}

			_, err = userConn.Write([]byte(message.ToJSONString()))
			if err != nil {
				fmt.Printf("Erro ao enviar mensagem ao cliente: %v\n", err)
				chat.Unregister(userConn)
				userConn.Close()
				continue
			}
		}
	}
}
