import { describe, expect, it } from 'vitest'
import { slugify } from './markdown'

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
