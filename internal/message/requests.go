package message

import (
	"encoding/json"
	"fmt"
)

type JoinRequest struct {
	SenderName string
}

type ChatRequest struct {
	MessageText string
}

func (r ChatRequest) Validate() error {
	return ValidateMessageText(r.MessageText)
}

func ParseJoinRequest(data string) (JoinRequest, error) {
	var request JoinRequest
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return JoinRequest{}, fmt.Errorf("JSON inválido: %w", err)
	}
	if err := ValidateSenderName(request.SenderName); err != nil {
		return JoinRequest{}, err
	}
	return request, nil
}

func ParseChatRequest(data string) (ChatRequest, error) {
	var request ChatRequest
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		return ChatRequest{}, fmt.Errorf("JSON inválido: %w", err)
	}
	if err := request.Validate(); err != nil {
		return ChatRequest{}, err
	}
	return request, nil
}
