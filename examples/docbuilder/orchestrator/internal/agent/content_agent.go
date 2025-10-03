package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/example/document-builder/internal/cache"
	"github.com/example/document-builder/internal/logging"
)

type contentAgent struct {
	cache     *cache.SQLiteCache
	logger    *logging.Logger
	ollamaURL string
}

func NewContentAgent(c *cache.SQLiteCache, logger *logging.Logger) *contentAgent {
	return &contentAgent{cache: c, logger: logger, ollamaURL: "http://localhost:11434/completions"}
}

func (a *contentAgent) GenerateOutline(ctx context.Context, topic, audience, level string) (Outline, error) {
	prompt := fmt.Sprintf("Create a detailed outline for a document about: %s. Audience: %s. Level: %s", topic, audience, level)
	if cached, ok := a.cache.Get(prompt); ok {
		var outline Outline
		if err := json.Unmarshal([]byte(cached), &outline); err == nil {
			return outline, nil
		}
	}
	// naive default outline for demo
	outline := Outline{Title: topic, Sections: []string{"Introduction", "System Architecture", "LLM Integration", "Document Output Plan", "Milestones", "Appendix"}}
	raw, _ := json.Marshal(outline)
	a.cache.Set(prompt, string(raw))
	return outline, nil
}

func (a *contentAgent) GenerateSection(ctx context.Context, topic, section string) (string, error) {
	prompt := fmt.Sprintf("Write a detailed section for document topic '%s' for the section '%s'.", topic, section)
	if cached, ok := a.cache.Get(prompt); ok {
		return cached, nil
	}
	// Demo placeholder text
	text := fmt.Sprintf("%s\n\nThis is an auto-generated section for '%s'. Replace with LLM output in production.", section, topic)
	a.cache.Set(prompt, text)
	return text, nil
}
