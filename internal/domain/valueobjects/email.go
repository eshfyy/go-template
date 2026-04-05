package valueobjects

import (
	"fmt"
	"net/mail"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Email{}, fmt.Errorf("email: cannot be empty")
	}

	addr, err := mail.ParseAddress(value)
	if err != nil {
		return Email{}, fmt.Errorf("email: invalid format: %w", err)
	}

	return Email{value: addr.Address}, nil
}

func (e Email) String() string { return e.value }
