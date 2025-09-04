package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Simple PocketFlow -> Ollama generator in Go.
// Features:
// - clone a public GitHub repo (or use local dir)
// - scan files by include/exclude patterns and max file size
// - summarize each file with a local LLM (Ollama)
// - synthesize abstractions & relationships
// - generate a multi-chapter tutorial (Markdown) and write to output
//
// Requirements (on your machine):
// - git (in PATH)
// - Ollama running locally (default: http://localhost:11434). Pull the model you want (eg. llama3.2 or codellama) in Ollama.
//
// Example usage:
// go run main.go --repo https://github.com/The-Pocket/PocketFlow-Tutorial-Codebase-Knowledge --model llama3.2 --output ./output

// Configurable defaults
var DefaultInclude = []string{"*.py", "*.js", "*.jsx", "*.ts", "*.tsx", "*.go", "*.java", "*.c", "*.cpp", "*.h", "*.md", "Dockerfile", "Makefile", "*.yaml", "*.yml"}
var DefaultExclude = []string{"assets/", "data/", "images/", "public/", "static/", "temp/", "docs/", "venv/", "node_modules/", "tests/", "examples/", "dist/", "build/", ".git/", ".github/"}

// Ollama API defaults
const DefaultOllamaURL = "http://localhost:11434"

func main() {
	repo := flag.String("repo", "", "GitHub repo URL (public)")
	dir := flag.String("dir", "", "Local directory path to analyze")
	model := flag.String("model", "llama3.2", "Ollama model name (eg. llama3.2 or codellama)")
	ollamaURL := flag.String("ollama-url", DefaultOllamaURL, "Ollama API base URL")
	output := flag.String("output", "output", "Output directory")
	maxSize := flag.Int("max-size", 200000, "Maximum file size (bytes) to read")
	maxAbs := flag.Int("max-abstractions", 15, "Maximum abstractions to identify")
	language := flag.String("language", "english", "Language for tutorial")
	flag.Parse()

	if *repo == "" && *dir == "" {
		log.Fatal("provide either --repo or --dir")
	}

	workDir := ""
	var err error
	if *repo != "" {
		workDir, err = cloneRepo(*repo)
		if err != nil {
			log.Fatalf("failed to clone repo: %v", err)
		}
		defer os.RemoveAll(workDir)
	} else {
		workDir = *dir
	}

	files, err := collectFiles(workDir, DefaultInclude, DefaultExclude, int64(*maxSize))
	if err != nil {
		log.Fatalf("collect files: %v", err)
	}

	if len(files) == 0 {
		log.Fatalf("no files matched patterns in %s", workDir)
	}

	os.MkdirAll(*output, 0755)

	// Summarize files
	fileSummaries := make([]string, 0, len(files))
	fmt.Printf("Found %d files — summarizing with model %s at %s\n", len(files), *model, *ollamaURL)

	for i, f := range files {
		fmt.Printf("[%d/%d] Summarizing %s\n", i+1, len(files), f)
		content, err := os.ReadFile(f)
		if err != nil {
			log.Printf("warning: read %s: %v — skipping", f, err)
			continue
		}
		prompt := buildFileSummaryPrompt(filepath.Base(f), f, string(content), *language)
		resp, err := callOllama(*ollamaURL, *model, prompt)
		if err != nil {
			log.Printf("warning: LLM summary failed for %s: %v — saving fallback summary", f, err)
			resp = fmt.Sprintf("(error summarizing file: %v)", err)
		}
		fileSummaries = append(fileSummaries, fmt.Sprintf("### File: %s\nPath: %s\n\n%s\n---\n", filepath.Base(f), f, resp))
		// be polite to local models — small pause
		time.Sleep(200 * time.Millisecond)
	}

	// Synthesize abstractions
	summaryAggregate := strings.Join(fileSummaries, "\n\n")
	absPrompt := buildAbstractionPrompt(summaryAggregate, *maxAbs, *language)
	fmt.Println("Identifying abstractions and relationships...")
	absResp, err := callOllama(*ollamaURL, *model, absPrompt)
	if err != nil {
		log.Printf("warning: abstraction extraction failed: %v", err)
		absResp = "(error extracting abstractions)"
	}

	// Generate chapter plan
	planPrompt := buildChapterPlanPrompt(absResp, *language)
	planResp, err := callOllama(*ollamaURL, *model, planPrompt)
	if err != nil {
		log.Printf("warning: chapter planning failed: %v", err)
		planResp = "(error generating chapter plan)"
	}

	// Generate final tutorial (detailed)
	finalPrompt := buildFinalTutorialPrompt(summaryAggregate, absResp, planResp, *language)
	fmt.Println("Generating final tutorial (this may take a while)...")
	finalResp, err := callOllama(*ollamaURL, *model, finalPrompt)
	if err != nil {
		log.Fatalf("tutorial generation failed: %v", err)
	}

	outPath := filepath.Join(*output, "TUTORIAL.md")
	if err := os.WriteFile(outPath, []byte(finalResp), 0644); err != nil {
		log.Fatalf("write output: %v", err)
	}

	fmt.Printf("Tutorial written to %s\n", outPath)
}

