package message

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func ValidateSenderName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("nome é obrigatório")
	}
	if utf8.RuneCountInString(name) > MaxSenderNameLength {
		return fmt.Errorf("nome deve ter no máximo %d caracteres", MaxSenderNameLength)
	}
	return nil
}

func ValidateMessageText(text string) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("o texto da mensagem é obrigatório")
	}
	if utf8.RuneCountInString(text) > MaxMessageTextLength {
		return fmt.Errorf("o texto da mensagem deve ter no máximo %d caracteres", MaxMessageTextLength)
	}
	return nil
}
