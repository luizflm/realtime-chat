package message

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	SenderName  string
	MessageText string
}

func FromJSONString(data string) (Message, error) {
	var message Message

	err := json.Unmarshal([]byte(data), &message)
	if err != nil {
		return Message{}, fmt.Errorf("invalid data")
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
