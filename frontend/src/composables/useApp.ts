// Check if Wails runtime is available
const isWailsAvailable = () => typeof window !== 'undefined' && 'go' in window

export function useApp() {
  const quit = async () => {
    console.log('[useApp] quit called')
    if (isWailsAvailable()) {
      try {
        await window.go.app.App.Quit()
        console.log('[useApp] Quit success')
      }
      catch (err) {
        console.error('[useApp] Quit error:', err)
      }
    }
  }

  return {
    quit,
  }
}
