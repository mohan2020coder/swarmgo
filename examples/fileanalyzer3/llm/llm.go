package llm

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Thread-safe cache
var (
	cache       = make(map[string]string)
	cacheLoaded = false
	cacheMutex  sync.Mutex
	CacheFile   = "llm_cache.json"
	LogDir      = "logs"
)

// Load cache from disk
func loadCache() {
	if cacheLoaded {
		return
	}
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	data, err := os.ReadFile(CacheFile)
	if err == nil {
		_ = json.Unmarshal(data, &cache)
	}
	cacheLoaded = true
}

// Save cache to disk
func saveCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	data, err := json.MarshalIndent(cache, "", "  ")
	if err == nil {
		_ = os.WriteFile(CacheFile, data, 0644)
	}
}

// Generate hash key for prompt + model
func hashPrompt(model, prompt string) string {
	h := sha1.New()
	io.WriteString(h, model+"|"+prompt)
	return hex.EncodeToString(h.Sum(nil))
}

// Log prompt & response to daily log
func logPrompt(model, prompt, response string) {
	dir := os.Getenv("LOG_DIR")
	if dir != "" {
		LogDir = dir
	}
	_ = os.MkdirAll(LogDir, 0755)
	file := filepath.Join(LogDir, time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open log file: %v", err)
		return
	}
	defer f.Close()
	entry := "\n=== [" + time.Now().Format(time.RFC3339) + "] MODEL: " + model + " ===\n" +
		"PROMPT:\n" + prompt + "\n\nRESPONSE:\n" + response + "\n---\n"
	f.WriteString(entry)
}

// APICaller type
type APICaller func(prompt string) (string, error)

// CachedCall wraps LLM call with cache & logging
func CachedCall(model string, prompt string, apiType string, call APICaller) (string, error) {
	loadCache()
	key := hashPrompt(model+"|"+apiType, prompt)

	// Cache hit
	cacheMutex.Lock()
	if resp, ok := cache[key]; ok {
		cacheMutex.Unlock()
		log.Printf("[CACHE HIT] %s", key)
		return resp, nil
	}
	cacheMutex.Unlock()

	// Call LLM
	resp, err := call(prompt)
	if err != nil {
		return "", err
	}

	// Save cache & log
	cacheMutex.Lock()
	cache[key] = resp
	cacheMutex.Unlock()
	saveCache()
	logPrompt(model, prompt, resp)

	return resp, nil
}
