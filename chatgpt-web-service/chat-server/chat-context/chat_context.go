package chatcontext

import "github.com/sashabaranov/go-openai"

type Message struct {
	ID      string                       `json:"id,omitempty"`
	PID     string                       `json:"pid,omitempty"`
	Message openai.ChatCompletionMessage `json:"message,omitempty"`
	Tokens  int                          `json:"tokens,omitempty"`
}

type ContextCache interface {
	Set(string, Message, int) error
	Get(string) (*Message, error)
}
