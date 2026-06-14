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
