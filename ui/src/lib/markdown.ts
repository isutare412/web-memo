/**
 * Convert a header's raw text into a URL-fragment slug.
 *
 * Lowercases, turns whitespace runs into single hyphens, and strips
 * everything that is not a Unicode letter, Unicode number, or hyphen — so
 * Korean/CJK headers stay readable. Returns "section" when nothing usable
 * remains (e.g. an emoji-only header).
 */
export function slugify(raw: string): string {
  const slug = raw
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-')
    .replace(/[^\p{L}\p{N}-]+/gu, '')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '')

  return slug || 'section'
}

/**
 * Create a stateful slugger that guarantees unique slugs within one document.
 *
 * The first time a base slug is seen it is returned as-is; later collisions
 * (including against explicit "-N" headers) get the next free "-1", "-2", ...
 * suffix. Construct a fresh slugger per parse so counts reset per render.
 */
export function makeSlugger(): (raw: string) => string {
  const used = new Set<string>()

  return (raw: string): string => {
    const base = slugify(raw)
    let candidate = base
    let i = 1
    while (used.has(candidate)) {
      candidate = `${base}-${i}`
      i++
    }
    used.add(candidate)
    return candidate
  }
}
