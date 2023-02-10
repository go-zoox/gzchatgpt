package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

func Run(apiKey string) error {
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	scanner := bufio.NewScanner(os.Stdin)
	quit := false

	for !quit {
		fmt.Printf("[问题] ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		typ, payload := parse(input)
		switch typ {
		case "quit", "exit":
			quit = true
		case "question":
			answer(ctx, client, payload)
		default:
			return fmt.Errorf("invalid type: %v", typ)
		}
	}

	return nil
}

func parse(input string) (typ string, payload string) {
	if input == "exit" {
		return "exit", ""
	}

	return "question", input
}

func answer(ctx context.Context, client gpt3.Client, quesiton string) {
	resp, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}

	fmt.Printf("[回答] %s\n", resp.Choices[0].Text)
}
