import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export const useThemeStore = create(
  persist(
    (set, get) => ({
      // State
      isDark: false,

      // Toggle theme
      toggle: () => {
        set((state) => {
          const newIsDark = !state.isDark
          if (newIsDark) {
            document.documentElement.classList.add('dark')
          } else {
            document.documentElement.classList.remove('dark')
          }
          return { isDark: newIsDark }
        })
      },

      // Set theme
      setDark: (isDark) => {
        set({ isDark })
        if (isDark) {
          document.documentElement.classList.add('dark')
        } else {
          document.documentElement.classList.remove('dark')
        }
      },

      // Initialize theme from system preference
      initFromSystem: () => {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        get().setDark(prefersDark)
      },
    }),
    {
      name: 'theme-storage', // localStorage key
    }
  )
)

// Initialize theme on app load
export const initializeTheme = () => {
  const theme = useThemeStore.getState()
  // Apply stored theme
  if (theme.isDark) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

