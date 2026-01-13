import {
  computed,
  ref,
  watch,
} from 'vue'
import type { AppSettings } from '@/types/settings'
import { DEFAULT_SETTINGS } from '@/types/settings'
import { generateColorScheme } from '@/utils/colors'

const STORAGE_KEY = 'round-sound-settings'

// Global reactive state
const settings = ref<AppSettings>(loadSettings())

/**
 * Load settings from localStorage
 */
function loadSettings(): AppSettings {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored)
      return {
        audio: {
          ...DEFAULT_SETTINGS.audio,
          ...parsed.audio,
        },
        colors: {
          ...DEFAULT_SETTINGS.colors,
          ...parsed.colors,
        },
        wnp: {
          ...DEFAULT_SETTINGS.wnp,
          ...parsed.wnp,
        },
      }
    }
  }
  catch (error) {
    console.error('[Settings] Failed to load settings:', error)
  }
  return DEFAULT_SETTINGS
}

/**
 * Save settings to localStorage
 */
function saveSettings() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings.value))
  }
  catch (error) {
    console.error('[Settings] Failed to save settings:', error)
  }
}

/**
 * Apply CSS variables from color scheme
 */
function applyCSSVariables() {
  const root = document.documentElement
  const colors = settings.value.colors

  root.style.setProperty('--color-primary', colors.primary)
  root.style.setProperty('--color-primary-glow', colors.primaryGlow)
  root.style.setProperty('--color-secondary', colors.secondary)
  root.style.setProperty('--color-accent', colors.accent)
}

// Watch for settings changes and save
watch(settings, () => {
  saveSettings()
  applyCSSVariables()
}, { deep: true })

// Apply CSS variables on load
applyCSSVariables()

export function useSettings() {
  const audioSettings = computed(() => settings.value.audio)
  const colorScheme = computed(() => settings.value.colors)
  const wnpSettings = computed(() => settings.value.wnp)

  function updateAudioSettings(partial: Partial<typeof DEFAULT_SETTINGS.audio>) {
    settings.value.audio = {
      ...settings.value.audio,
      ...partial,
    }
  }

  function updatePrimaryColor(hex: string) {
    const scheme = generateColorScheme(hex)
    settings.value.colors = scheme
  }

  function updateColorScheme(scheme: Partial<typeof DEFAULT_SETTINGS.colors>) {
    settings.value.colors = {
      ...settings.value.colors,
      ...scheme,
    }
  }

  function updateWNPSettings(partial: Partial<typeof DEFAULT_SETTINGS.wnp>) {
    settings.value.wnp = {
      ...settings.value.wnp,
      ...partial,
    }
  }

  function resetToDefaults() {
    settings.value = JSON.parse(JSON.stringify(DEFAULT_SETTINGS))
  }

  return {
    audioSettings,
    colorScheme,
    wnpSettings,
    updateAudioSettings,
    updatePrimaryColor,
    updateColorScheme,
    updateWNPSettings,
    resetToDefaults,
  }
}
