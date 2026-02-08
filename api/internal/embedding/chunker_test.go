package embedding

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Chunker", func() {
	Context("prepareText", func() {
		It("combines title and content", func() {
			result := prepareText("My Title", "Some content")
			Expect(result).To(Equal("My Title\n\nSome content"))
		})

		It("strips image links", func() {
			result := prepareText("Title", "before ![alt text](https://example.com/img.png) after")
			Expect(result).To(Equal("Title\n\nbefore  after"))
		})

		It("strips image links with empty alt text", func() {
			result := prepareText("Title", "before ![](https://example.com/img.png) after")
			Expect(result).To(Equal("Title\n\nbefore  after"))
		})

		It("strips raw URLs", func() {
			result := prepareText("Title", "visit https://example.com/page and http://foo.bar/baz for more")
			Expect(result).To(Equal("Title\n\nvisit  and  for more"))
		})

		It("strips both image links and raw URLs", func() {
			result := prepareText("Title", "![img](https://a.com/b.png)\nhttps://c.com/d")
			Expect(result).To(Equal("Title"))
		})

		It("trims whitespace from result", func() {
			result := prepareText("Title", "  \n  ")
			Expect(result).To(Equal("Title"))
		})
	})

	Context("isHeading", func() {
		DescribeTable("detects markdown headings",
			func(line string, expected bool) {
				Expect(isHeading(line)).To(Equal(expected))
			},
			Entry("h1", "# Heading", true),
			Entry("h2", "## Heading", true),
			Entry("h3", "### Heading", true),
			Entry("h4", "#### Heading", true),
			Entry("h5", "##### Heading", true),
			Entry("h6", "###### Heading", true),
			Entry("too many hashes", "####### Heading", false),
			Entry("no space after hash", "#NoSpace", false),
			Entry("empty line", "", false),
			Entry("plain text", "just text", false),
			Entry("hash in middle", "not a # heading", false),
		)
	})

	Context("splitByHeadings", func() {
		It("splits text at heading boundaries", func() {
			text := "intro\n## Section A\ncontent a\n## Section B\ncontent b"
			sections := splitByHeadings(text)
			Expect(sections).To(Equal([]string{
				"intro",
				"## Section A\ncontent a",
				"## Section B\ncontent b",
			}))
		})

		It("returns single section when no headings", func() {
			text := "just some text\nwith lines"
			sections := splitByHeadings(text)
			Expect(sections).To(Equal([]string{"just some text\nwith lines"}))
		})

		It("handles heading at the start", func() {
			text := "## First\ncontent"
			sections := splitByHeadings(text)
			Expect(sections).To(Equal([]string{"## First\ncontent"}))
		})

		It("handles multiple heading levels", func() {
			text := "# Top\ncontent\n### Sub\nmore"
			sections := splitByHeadings(text)
			Expect(sections).To(Equal([]string{
				"# Top\ncontent",
				"### Sub\nmore",
			}))
		})
	})

	Context("mergeSmallChunks", func() {
		It("merges small parts into chunks within limit", func() {
			parts := []string{"aaa", "bbb", "ccc"}
			chunks := mergeSmallChunks(parts)
			Expect(chunks).To(Equal([]string{"aaa\nbbb\nccc"}))
		})

		It("skips empty parts", func() {
			parts := []string{"aaa", "", "  ", "bbb"}
			chunks := mergeSmallChunks(parts)
			Expect(chunks).To(Equal([]string{"aaa\nbbb"}))
		})

		It("splits when accumulated size exceeds limit", func() {
			large := strings.Repeat("x", maxChunkChars)
			parts := []string{large, "overflow"}
			chunks := mergeSmallChunks(parts)
			Expect(chunks).To(HaveLen(2))
			Expect(chunks[0]).To(Equal(large))
			Expect(chunks[1]).To(Equal("overflow"))
		})
	})

	Context("applyOverlap", func() {
		It("returns single chunk unchanged", func() {
			chunks := applyOverlap([]string{"only one"})
			Expect(chunks).To(Equal([]string{"only one"}))
		})

		It("returns empty slice unchanged", func() {
			chunks := applyOverlap(nil)
			Expect(chunks).To(BeNil())
		})

		It("prepends overlap from previous chunk", func() {
			prev := strings.Repeat("a", 300)
			curr := "current chunk"
			chunks := applyOverlap([]string{prev, curr})

			Expect(chunks).To(HaveLen(2))
			Expect(chunks[0]).To(Equal(prev))

			expectedOverlap := prev[len(prev)-chunkOverlapChars:]
			Expect(chunks[1]).To(Equal(expectedOverlap + "\n" + curr))
		})

		It("uses full previous chunk when shorter than overlap size", func() {
			prev := "short"
			curr := "next"
			chunks := applyOverlap([]string{prev, curr})

			Expect(chunks).To(HaveLen(2))
			Expect(chunks[0]).To(Equal("short"))
			Expect(chunks[1]).To(Equal("short\nnext"))
		})
	})

	Context("chunkText", func() {
		It("returns single chunk for short text", func() {
			text := "short text"
			chunks := chunkText(text)
			Expect(chunks).To(Equal([]string{"short text"}))
		})

		It("returns single chunk at exactly max size", func() {
			text := strings.Repeat("x", maxChunkChars)
			chunks := chunkText(text)
			Expect(chunks).To(Equal([]string{text}))
		})

		It("splits long text by headings", func() {
			section1 := "## Section 1\n" + strings.Repeat("a", 3000)
			section2 := "## Section 2\n" + strings.Repeat("b", 3000)
			text := section1 + "\n" + section2
			chunks := chunkText(text)

			Expect(len(chunks)).To(BeNumerically(">=", 2))
			Expect(chunks[0]).To(ContainSubstring("Section 1"))
		})

		It("splits long section by paragraphs", func() {
			para1 := strings.Repeat("a", 2000)
			para2 := strings.Repeat("b", 2000)
			para3 := strings.Repeat("c", 2000)
			text := para1 + "\n\n" + para2 + "\n\n" + para3
			chunks := chunkText(text)

			Expect(len(chunks)).To(BeNumerically(">=", 2))
		})

		It("returns single chunk for short Korean text", func() {
			text := "안녕하세요. 이것은 한국어 테스트입니다."
			chunks := chunkText(text)
			Expect(chunks).To(Equal([]string{text}))
		})

		It("splits long Korean text exceeding byte limit", func() {
			// Each Korean char is 3 bytes in UTF-8.
			// Build lines of Korean text that together exceed maxChunkChars.
			koreanLine := strings.Repeat("가", 100) // 300 bytes per line
			var lines []string
			for range 20 { // 20 × 300 = 6000 bytes > maxChunkChars
				lines = append(lines, koreanLine)
			}
			text := strings.Join(lines, "\n")
			chunks := chunkText(text)
			Expect(len(chunks)).To(BeNumerically(">=", 2))
		})

		It("splits Korean text by Korean headings", func() {
			section1 := "## 소개\n" + strings.Repeat("가", 1200)
			section2 := "## 본론\n" + strings.Repeat("나", 1200)
			text := section1 + "\n" + section2
			chunks := chunkText(text)

			Expect(len(chunks)).To(BeNumerically(">=", 2))
			Expect(chunks[0]).To(ContainSubstring("소개"))
		})

		It("splits Korean paragraphs", func() {
			para1 := strings.Repeat("가", 700)
			para2 := strings.Repeat("나", 700)
			para3 := strings.Repeat("다", 700)
			text := para1 + "\n\n" + para2 + "\n\n" + para3
			chunks := chunkText(text)

			Expect(len(chunks)).To(BeNumerically(">=", 2))
		})

		It("handles mixed Korean and English content", func() {
			section1 := "## Overview 개요\n" + strings.Repeat("가", 1200)
			section2 := "## Details 상세\n" + strings.Repeat("a", 3000)
			text := section1 + "\n" + section2
			chunks := chunkText(text)

			Expect(len(chunks)).To(BeNumerically(">=", 2))
			Expect(chunks[0]).To(ContainSubstring("개요"))
		})

		It("splits long paragraph by newlines", func() {
			var lines []string
			for i := 0; i < 100; i++ {
				lines = append(lines, strings.Repeat("x", 100))
			}
			text := strings.Join(lines, "\n")
			chunks := chunkText(text)

			Expect(len(chunks)).To(BeNumerically(">=", 2))
			for _, chunk := range chunks {
				// Each chunk (before overlap) should be within limit.
				// After overlap, it may slightly exceed, which is expected.
				Expect(len(chunk)).To(BeNumerically("<=", maxChunkChars+chunkOverlapChars+1))
			}
		})
	})
})
