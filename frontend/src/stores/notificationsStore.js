import { create } from 'zustand'

let notificationId = 0

export const useNotificationsStore = create((set) => ({
  // State
  notifications: [],

  // Add a notification
  addNotification: (notification) => {
    const id = ++notificationId
    const newNotification = {
      id,
      type: notification.type || 'info', // info, success, warning, error
      message: notification.message,
      duration: notification.duration || 5000,
      timestamp: new Date().toISOString(),
    }

    set((state) => ({
      notifications: [...state.notifications, newNotification],
    }))

    // Auto remove after duration
    if (newNotification.duration > 0) {
      setTimeout(() => {
        set((state) => ({
          notifications: state.notifications.filter((n) => n.id !== id),
        }))
      }, newNotification.duration)
    }

    return id
  },

  // Remove a notification
  removeNotification: (id) => {
    set((state) => ({
      notifications: state.notifications.filter((n) => n.id !== id),
    }))
  },

  // Clear all notifications
  clearAll: () => {
    set({ notifications: [] })
  },

  // Helper methods for different notification types
  success: (message, duration) => {
    return useNotificationsStore.getState().addNotification({
      type: 'success',
      message,
      duration,
    })
  },

  error: (message, duration) => {
    return useNotificationsStore.getState().addNotification({
      type: 'error',
      message,
      duration,
    })
  },

  warning: (message, duration) => {
    return useNotificationsStore.getState().addNotification({
      type: 'warning',
      message,
      duration,
    })
  },

  info: (message, duration) => {
    return useNotificationsStore.getState().addNotification({
      type: 'info',
      message,
      duration,
    })
  },
}))

