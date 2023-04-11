package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Command struct {
	command string
}

func NewCommand(input string) *Command {
	return &Command{
		command: strings.Trim(input, "\n"),
	}
}

func (c *Command) HandleSlashCommand() {

	switch c.command {
	case "/quit", "/exit":
		fmt.Println("Thanks for Chatting")
		os.Exit(0)

	}

}

func (c *Command) HandleChatCommand(client *openai.Client) (openai.ChatCompletionResponse, error) {

	return client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: c.command + "\n",
				},
			},
		},
	)

}
