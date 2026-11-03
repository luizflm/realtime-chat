package message

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	MaxSenderNameLength  = 32
	MaxMessageTextLength = 2000
)

type Message struct {
	SenderName  string
	MessageText string
}

func ValidateSenderName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("nome é obrigatório")
	}
	if utf8.RuneCountInString(name) > MaxSenderNameLength {
		return fmt.Errorf("nome deve ter no máximo %d caracteres", MaxSenderNameLength)
	}
	return nil
}

func (m Message) Validate() error {
	if err := ValidateSenderName(m.SenderName); err != nil {
		return err
	}
	if strings.TrimSpace(m.MessageText) == "" {
		return fmt.Errorf("o texto da mensagem é obrigatório")
	}
	if utf8.RuneCountInString(m.MessageText) > MaxMessageTextLength {
		return fmt.Errorf("o texto da mensagem deve ter no máximo %d caracteres", MaxMessageTextLength)
	}
	return nil
}

func FromJSONString(data string) (Message, error) {
	var message Message

	err := json.Unmarshal([]byte(data), &message)
	if err != nil {
		return Message{}, fmt.Errorf("JSON inválido: %w", err)
	}
	if err := message.Validate(); err != nil {
		return Message{}, err
	}

	return message, nil
}

func (m Message) ToJSONString() string {
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s\n", string(data))
}
