const themeKey = 'preferredTheme'

export const themes = [
  'light',
  'dark',
  'cupcake',
  'retro',
  'halloween',
  'sunset',
  'dim',
  'lemonade',
] as const

export type Theme = (typeof themes)[number]

export function getPreferredTheme(): Theme {
  const stored = localStorage.getItem(themeKey)
  if (stored && themes.includes(stored as Theme)) {
    return stored as Theme
  }

  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light'
}

export function setPreferredTheme(theme: Theme) {
  localStorage.setItem(themeKey, theme)
}

export function setDocumentDataTheme(theme: Theme) {
  document.documentElement.setAttribute('data-theme', theme)
}
