package embedding

import (
	"regexp"
	"strings"
)

const (
	maxChunkChars     = 4096
	chunkOverlapChars = 256
)

var (
	imagePattern = regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)
	urlPattern   = regexp.MustCompile(`https?://\S+`)
)

func prepareText(title, content string) string {
	text := title + "\n\n" + content
	text = imagePattern.ReplaceAllString(text, "")
	text = urlPattern.ReplaceAllString(text, "")
	return strings.TrimSpace(text)
}

func chunkText(text string) []string {
	if len(text) <= maxChunkChars {
		return []string{text}
	}

	sections := splitByHeadings(text)

	var chunks []string
	for _, section := range sections {
		if len(section) <= maxChunkChars {
			chunks = append(chunks, section)
			continue
		}
		chunks = append(chunks, splitByParagraphs(section)...)
	}

	refined := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		if len(chunk) <= maxChunkChars {
			refined = append(refined, chunk)
			continue
		}
		refined = append(refined, splitByNewlines(chunk)...)
	}

	return applyOverlap(refined)
}

func splitByHeadings(text string) []string {
	lines := strings.Split(text, "\n")
	var sections []string
	var current strings.Builder

	for _, line := range lines {
		if isHeading(line) && current.Len() > 0 {
			sections = append(sections, strings.TrimSpace(current.String()))
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteByte('\n')
		}
		current.WriteString(line)
	}
	if current.Len() > 0 {
		sections = append(sections, strings.TrimSpace(current.String()))
	}

	return sections
}

func splitByParagraphs(text string) []string {
	paragraphs := strings.Split(text, "\n\n")
	return mergeSmallChunks(paragraphs)
}

func splitByNewlines(text string) []string {
	lines := strings.Split(text, "\n")
	return mergeSmallChunks(lines)
}

func mergeSmallChunks(parts []string) []string {
	var chunks []string
	var current strings.Builder

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if current.Len() > 0 && current.Len()+len(part)+1 > maxChunkChars {
			chunks = append(chunks, strings.TrimSpace(current.String()))
			current.Reset()
		}

		if current.Len() > 0 {
			current.WriteByte('\n')
		}
		current.WriteString(part)
	}
	if current.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(current.String()))
	}

	return chunks
}

func applyOverlap(chunks []string) []string {
	if len(chunks) <= 1 {
		return chunks
	}

	result := make([]string, len(chunks))
	result[0] = chunks[0]

	for i := 1; i < len(chunks); i++ {
		prev := chunks[i-1]
		overlap := prev
		if len(overlap) > chunkOverlapChars {
			overlap = overlap[len(overlap)-chunkOverlapChars:]
		}
		result[i] = overlap + "\n" + chunks[i]
	}

	return result
}

func isHeading(line string) bool {
	trimmed := strings.TrimLeft(line, "#")
	return len(trimmed) < len(line) && len(line)-len(trimmed) <= 6 && strings.HasPrefix(trimmed, " ")
}
