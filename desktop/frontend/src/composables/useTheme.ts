import { ref } from 'vue'

export type ThemeMode = 'auto' | 'light' | 'dark'

const STORAGE_KEY = 'ccx-desktop-theme'

const theme = ref<ThemeMode>('auto')

const mql = window.matchMedia('(prefers-color-scheme: dark)')

function resolvedDark(mode: ThemeMode): boolean {
  if (mode === 'auto') return mql.matches
  return mode === 'dark'
}

function applyTheme(mode: ThemeMode) {
  const html = document.documentElement
  if (resolvedDark(mode)) {
    html.classList.add('dark')
  } else {
    html.classList.remove('dark')
  }
}

function onSystemChange() {
  if (theme.value === 'auto') applyTheme('auto')
}

function init() {
  const stored = localStorage.getItem(STORAGE_KEY) as ThemeMode | null
  theme.value = stored === 'light' || stored === 'dark' || stored === 'auto' ? stored : 'auto'
  applyTheme(theme.value)
  mql.addEventListener('change', onSystemChange)
}

function setTheme(mode: ThemeMode) {
  theme.value = mode
  applyTheme(mode)
  localStorage.setItem(STORAGE_KEY, mode)
}

export function useTheme() {
  return { theme, init, setTheme }
}
