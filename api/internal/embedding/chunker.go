package embedding

import (
	"regexp"
	"strings"
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

func chunkText(text string, prof modelProfile) []string {
	if len(text) <= prof.maxChunkBytes {
		return []string{text}
	}

	sections := splitByHeadings(text)

	var chunks []string
	for _, section := range sections {
		if len(section) <= prof.maxChunkBytes {
			chunks = append(chunks, section)
			continue
		}
		chunks = append(chunks, splitByParagraphs(section, prof.maxChunkBytes)...)
	}

	refined := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		if len(chunk) <= prof.maxChunkBytes {
			refined = append(refined, chunk)
			continue
		}
		refined = append(refined, splitByNewlines(chunk, prof.maxChunkBytes)...)
	}

	return applyOverlap(refined, prof.chunkOverlapBytes)
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

func splitByParagraphs(text string, maxBytes int) []string {
	paragraphs := strings.Split(text, "\n\n")
	return mergeSmallChunks(paragraphs, maxBytes)
}

func splitByNewlines(text string, maxBytes int) []string {
	lines := strings.Split(text, "\n")
	return mergeSmallChunks(lines, maxBytes)
}

func mergeSmallChunks(parts []string, maxBytes int) []string {
	var chunks []string
	var current strings.Builder

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if current.Len() > 0 && current.Len()+len(part)+1 > maxBytes {
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

func applyOverlap(chunks []string, overlapBytes int) []string {
	if len(chunks) <= 1 {
		return chunks
	}

	result := make([]string, len(chunks))
	result[0] = chunks[0]

	for i := 1; i < len(chunks); i++ {
		prev := chunks[i-1]
		overlap := prev
		if len(overlap) > overlapBytes {
			overlap = overlap[len(overlap)-overlapBytes:]
		}
		result[i] = overlap + "\n" + chunks[i]
	}

	return result
}

func isHeading(line string) bool {
	trimmed := strings.TrimLeft(line, "#")
	return len(trimmed) < len(line) && len(line)-len(trimmed) <= 6 && strings.HasPrefix(trimmed, " ")
}
