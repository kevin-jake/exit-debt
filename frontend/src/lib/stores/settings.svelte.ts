import { browser } from "$app/environment";
import { writable } from "svelte/store";

export interface Settings {
  // General Settings
  currency: {
    code: string;
    symbol: string;
    position: "before" | "after";
    decimalPlaces: number;
    thousandsSeparator: string;
    decimalSeparator: string;
  };
  dateFormat: "MM/DD/YYYY" | "DD/MM/YYYY" | "YYYY-MM-DD";
  timezone: string;

  // Basic Notification Settings
  notifications: {
    email: {
      enabled: boolean;
      dueDateReminders: number[];
      paymentConfirmations: boolean;
      newDebtLists: boolean;
      summaries: "weekly" | "monthly" | "never";
      marketing: boolean;
    };
  };

  // Display Preferences
  display: {
    theme: "light" | "dark" | "system";
  };

  // Data Management
  dataManagement: {
    autoBackup: boolean;
    backupFrequency: "daily" | "weekly" | "monthly" | "never";
  };
}

function getDefaultSettings(): Settings {
  return {
    currency: {
      code: "USD",
      symbol: "$",
      position: "before",
      decimalPlaces: 2,
      thousandsSeparator: ",",
      decimalSeparator: ".",
    },
    dateFormat: "MM/DD/YYYY",
    timezone: browser
      ? Intl.DateTimeFormat().resolvedOptions().timeZone
      : "UTC",
    notifications: {
      email: {
        enabled: true,
        dueDateReminders: [1, 3, 7],
        paymentConfirmations: true,
        newDebtLists: true,
        summaries: "weekly",
        marketing: false,
      },
    },
    display: {
      theme: "system",
    },
    dataManagement: {
      autoBackup: true,
      backupFrequency: "weekly",
    },
  };
}

function createSettingsStore() {
  const { subscribe, set, update } = writable<Settings>(getDefaultSettings());
  const { subscribe: subscribeLoading, set: setLoading } = writable(false);
  const { subscribe: subscribeChanges, set: setChanges } = writable(false);

  return {
    subscribe,
    isLoading: { subscribe: subscribeLoading },
    hasUnsavedChanges: { subscribe: subscribeChanges },

    async loadSettings() {
      if (!browser) return;

      setLoading(true);
      try {
        const stored = localStorage.getItem("app-settings");
        if (stored) {
          const parsedSettings = JSON.parse(stored);
          set({ ...getDefaultSettings(), ...parsedSettings });
        }
        this.applySettings();
      } catch (error) {
        console.error("Failed to load settings:", error);
      } finally {
        setLoading(false);
      }
    },

    async saveSettings() {
      if (!browser) return;

      try {
        const currentSettings = get(this);
        localStorage.setItem("app-settings", JSON.stringify(currentSettings));
        setChanges(false);
        this.applySettings();
      } catch (error) {
        console.error("Failed to save settings:", error);
      }
    },

    applySettings() {
      const currentSettings = get(this);

      // Apply theme
      if (currentSettings.display?.theme) {
        // Integrate with existing theme system
        if (currentSettings.display.theme === "system") {
          const prefersDark = window.matchMedia(
            "(prefers-color-scheme: dark)"
          ).matches;
          document.documentElement.classList.toggle("dark", prefersDark);
        } else {
          document.documentElement.classList.toggle(
            "dark",
            currentSettings.display.theme === "dark"
          );
        }
      }

      // Apply currency formatting
      if (currentSettings.currency) {
        // Update currency display throughout app
        // This will be handled by components that need currency formatting
      }
    },

    resetToDefaults() {
      set(getDefaultSettings());
      setChanges(true);
    },

    updateSetting(path: string, value: any) {
      update((settings) => {
        const keys = path.split(".");
        let current: any = settings;

        for (let i = 0; i < keys.length - 1; i++) {
          if (!(keys[i] in current)) {
            current[keys[i]] = {};
          }
          current = current[keys[i]];
        }

        current[keys[keys.length - 1]] = value;
        setChanges(true);
        return settings;
      });
    },
  };
}

export const settingsStore = createSettingsStore();

// Helper function to get current store value
function get(store: any) {
  let value: any;
  store.subscribe((v: any) => (value = v))();
  return value;
}
