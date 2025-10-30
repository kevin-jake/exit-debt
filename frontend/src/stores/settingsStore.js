import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export const useSettingsStore = create(
  persist(
    (set) => ({
      // State
      theme: 'light', // 'light' | 'dark' | 'system'
      currency: 'USD',
      language: 'en',
      notifications: {
        email: true,
        push: true,
        reminders: true,
      },

      // Set theme
      setTheme: (theme) => {
        set({ theme })
        // Apply theme to document
        if (theme === 'dark') {
          document.documentElement.classList.add('dark')
        } else if (theme === 'light') {
          document.documentElement.classList.remove('dark')
        } else {
          // System preference
          const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
          if (prefersDark) {
            document.documentElement.classList.add('dark')
          } else {
            document.documentElement.classList.remove('dark')
          }
        }
      },

      // Set currency
      setCurrency: (currency) => {
        set({ currency })
      },

      // Set language
      setLanguage: (language) => {
        set({ language })
      },

      // Update notification settings
      setNotificationSettings: (settings) => {
        set((state) => ({
          notifications: {
            ...state.notifications,
            ...settings,
          },
        }))
      },

      // Reset to defaults
      resetSettings: () => {
        set({
          theme: 'light',
          currency: 'USD',
          language: 'en',
          notifications: {
            email: true,
            push: true,
            reminders: true,
          },
        })
      },
    }),
    {
      name: 'settings-storage', // localStorage key
    }
  )
)

// Initialize theme on app load
export const initializeTheme = () => {
  const settings = useSettingsStore.getState()
  settings.setTheme(settings.theme)
}

