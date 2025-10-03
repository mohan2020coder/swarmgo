package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"pocketflow/llm"
	"pocketflow/prompt"
	"pocketflow/repo"
	"pocketflow/utils"
)

func main() {
	// CLI flags (same as original)
	repoURL := flag.String("repo", "", "GitHub repo URL (public)")
	dir := flag.String("dir", "", "Local directory path to analyze")
	model := flag.String("model", "llama3.2", "Ollama model name (eg. llama3.2 or codellama)")
	ollamaURL := flag.String("ollama-url", "http://localhost:11434", "Ollama API base URL")
	output := flag.String("output", "output", "Output directory")
	maxSize := flag.Int("max-size", 200000, "Maximum file size (bytes) to read")
	maxAbs := flag.Int("max-abstractions", 15, "Maximum abstractions to identify")
	language := flag.String("language", "english", "Language for tutorial")
	apiChoice := flag.String("api", "ollama", "Which API to use: 'ollama' or 'longcat'")
	longcatKey := flag.String("api-key", "", "API key for LongCat API")

	flag.Parse()

	if *repoURL == "" && *dir == "" {
		log.Fatal("Provide either --repo or --dir")
	}

	if *apiChoice != "ollama" && *apiChoice != "longcat" {
		log.Fatalf("Invalid --api value: %s; must be 'ollama' or 'longcat'", *apiChoice)
	}

	if *apiChoice == "longcat" && *longcatKey == "" {
		log.Fatal("--api-key is required when using LongCat API")
	}

	// Clone repo or use local dir
	workDir := *dir
	if *repoURL != "" {
		tmp, err := repo.CloneRepo(*repoURL)
		if err != nil {
			log.Fatalf("Failed to clone repo: %v", err)
		}
		defer os.RemoveAll(tmp)
		workDir = tmp
	}

	// Collect files
	files, err := repo.CollectFiles(workDir,
		[]string{"*.go", "*.py", "*.js", "*.ts", "*.tsx", "*.java", "*.c", "*.cpp", "*.h", "*.md", "Dockerfile", "Makefile", "*.yaml", "*.yml"},
		[]string{"assets/", "data/", "images/", "public/", "static/", "temp/", "docs/", "venv/", "node_modules/", "tests/", "examples/", "dist/", "build/", ".git/", ".github/"},
		int64(*maxSize),
	)
	if err != nil || len(files) == 0 {
		log.Fatalf("No files matched patterns in %s", workDir)
	}

	os.MkdirAll(*output, 0755)
	fmt.Printf("Found %d files — summarizing with model %s using %s API\n", len(files), *model, *apiChoice)

	// 1️⃣ Summarize files
	fileSummaries := []string{}
	for i, f := range files {
		content, _ := os.ReadFile(f)
		promptText := prompt.BuildFileSummaryPrompt(filepath.Base(f), f, string(content), *language)

		resp, err := llm.CachedCall(*model, promptText, *apiChoice, func(p string) (string, error) {
			if *apiChoice == "ollama" {
				return callOllama(*ollamaURL, *model, p)
			}
			return callLongCat(*longcatKey, p)
		})
		if err != nil {
			resp = fmt.Sprintf("(error summarizing file: %v)", err)
		}

		outFile := filepath.Join(*output, fmt.Sprintf("summary_%02d_%s.md", i+1, utils.SanitizeFilename(filepath.Base(f))))
		_ = os.WriteFile(outFile, []byte(resp), 0644)
		fileSummaries = append(fileSummaries, fmt.Sprintf("### File: %s\nPath: %s\n\n%s\n---\n", filepath.Base(f), f, resp))
		time.Sleep(300 * time.Millisecond)
	}

	// 2️⃣ Generate abstractions
	absPrompt := prompt.BuildAbstractionPrompt(joinSummaries(fileSummaries), *maxAbs, *language)
	absResp, _ := llm.CachedCall(*model, absPrompt, *apiChoice, func(p string) (string, error) {
		if *apiChoice == "ollama" {
			return callOllama(*ollamaURL, *model, p)
		}
		return callLongCat(*longcatKey, p)
	})
	_ = os.WriteFile(filepath.Join(*output, "abstractions.md"), []byte(absResp), 0644)

	// 3️⃣ Chapter plan
	planPrompt := prompt.BuildChapterPlanPrompt(absResp, *language)
	planResp, _ := llm.CachedCall(*model, planPrompt, *apiChoice, func(p string) (string, error) {
		if *apiChoice == "ollama" {
			return callOllama(*ollamaURL, *model, p)
		}
		return callLongCat(*longcatKey, p)
	})
	planFile := filepath.Join(*output, "chapter_plan.md")
	_ = os.WriteFile(planFile, []byte(planResp), 0644)

	chapters := utils.ParseChapters(planResp)
	if len(chapters) == 0 {
		log.Printf("No chapters parsed from plan; check %s for raw output", planFile)
		return
	}

	// 4️⃣ Generate chapters
	tutorialParts := []string{}
	for i, chapter := range chapters {
		fmt.Printf("Generating chapter %d/%d\n", i+1, len(chapters))
		chPrompt := prompt.BuildChapterContentPrompt(chapter, joinSummaries(fileSummaries), absResp, *language)

		content, _ := llm.CachedCall(*model, chPrompt, *apiChoice, func(p string) (string, error) {
			if *apiChoice == "ollama" {
				return callOllama(*ollamaURL, *model, p)
			}
			return callLongCat(*longcatKey, p)
		})

		chapterFile := filepath.Join(*output, fmt.Sprintf("chapter_%02d.md", i+1))
		_ = os.WriteFile(chapterFile, []byte(content), 0644)
		tutorialParts = append(tutorialParts, content)
		time.Sleep(500 * time.Millisecond)
	}

	fullTutorial := joinSummaries(tutorialParts)
	outPath := filepath.Join(*output, "TUTORIAL.md")
	_ = os.WriteFile(outPath, []byte(fullTutorial), 0644)
	fmt.Printf("Tutorial generated successfully at %s\n", outPath)
}

// Helpers
func joinSummaries(parts []string) string {
	res := ""
	for i, s := range parts {
		if i > 0 {
			res += "\n\n"
		}
		res += s
	}
	return res
}
