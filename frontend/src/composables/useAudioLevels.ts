import {
  onMounted,
  onUnmounted,
  ref,
  watch,
} from 'vue'
import { EventsEmit } from '../../wailsjs/runtime/runtime'
import { useSettings } from './useSettings'

// Check if Wails runtime is available
const isWailsAvailable = () => typeof window !== 'undefined' && 'runtime' in window

export function useAudioLevels(bandCount = 64) {
  const { audioSettings } = useSettings()
  const levels = ref<number[]>(new Array(bandCount).fill(0))
  const isActive = ref(false)

  let unsubscribe: (() => void) | null = null
  let animationId: number | null = null

  // Send audio settings to backend when changed
  function updateBackendSettings() {
    if (!isWailsAvailable()) return

    EventsEmit('audio:config', {
      fftSize: audioSettings.value.fftSize,
      freqMin: audioSettings.value.freqMin,
      freqMax: audioSettings.value.freqMax,
    })
  }

  onMounted(() => {
    if (!isWailsAvailable()) {
      // Mock data for development - simulate audio levels
      const simulateLevels = () => {
        const newLevels = new Array(bandCount).fill(0).map(() =>
          Math.random() * 0.5 + Math.sin(Date.now() / 500) * 0.3 + 0.2,
        )
        levels.value = newLevels
        animationId = requestAnimationFrame(simulateLevels)
      }
      simulateLevels()
      isActive.value = true
      return
    }

    // Subscribe to audio level updates from backend
    unsubscribe = window.runtime.EventsOn('audio:levels', (...args: unknown[]) => {
      const data = args[0] as number[] | undefined
      if (data && Array.isArray(data)) {
        levels.value = data
        isActive.value = true
      }
    })

    // Send initial settings
    updateBackendSettings()
  })

  onUnmounted(() => {
    if (unsubscribe) unsubscribe()
    if (animationId) cancelAnimationFrame(animationId)
  })

  // Watch for settings changes and update backend
  watch(audioSettings, () => {
    updateBackendSettings()
  }, { deep: true })

  return {
    levels,
    isActive,
  }
}
