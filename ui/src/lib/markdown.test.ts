import { describe, expect, it } from 'vitest'
import { Marked } from '@ts-stack/markdown'
import { createHeadingRenderer, makeSlugger, slugify } from './markdown'

describe('slugify', () => {
  it('kebab-cases an ASCII header', () => {
    expect(slugify('Hello World')).toBe('hello-world')
  })

  it('preserves Korean characters', () => {
    expect(slugify('회의 노트')).toBe('회의-노트')
  })

  it('handles mixed ASCII and CJK', () => {
    expect(slugify('Project 프로젝트 2')).toBe('project-프로젝트-2')
  })

  it('strips punctuation but keeps letters and numbers', () => {
    expect(slugify('Hello, World! (v2)')).toBe('hello-world-v2')
  })

  it('collapses whitespace runs and trims edges', () => {
    expect(slugify('  Multiple   Spaces  ')).toBe('multiple-spaces')
  })

  it('falls back to "section" for slug-less input', () => {
    expect(slugify('🎉')).toBe('section')
  })
})

describe('makeSlugger', () => {
  it('returns the base slug for the first occurrence', () => {
    const slug = makeSlugger()
    expect(slug('Notes')).toBe('notes')
  })

  it('suffixes duplicates with -1, -2, ...', () => {
    const slug = makeSlugger()
    expect(slug('Notes')).toBe('notes')
    expect(slug('Notes')).toBe('notes-1')
    expect(slug('Notes')).toBe('notes-2')
  })

  it('avoids colliding with an explicit suffixed header', () => {
    const slug = makeSlugger()
    expect(slug('Notes')).toBe('notes')
    expect(slug('Notes 1')).toBe('notes-1')
    expect(slug('Notes')).toBe('notes-2')
  })

  it('de-duplicates the "section" fallback', () => {
    const slug = makeSlugger()
    expect(slug('🎉')).toBe('section')
    expect(slug('🚀')).toBe('section-1')
  })
})

describe('createHeadingRenderer', () => {
  it('adds slug ids to headings', () => {
    const html = Marked.parse('# Hello World', { renderer: createHeadingRenderer() })
    expect(html).toContain('<h1 id="hello-world">')
  })

  it('produces unique ids for duplicate headers', () => {
    const html = Marked.parse('# Notes\n\n# Notes', { renderer: createHeadingRenderer() })
    expect(html).toContain('<h1 id="notes">')
    expect(html).toContain('<h1 id="notes-1">')
  })

  it('keeps Korean headers readable', () => {
    const html = Marked.parse('## 회의 노트', { renderer: createHeadingRenderer() })
    expect(html).toContain('<h2 id="회의-노트">')
  })

  it('preserves inline markup in heading text', () => {
    const html = Marked.parse('# Use `code` here', { renderer: createHeadingRenderer() })
    expect(html).toContain('id="use-code-here"')
    expect(html).toContain('<code>code</code>')
  })
})
