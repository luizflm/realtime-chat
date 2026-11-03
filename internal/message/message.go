package message

import (
	"encoding/json"
	"fmt"
)

const (
	MaxSenderNameLength  = 32
	MaxMessageTextLength = 2000
)

type Message struct {
	SenderName  string
	MessageText string
}

func (m Message) Validate() error {
	if err := ValidateSenderName(m.SenderName); err != nil {
		return err
	}
	return ValidateMessageText(m.MessageText)
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
