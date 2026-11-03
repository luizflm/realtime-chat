package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"realtime-chat/internal/message"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor")
		return
	}
	defer connection.Close()

	input := bufio.NewScanner(os.Stdin)
	name, err := readSenderName(input)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("Erro ao ler nome: %v\n", err)
		}
		return
	}

	encoder := json.NewEncoder(connection)
	if err := encoder.Encode(message.JoinRequest{SenderName: name}); err != nil {
		fmt.Printf("Erro ao entrar no chat: %v\n", err)
		return
	}

	go receiveMessages(connection)
	if err := sendMessages(input, encoder); err != nil {
		fmt.Println(err)
	}
}

func readSenderName(input *bufio.Scanner) (string, error) {
	for {
		fmt.Println("Escreva o seu nome para entrar no chat")
		if !input.Scan() {
			if err := input.Err(); err != nil {
				return "", err
			}
			return "", io.EOF
		}
		name := input.Text()
		if err := message.ValidateSenderName(name); err != nil {
			fmt.Printf("Nome inválido: %v\n", err)
			continue
		}
		return name, nil
	}
}

func receiveMessages(connection io.Reader) {
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		incoming, err := message.FromJSONString(scanner.Text())
		if err != nil {
			fmt.Printf("Mensagem recebida inválida: %v\n", err)
			continue
		}
		fmt.Printf("[%s]: %s", incoming.SenderName, incoming.MessageText)
	}
}

func sendMessages(input *bufio.Scanner, encoder *json.Encoder) error {
	fmt.Println("Digite sua mensagem:")
	for input.Scan() {
		request := message.ChatRequest{MessageText: input.Text()}
		if err := request.Validate(); err != nil {
			fmt.Printf("Mensagem inválida: %v\n", err)
			continue
		}
		if err := encoder.Encode(request); err != nil {
			return fmt.Errorf("Erro ao enviar mensagem: %w", err)
		}
	}
	if err := input.Err(); err != nil {
		return fmt.Errorf("Erro ao ler entrada: %w", err)
	}
	return nil
}
