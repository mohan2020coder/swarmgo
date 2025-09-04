package main

import (
	"fmt"
	"log"
	"os"
)

func WriteCodeFile(lang, code string) {
	ext := map[string]string{
		"python": ".py",
	}[lang]

	filename := "main" + ext
	err := os.WriteFile(filename, []byte(code), 0644)
	if err != nil {
		log.Printf("[Writer] ERROR writing %s: %v\n", filename, err)
		return
	}
	logStep("Writer", fmt.Sprintf("Code written to %s", filename))
}

func WriteDocs(file, section, content string) {
	_ = os.WriteFile(file, []byte(fmt.Sprintf("## [%s]\n\n%s\n\n", section, content)), 0644)
	logStep("Writer", fmt.Sprintf("%s created with section %s", file, section))
}

func AppendDocs(file, section, content string) {
	f, _ := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("## [%s]\n\n%s\n\n", section, content))
	logStep("Writer", fmt.Sprintf("%s appended with section %s", file, section))
}

func WriteTestReport(file string, result *RunResult) {
	content := fmt.Sprintf(`# Execution Report

**Language:** %s  
**Exit Code:** %d  

### Stdout
%s

### Stderr
%s
`, result.Language, result.ExitCode, result.Stdout, result.Stderr)

	_ = os.WriteFile(file, []byte(content), 0644)
	logStep("Writer", fmt.Sprintf("Execution report written to %s", file))
}
