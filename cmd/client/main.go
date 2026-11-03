package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"realtime-chat/cmd/client/ui"
	"realtime-chat/internal/message"

	tea "charm.land/bubbletea/v2"
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

	program := tea.NewProgram(ui.InitModel(connection, name))

	go func() {
		scanner := bufio.NewScanner(connection)
		for scanner.Scan() {
			incoming, err := message.FromJSONString(scanner.Text())
			if err != nil {
				continue
			}
			program.Send(incoming)
		}
		program.Quit()
	}()

	if _, err := program.Run(); err != nil {
		fmt.Printf("Erro na interface: %v\n", err)
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
