package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dathan/go-openai-prompt-git-save/pkg/githubstorage"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with OpenAI GPT-3",
	Run: func(cmd *cobra.Command, args []string) {
		// inialize the client
		client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		fmt.Println("Welcome to the OpenAI GPT-3 chatbot. Type 'exit' to quit.")
		for {
			fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input", err)
				os.Exit(1)
			}

			if strings.Trim(input, "\n") == "/exit" {
				fmt.Println("Thanks for chatting!")
				os.Exit(0)
			}

			resp, err := chatCommand(client, input)
			if err != nil {
				fmt.Println("Error with GPT4")
				os.Exit(1)
			}

			fmt.Printf("Response from GPT-4\n %s\n", resp.Choices[0].Message.Content)
			err = githubstorage.SaveInput(input)
			if err != nil {
				fmt.Println("Error saving input to GitHub")
				os.Exit(1)
			}
		}

	},
}

func chatCommand(client *openai.Client, content string) (resp openai.ChatCompletionResponse, err error) {

	return client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content + "\n",
				},
			},
		},
	)

}
func main() {

	err := chatCmd.Execute()
	if err != nil {
		fmt.Println("Error from chatCmd", err)
		os.Exit(2)
	}
}