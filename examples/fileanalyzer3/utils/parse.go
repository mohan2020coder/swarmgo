package utils

import (
	"regexp"
	"strings"
)

func ParseChapters(plan string) []string {
	lines := strings.Split(plan, "\n")
	chapterRE := regexp.MustCompile(`(?i)^##\s*\**chapter\s*\d+.*`)
	chapters := []string{}

	for i := 0; i < len(lines); i++ {
		if chapterRE.MatchString(lines[i]) {
			var sb strings.Builder
			sb.WriteString(lines[i])
			sb.WriteString("\n")

			for j := i + 1; j < len(lines); j++ {
				if chapterRE.MatchString(lines[j]) {
					break
				}
				if strings.TrimSpace(lines[j]) == "" && sb.Len() > 0 {
					break
				}
				sb.WriteString(lines[j])
				sb.WriteString("\n")
			}
			chapters = append(chapters, sb.String())
		}
	}
	return chapters
}
