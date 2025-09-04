package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	swarmgo "github.com/mohan2020coder/swarmgo"
	"github.com/mohan2020coder/swarmgo/llm"
)

var codeFenceRx = regexp.MustCompile("(?s)```([a-zA-Z0-9_+-]*)\\s*(.*?)```")

func FindCodeBlock(s string) (lang, code string, found bool) {
	m := codeFenceRx.FindStringSubmatch(s)
	if len(m) == 3 {
		lang = strings.ToLower(strings.TrimSpace(m[1]))
		code = m[2]
		if lang == "" {
			lang = "python"
		}
		return lang, code, true
	}
	return "", "", false
}

// AskAgent with persistent memory
func AskAgent(ctx context.Context, client *swarmgo.Swarm, agent *swarmgo.Agent, prompt string, mem MemoryProvider) (string, error) {
	logStep(agent.Name, "processing request...")

	// Save user prompt
	if err := mem.Add(agent.Name, llm.RoleUser, prompt); err != nil {
		return "", err
	}

	// Retrieve conversation
	messages, err := mem.Get(agent.Name)
	if err != nil {
		return "", err
	}

	ctx2, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	resp, err := client.Run(ctx2, agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		return "", err
	}
	if len(resp.Messages) == 0 {
		return "", fmt.Errorf("empty response")
	}

	answer := resp.Messages[len(resp.Messages)-1].Content

	// Save assistant response
	if err := mem.Add(agent.Name, llm.RoleAssistant, answer); err != nil {
		return "", err
	}

	logStep(agent.Name, "response received")
	return answer, nil
}

func logStep(agent, msg string) {
	log.Printf("[%s] %s\n", agent, msg)
}
