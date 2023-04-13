package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// TODO seperation of concerns candidate
// this is package global, that grows unbounded
var options map[string]string = make(map[string]string)

// the command changes but options does not
type Command struct {
	command string
	options map[string]string // this is a antipattern as this field stays global between different instances of the command object TODO refactor
}

// options stays persistent between new objects
func NewCommand(input string) *Command {
	return &Command{
		command: strings.Trim(input, "\n"),
		options: options,
	}
}

// all the parsing for slash commands in this method
func (c *Command) HandleSlashCommand() bool {

	switch {
	case strings.Contains(c.command, "/quit"), strings.Contains(c.command, "/exit"):
		fmt.Println("Thanks for Chatting")
		os.Exit(0)
	case strings.Contains(c.command, "/option"):
		words := strings.Split(c.command, " ")
		if len(words) == 3 {
			c.AddOption(words[1], words[2])
			return true
		}
	}
	return false
}

func (c *Command) AddOption(key string, value string) {
	c.options[key] = value
	fmt.Printf("%s=%s added\n", key, value)
}

func (c *Command) GetOption(key string) string {
	if value, ok := c.options[key]; ok {
		return value
	}
	return ""
}

func (c *Command) HandleChatCommand(client *openai.Client) (openai.ChatCompletionResponse, error) {

	model := c.GetOption("model")
	if len(model) == 0 {
		model = openai.GPT3Dot5Turbo
	}

	if model != openai.GPT3Dot5Turbo && model != openai.GPT4 {
		return openai.ChatCompletionResponse{}, fmt.Errorf("Invalid model we only support %s and %s but not %s", openai.GPT3Dot5Turbo, openai.GPT4, model)
	}

	return client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: c.command + "\n",
				},
			},
		},
	)
}
