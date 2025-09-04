package main

import (
	"context"
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	swarmgo "github.com/mohan2020coder/swarmgo"
	"github.com/mohan2020coder/swarmgo/llm"
)

func main() {
	dotenv.Load()

	// use OpenRouter
	client := swarmgo.NewSwarm(os.Getenv("OPENROUTER_API_KEY"), llm.OpenRouter)

	agent := &swarmgo.Agent{
		Name:         "Agent",
		Instructions: "You are a helpful agent.",
		Model:        "deepseek/deepseek-chat", // check OpenRouter model catalog for correct slug
	}

	messages := []llm.Message{
		{Role: llm.RoleUser, Content: "Hi!"},
	}

	ctx := context.Background()
	response, err := client.Run(ctx, agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)
}
