package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	swarmgo "github.com/prathyushnallamothu/swarmgo"
	"github.com/prathyushnallamothu/swarmgo/llm"
)

const (
	OllamaBaseURL = "http://localhost:11434"
	ModelName     = "codellama"
	MaxTurns      = 6
	MemoryPath    = "memory" // folder for persisting agent memory
)

var memory MemoryProvider

func main() {
	// Ensure Docker is installed
	if _, err := exec.LookPath("docker"); err != nil {
		log.Fatal("❌ Docker not found. Please install and start Docker.")
	}

	fmt.Println("🚀 SDLC Crew is starting with model:", ModelName)

	client := swarmgo.NewSwarm(OllamaBaseURL, llm.Ollama)
	ctx := context.Background()

	// Init file-based memory
	memory = NewFileMemoryProvider(MemoryPath)

	// Step 1: Planning
	plan, err := AskWithMemory(ctx, client, Planner, "Build a CLI Fibonacci generator in Python.")
	if err != nil {
		log.Fatal(err)
	}
	WriteDocs("DOCS.md", "Planner", plan)

	// Step 2: Architecture
	arch, err := AskWithMemory(ctx, client, Architect, plan)
	if err != nil {
		log.Fatal(err)
	}
	AppendDocs("DOCS.md", "Architect", arch)

	// Step 3: Coding
	codeOutput, err := AskWithMemory(ctx, client, Coder, plan+"\n"+arch)
	if err != nil {
		log.Fatal(err)
	}
	lang, code, found := FindCodeBlock(codeOutput)
	if found {
		WriteCodeFile(lang, code)
		WriteDocs("CODE.md", "Coder", codeOutput)
	}

	// Step 4: Review
	review, err := AskWithMemory(ctx, client, Reviewer, codeOutput)
	if err != nil {
		log.Fatal(err)
	}
	WriteDocs("CODE.md", "Reviewer", review)

	// Step 5: Execute code in Docker sandbox
	if found {
		result, err := RunInDocker(ctx, lang, code)
		if err != nil {
			log.Fatal(err)
		}
		WriteTestReport("TEST_REPORT.md", result)
	}

	fmt.Println("✅ SDLC process completed. Generated: DOCS.md, CODE.md, TEST_REPORT.md, and source file.")
}

// AskWithMemory uses persistent memory for each agent
func AskWithMemory(ctx context.Context, client *swarmgo.Swarm, agent *swarmgo.Agent, prompt string) (string, error) {
	logStep(agent.Name, "loading memory...")

	// Load previous conversation
	history, err := memory.Get(agent.Name)
	if err != nil {
		return "", err
	}

	// Append new user prompt
	history = append(history, llm.Message{
		Role:    llm.RoleUser,
		Content: prompt,
	})

	logStep(agent.Name, "processing request...")

	// Run with full history + timeout
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	resp, err := client.Run(ctx2, agent, history, nil, "", false, false, MaxTurns, true)
	if err != nil {
		return "", err
	}

	if len(resp.Messages) == 0 {
		return "", fmt.Errorf("empty response")
	}

	// Take last message as answer
	answer := resp.Messages[len(resp.Messages)-1].Content
	logStep(agent.Name, "response received")

	// Persist both prompt + response
	if err := memory.Add(agent.Name, llm.RoleUser, prompt); err != nil {
		return "", err
	}
	if err := memory.Add(agent.Name, llm.RoleAssistant, answer); err != nil {
		return "", err
	}

	return answer, nil
}
