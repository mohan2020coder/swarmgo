package prompt

import (
	"fmt"
)

// File summary prompt
func BuildFileSummaryPrompt(filename, path, content, language string) string {
	fence := "```"
	return fmt.Sprintf(`You are a code summarization assistant. Produce a clear, developer-friendly summary in %s for the file %s (path: %s).

Instructions:
1) Provide a short "one_line" summary.
2) Purpose of the file.
3) Major functions/classes.
4) Key technical details & TODOs.
5) Short usage example if applicable.

Respond in Markdown. Begin with "%s - File Summary" as a header.

%s
%s
%s
`,
		language, filename, path, filename, fence, snippet(content, 20000), fence)
}

// Abstraction prompt
func BuildAbstractionPrompt(summaries string, maxAbs int, language string) string {
	return fmt.Sprintf(`You are an assistant that identifies architectural abstractions and relationships for a codebase. Identify up to %d core abstractions/components and for each provide:
- name
- one-line description
- files
- dependencies
- responsibilities

Output as Markdown. Use %s language.

File summaries:
%s
`, maxAbs, language, summaries)
}

// Chapter plan prompt
func BuildChapterPlanPrompt(abstractions, language string) string {
	return fmt.Sprintf(`You are writing a tutorial plan for a codebase. Generate a detailed chapter outline in %s language.

Abstractions:
%s
`, language, abstractions)
}

// Chapter content prompt
func BuildChapterContentPrompt(chapter, summaries, abstractions, language string) string {
	return fmt.Sprintf(`You are writing a developer tutorial in %s language. Use the chapter plan below plus the file summaries and architectural abstractions.

Chapter plan:
%s

File summaries:
%s

Architectural abstractions:
%s

Write a clear and concise tutorial chapter.
`,
		language, chapter, summaries, abstractions)
}

// snippet limits file content
func snippet(text string, max int) string {
	if len(text) > max {
		return text[:max] + "\n... [truncated]"
	}
	return text
}
