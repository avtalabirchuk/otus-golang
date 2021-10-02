package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Counter struct {
		Counter int `validate:"max:100"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{"1.0.0"},
			nil,
		},
		{
			Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			nil,
		},
		{
			User{
				ID:     "d5a668db-1da8-43b5-83d6-ecd7ee3c2cf8",
				Name:   "Andrey",
				Age:    36,
				Email:  "ya@yandex.ru",
				Role:   "stuff",
				Phones: []string{"89222223329"},
				meta:   nil,
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case Number %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}
