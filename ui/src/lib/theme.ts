const themeKey = 'preferredTheme'

export enum ThemeMode {
  Light = 'light',
  Dark = 'dark',
}

const themeNames = new Map([
  [ThemeMode.Light, 'emerald'],
  [ThemeMode.Dark, 'dracula'],
])

export function getPreferredTheme(): ThemeMode {
  const theme = localStorage.getItem(themeKey)
  switch (theme) {
    case 'light':
      return ThemeMode.Light
    case 'dark':
      return ThemeMode.Dark
  }

  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
    ? ThemeMode.Dark
    : ThemeMode.Light
}

export function setPreferredTheme(theme: ThemeMode) {
  localStorage.setItem(themeKey, theme)
}

export function setDocumentDataTheme(theme: ThemeMode) {
  document.documentElement.setAttribute('data-theme', themeNames.get(theme)!)
}
