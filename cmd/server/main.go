package main

import (
	"bufio"
	"fmt"
	"net"
	"realtime-chat/internal/message"
)

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
	if !scanner.Scan() {
		return
	}
	join, err := message.ParseJoinRequest(scanner.Text())
	if err != nil {
		return
	}
	chat.Register(connection)

	for scanner.Scan() {
		request, err := message.ParseChatRequest(scanner.Text())
		if err != nil {
			return
		}
		outgoing := message.Message{
			SenderName:  join.SenderName,
			MessageText: request.MessageText,
		}
		fmt.Printf("Mensagem recebida do usuário: [%s] - Mensagem [%s]\n", outgoing.SenderName, outgoing.MessageText)
		chat.Broadcast(connection, outgoing)
	}
}
