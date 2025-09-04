package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/prathyushnallamothu/swarmgo/llm"
)

// Interface for memory providers
type MemoryProvider interface {
	Add(agentName string, role llm.Role, content string) error
	Get(agentName string) ([]llm.Message, error)
}

///////////////////////////////////////////////////////////
// File-based implementation
///////////////////////////////////////////////////////////

type FileMemoryProvider struct {
	mu       sync.Mutex
	basePath string
}

func NewFileMemoryProvider(basePath string) *FileMemoryProvider {
	os.MkdirAll(basePath, 0755)
	log.Printf("🗂️  FileMemoryProvider initialized at %s", basePath)
	return &FileMemoryProvider{basePath: basePath}
}

func (f *FileMemoryProvider) filePath(agentName string) string {
	return filepath.Join(f.basePath, agentName+".json")
}

// internal helper without locks (to avoid deadlocks)
func (f *FileMemoryProvider) load(agentName string) ([]llm.Message, error) {
	path := f.filePath(agentName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("📂 [%s] no memory file, starting fresh", agentName)
			return []llm.Message{}, nil
		}
		return nil, err
	}

	var messages []llm.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}
	log.Printf("📖 [%s] loaded %d messages from memory", agentName, len(messages))
	return messages, nil
}

func (f *FileMemoryProvider) Add(agentName string, role llm.Role, content string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	messages, err := f.load(agentName)
	if err != nil {
		return err
	}

	// Append new entry
	messages = append(messages, llm.Message{Role: role, Content: content})
	log.Printf("✏️  [%s] adding message: role=%s content=%q", agentName, role, truncate(content, 80))

	// Save back to file
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(f.filePath(agentName), data, 0644); err != nil {
		return err
	}

	log.Printf("💾 [%s] saved %d messages to memory", agentName, len(messages))
	return nil
}

func (f *FileMemoryProvider) Get(agentName string) ([]llm.Message, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.load(agentName)
}

///////////////////////////////////////////////////////////
// Helper
///////////////////////////////////////////////////////////

// truncate long content in logs
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
