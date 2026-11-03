package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"realtime-chat/internal/message"
)

func main() {
	serverConnection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor")
		return
	}

	defer serverConnection.Close()

	input := bufio.NewScanner(os.Stdin)

	fmt.Println("Escreva o seu nome para entrar no chat")
	if !input.Scan() {
		if err := input.Err(); err != nil {
			fmt.Printf("Erro ao ler nome: %v\n", err)
		}
		return
	}
	user := input.Text()

	go func() {
		serverScanner := bufio.NewScanner(serverConnection)
		for serverScanner.Scan() {
			serverMessage, _ := message.FromJSONString(serverScanner.Text())
			fmt.Printf("[%s]: %s", serverMessage.SenderName, serverMessage.MessageText)
		}
	}()

	fmt.Println("Digite sua mensagem:")

	for input.Scan() {
		msg := message.Message{
			SenderName:  user,
			MessageText: input.Text(),
		}

		if _, err := serverConnection.Write([]byte(msg.ToJSONString())); err != nil {
			fmt.Printf("Erro ao enviar mensagem: %v\n", err)
			return
		}
	}

	if err := input.Err(); err != nil {
		fmt.Printf("Erro ao ler entrada: %v\n", err)
	}
}
