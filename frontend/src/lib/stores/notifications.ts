import { writable } from "svelte/store";

export interface Notification {
  id: string;
  type: "success" | "error" | "warning" | "info";
  title: string;
  message: string;
  duration?: number; // in milliseconds, undefined means manual dismiss
  timestamp: Date;
}

export interface NotificationsState {
  notifications: Notification[];
}

function createNotificationsStore() {
  const { subscribe, set, update } = writable<NotificationsState>({
    notifications: [],
  });

  return {
    subscribe,

    add(notification: Omit<Notification, "id" | "timestamp">) {
      const id = Math.random().toString(36).substr(2, 9);
      const newNotification: Notification = {
        ...notification,
        id,
        timestamp: new Date(),
      };

      update((state) => ({
        ...state,
        notifications: [...state.notifications, newNotification],
      }));

      // Auto-remove notification after duration
      if (notification.duration !== undefined) {
        setTimeout(() => {
          this.remove(id);
        }, notification.duration);
      }

      return id;
    },

    success(title: string, message: string, duration = 5000) {
      return this.add({
        type: "success",
        title,
        message,
        duration,
      });
    },

    error(title: string, message: string, duration = 8000) {
      return this.add({
        type: "error",
        title,
        message,
        duration,
      });
    },

    warning(title: string, message: string, duration = 6000) {
      return this.add({
        type: "warning",
        title,
        message,
        duration,
      });
    },

    info(title: string, message: string, duration = 4000) {
      return this.add({
        type: "info",
        title,
        message,
        duration,
      });
    },

    remove(id: string) {
      update((state) => ({
        ...state,
        notifications: state.notifications.filter((n) => n.id !== id),
      }));
    },

    clear() {
      set({ notifications: [] });
    },

    clearByType(type: Notification["type"]) {
      update((state) => ({
        ...state,
        notifications: state.notifications.filter((n) => n.type !== type),
      }));
    },
  };
}

export const notificationsStore = createNotificationsStore();
