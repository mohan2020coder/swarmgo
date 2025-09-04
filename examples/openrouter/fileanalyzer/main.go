package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	dotenv "github.com/joho/godotenv"
	swarmgo "github.com/mohan2020coder/swarmgo"
	"github.com/mohan2020coder/swarmgo/llm"
)

// Configurable defaults
var DefaultInclude = []string{"*.py", "*.js", "*.jsx", "*.ts", "*.tsx", "*.go", "*.java", "*.c", "*.cpp", "*.h", "*.md", "Dockerfile", "Makefile", "*.yaml", "*.yml"}
var DefaultExclude = []string{"assets/", "data/", "images/", "public/", "static/", "temp/", "docs/", "venv/", "node_modules/", "tests/", "examples/", "dist/", "build/", ".git/", ".github/"}

func main() {
	repo := flag.String("repo", "", "GitHub repo URL (public)")
	dir := flag.String("dir", "", "Local directory path to analyze")
	model := flag.String("model", "deepseek/deepseek-chat", "OpenRouter model slug (see https://openrouter.ai/models)")
	output := flag.String("output", "output", "Output directory")
	maxSize := flag.Int("max-size", 200000, "Maximum file size (bytes) to read")
	maxAbs := flag.Int("max-abstractions", 15, "Maximum abstractions to identify")
	language := flag.String("language", "english", "Language for tutorial")
	flag.Parse()

	if *repo == "" && *dir == "" {
		log.Fatal("provide either --repo or --dir")
	}

	// load env + init client
	dotenv.Load()
	client := swarmgo.NewSwarm(os.Getenv("OPENROUTER_API_KEY"), llm.OpenRouter)
	agent := &swarmgo.Agent{
		Name:         "TutorAgent",
		Instructions: "You are a helpful agent that generates summaries and tutorials.",
		Model:        *model,
	}

	// prepare repo
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
	fmt.Printf("Found %d files — summarizing with model %s (OpenRouter)\n", len(files), *model)

	for i, f := range files {
		fmt.Printf("[%d/%d] Summarizing %s\n", i+1, len(files), f)
		content, err := os.ReadFile(f)
		if err != nil {
			log.Printf("warning: read %s: %v — skipping", f, err)
			continue
		}
		prompt := buildFileSummaryPrompt(filepath.Base(f), f, string(content), *language)
		resp, err := callSwarm(client, agent, prompt)
		if err != nil {
			log.Printf("warning: LLM summary failed for %s: %v — saving fallback summary", f, err)
			resp = fmt.Sprintf("(error summarizing file: %v)", err)
		}
		fileSummaries = append(fileSummaries, fmt.Sprintf("### File: %s\nPath: %s\n\n%s\n---\n", filepath.Base(f), f, resp))
		time.Sleep(200 * time.Millisecond) // be polite
	}

	// Synthesize abstractions
	summaryAggregate := strings.Join(fileSummaries, "\n\n")
	absPrompt := buildAbstractionPrompt(summaryAggregate, *maxAbs, *language)
	fmt.Println("Identifying abstractions and relationships...")
	absResp, err := callSwarm(client, agent, absPrompt)
	if err != nil {
		log.Printf("warning: abstraction extraction failed: %v", err)
		absResp = "(error extracting abstractions)"
	}

	// Generate chapter plan
	planPrompt := buildChapterPlanPrompt(absResp, *language)
	planResp, err := callSwarm(client, agent, planPrompt)
	if err != nil {
		log.Printf("warning: chapter planning failed: %v", err)
		planResp = "(error generating chapter plan)"
	}

	// Generate final tutorial (detailed)
	finalPrompt := buildFinalTutorialPrompt(summaryAggregate, absResp, planResp, *language)
	fmt.Println("Generating final tutorial (this may take a while)...")
	finalResp, err := callSwarm(client, agent, finalPrompt)
	if err != nil {
		log.Fatalf("tutorial generation failed: %v", err)
	}

	outPath := filepath.Join(*output, "TUTORIAL.md")
	if err := os.WriteFile(outPath, []byte(finalResp), 0644); err != nil {
		log.Fatalf("write output: %v", err)
	}

	fmt.Printf("Tutorial written to %s\n", outPath)
}

// ------------------- helpers -------------------

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
			return nil
		}
		if info.IsDir() {
			for _, ex := range excludes {
				if strings.HasPrefix(path, filepath.Join(root, ex)) || strings.Contains(path, ex) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if info.Size() > maxSize {
			return nil
		}
		base := filepath.Base(path)
		for _, inc := range includes {
			if strings.HasSuffix(inc, "/") {
				continue
			}
			if ok, _ := filepath.Match(inc, base); ok {
				files = append(files, path)
				return nil
			}
			if strings.HasSuffix(inc, base) {
				files = append(files, path)
				return nil
			}
			if inc == "Dockerfile" || inc == "Makefile" {
				if base == inc {
					files = append(files, path)
					return nil
				}
			}
		}
		return nil
	})
	return files, err
}

func buildFileSummaryPrompt(filename, path, content, language string) string {
	fence := "```"
	return fmt.Sprintf(`You are a code summarization assistant. Produce a clear, developer-friendly summary in %s for the file %s (path: %s).

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
`, language, filename, path, filename, fence, snippet(content, 20000), fence)
}

func buildAbstractionPrompt(summaries string, maxAbs int, language string) string {
	return fmt.Sprintf(`You are an assistant that identifies architectural abstractions and relationships for a codebase. Given the following file summaries (Markdown), identify up to %d core abstractions/components (eg: "TaskScheduler", "Database Layer", "API Handlers", "Auth Module") and for each abstraction provide:
- name
- one-line description
- files that implement or relate to it (list)
- dependencies (what other abstractions it talks to)

Also produce a short adjacency list of relationships between abstractions.

Respond in Markdown with clear sections. Language: %s

File summaries:

%s
`, maxAbs, language, summaries)
}

func buildChapterPlanPrompt(abstractions string, language string) string {
	return fmt.Sprintf(`You are a tutorial planner. Given these abstractions and relationships (Markdown), create a progressive multi-chapter plan in %s for teaching this codebase to developers.

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
}

func buildFinalTutorialPrompt(summaries, abstractions, plan, language string) string {
	fence := "```"
	return fmt.Sprintf(`You are a senior engineer writing a multi-chapter tutorial in %s for this codebase.
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
}

func snippet(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n... (truncated)"
}

// ------------------- swarm wrapper -------------------

func callSwarm(client *swarmgo.Swarm, agent *swarmgo.Agent, prompt string) (string, error) {
	ctx := context.Background()
	msgs := []llm.Message{
		{Role: llm.RoleUser, Content: prompt},
	}

	resp, err := client.Run(ctx, agent, msgs, nil, "", false, false, 5, true)
	if err != nil {
		return "", err
	}

	if len(resp.Messages) == 0 {
		return "", fmt.Errorf("no messages returned")
	}

	return resp.Messages[len(resp.Messages)-1].Content, nil
}
