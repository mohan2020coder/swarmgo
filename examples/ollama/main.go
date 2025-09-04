package main

import (
	"context"
	"fmt"
	"log"

	swarmgo "github.com/mohan2020coder/swarmgo"
	"github.com/mohan2020coder/swarmgo/llm"
)

func main() {
	// Create Swarm client with Ollama backend (provide base URL instead of API key)
	client := swarmgo.NewSwarm("http://localhost:11434", llm.Ollama)

	agent := &swarmgo.Agent{
		Name:         "Agent",
		Instructions: "You are a helpful agent.",
		// Set to the model you pulled into Ollama (e.g. llama2, mistral, codellama)
		Model: "codellama",
	}

	messages := []llm.Message{
		{Role: llm.RoleUser, Content: "Hi!"},
	}

	ctx := context.Background()
	response, err := client.Run(ctx, agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)
}
