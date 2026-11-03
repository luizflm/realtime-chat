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

	var user string
	for {
		fmt.Println("Escreva o seu nome para entrar no chat")
		if !input.Scan() {
			if err := input.Err(); err != nil {
				fmt.Printf("Erro ao ler nome: %v\n", err)
			}
			return
		}
		user = input.Text()
		if err := message.ValidateSenderName(user); err != nil {
			fmt.Printf("Nome inválido: %v\n", err)
			continue
		}
		break
	}

	go func() {
		serverScanner := bufio.NewScanner(serverConnection)
		for serverScanner.Scan() {
			serverMessage, err := message.FromJSONString(serverScanner.Text())
			if err != nil {
				fmt.Printf("Mensagem recebida inválida: %v\n", err)
				continue
			}
			fmt.Printf("[%s]: %s", serverMessage.SenderName, serverMessage.MessageText)
		}
	}()

	fmt.Println("Digite sua mensagem:")

	for input.Scan() {
		msg := message.Message{
			SenderName:  user,
			MessageText: input.Text(),
		}
		if err := msg.Validate(); err != nil {
			fmt.Printf("Mensagem inválida: %v\n", err)
			continue
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