func cloneRepo(repo string) (string, error) {
	tmp, err := os.MkdirTemp("", "pocketflow_repo_")
	if err != nil {
		return "", err
	}
	fmt.Printf("Cloning %s into %s\n", repo, tmp)
	cmd := exec.Command("git", "clone", "--depth", "1", repo, tmp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	return tmp, nil
}

func collectFiles(root string, includes, excludes []string, maxSize int64) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip
		}
		if info.IsDir() {
			// check excludes
			for _, ex := range excludes {
				if strings.HasPrefix(path, filepath.Join(root, ex)) || strings.Contains(path, ex) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		// Skip large files
		if info.Size() > maxSize {
			return nil
		}
		// Match includes by suffix or glob
		matched := false
		base := filepath.Base(path)
		for _, inc := range includes {
			// very basic matching: suffix or filepath.Match
			if strings.HasSuffix(inc, "/") {
				continue
			}
			if ok, _ := filepath.Match(inc, base); ok {
				matched = true
				break
			}
			if strings.HasSuffix(inc, base) {
				matched = true
				break
			}
			// special case: Dockerfile, Makefile exact names
			if inc == "Dockerfile" || inc == "Makefile" {
				if base == inc {
					matched = true
					break
				}
			}
		}
		if matched {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func buildFileSummaryPrompt(filename, path, content, language string) string {
	fence := "```"
	p := fmt.Sprintf(`You are a code summarization assistant. Produce a clear, developer-friendly summary in %s for the file %s (path: %s).

Instructions:
1) Provide a short "one_line" summary (1 sentence).
2) Provide "purpose" (what this file is responsible for).
3) List major functions/classes (name + 1-line description).
4) List key technical details and any TODOs or important comments.
5) Provide a short example of how this file is used, if applicable.

Respond in Markdown. Begin with "%s - File Summary" as a header, then the sections.

File content:
%s
%s
%s
`,
		language, filename, path, filename, fence, snippet(content, 20000), fence)
	return p
}

func buildAbstractionPrompt(summaries string, maxAbs int, language string) string {
	p := fmt.Sprintf(`You are an assistant that identifies architectural abstractions and relationships for a codebase. Given the following file summaries (Markdown), identify up to %d core abstractions/components (eg: "TaskScheduler", "Database Layer", "API Handlers", "Auth Module") and for each abstraction provide:
- name
- one-line description
- files that implement or relate to it (list)
- dependencies (what other abstractions it talks to)

Also produce a short adjacency list of relationships between abstractions.

Respond in Markdown with clear sections. Language: %s

File summaries:

%s
`, maxAbs, language, summaries)
	return p
}

// func buildChapterPlanPrompt(abstractions string, language string) string {
// 	p := fmt.Sprintf(`You are a tutorial planner. Given these abstractions and relationships (Markdown), create a chapter-by-chapter plan for a beginner-friendly tutorial in %s. For each chapter give:
// 1) Chapter title
// 2) Estimated length in paragraphs (short)
// 3) Key learning objectives
// 4) Files / abstractions covered

// Respond in Markdown with a numbered chapter list.

// Abstractions/relationships:

// %s
// `, language, abstractions)
// 	return p
// }

func buildChapterPlanPrompt(abstractions string, language string) string {
	p := fmt.Sprintf(`You are a tutorial planner. Given these abstractions and relationships (Markdown), create a progressive multi-chapter plan in %s for teaching this codebase to developers.

Rules for the plan:
- Start simple (setup, project overview, basic files).
- Gradually introduce abstractions in logical order.
- Each chapter must have:
  1) Title
  2) Audience level (Beginner, Intermediate, Advanced)
  3) Learning objectives
  4) Files and abstractions covered
  5) Expected outcome after finishing chapter

Respond in Markdown as a numbered chapter list.

Abstractions/relationships:

%s
`, language, abstractions)
	return p
}

func buildFinalTutorialPrompt(summaries, abstractions, plan, language string) string {
	fence := "```"
	p := fmt.Sprintf(`You are a senior engineer writing a multi-chapter tutorial in %s for this codebase.
Base it on the provided chapter plan, file summaries, and abstractions.

Requirements for EACH chapter:
1. Begin with a short introduction (why this chapter matters).
2. Explain key abstractions & files with plain language.
3. Show **real code excerpts** in fenced blocks. Format:
  %s%s
   // code here
   %s
4. Add ASCII diagrams or tables where useful.
5. Provide **step-by-step walkthroughs** (what to run, how to test).
6. Add "Why it works" and "Common pitfalls".
7. End with a short **Hands-on Exercise**.

At the end of the tutorial:
- Provide a "Where to go next" section (ideas to extend the project).
- Add a summary recap of what the reader has learned.

Inputs you MUST use:

File summaries:
%s

Abstractions:
%s

Chapter plan:
%s
`, language, fence, language, fence, summaries, abstractions, plan)
	return p
}

// func buildFinalTutorialPrompt(summaries, abstractions, plan, language string) string {
// 	p := fmt.Sprintf(`You are a senior engineer writing a detailed tutorial for this codebase in %s. Use the following inputs:

// 1) File summaries (use them for accurate code/line references)
// 2) Identified abstractions and relationships
// 3) Chapter plan

// Produce a comprehensive tutorial in Markdown. For each chapter in the plan, produce:
// - A thorough explanation with code excerpts (use triple-backtick fences and include the file path before the fence)
// - Diagrams described in ASCII if helpful
// - Example usage and step-by-step walkthroughs
// - "Why it works" and "Common pitfalls" subsections

// Make the tutorial actionable so a reader can understand architecture, run the project, and extend it. Use the file summaries and abstractions to reference real files. Language: %s

// File summaries:
// %s

// Abstractions:
// %s

// Chapter plan:
// %s
// `, language, language, summaries, abstractions, plan)
// 	return p
// }

func snippet(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n... (truncated)"
}

// callOllama posts to /api/generate with stream=false and returns the response text.
func callOllama(baseURL, model, prompt string) (string, error) {
	url := strings.TrimRight(baseURL, "/") + "/api/generate"
	body := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}
	b, _ := json.Marshal(body)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error: %s: %s", resp.Status, string(data))
	}
	// Ollama returns a JSON object. We'll parse common fields. Response shape may vary by version.
	var parsed map[string]interface{}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		// Sometimes Ollama returns plain text — fallback to raw body
		return string(data), nil
	}
	// Try fields in order: response, message.content, message, result
	if v, ok := parsed["response"]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}
	}
	if v, ok := parsed["message"]; ok {
		// message can be object with content
		switch m := v.(type) {
		case string:
			return m, nil
		case map[string]interface{}:
			if c, ok := m["content"].(string); ok {
				return c, nil
			}
		}
	}
	if v, ok := parsed["error"]; ok {
		return "", errors.New(fmt.Sprint(v))
	}
	// As a last resort return the whole JSON
	pretty := new(bytes.Buffer)
	json.Indent(pretty, data, "", "  ")
	return pretty.String(), nil
}
