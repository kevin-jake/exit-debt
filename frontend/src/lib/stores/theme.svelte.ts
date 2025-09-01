import { browser } from "$app/environment";
import { writable } from "svelte/store";
import { getTheme, setTheme } from "$lib/utils";

function createThemeStore() {
  const { subscribe, set, update } = writable<"dark" | "light">(getTheme());

  return {
    subscribe,
    toggle: () => {
      update((currentTheme) => {
        const newTheme = currentTheme === "dark" ? "light" : "dark";
        setTheme(newTheme);
        return newTheme;
      });
    },
    set: (newTheme: "dark" | "light") => {
      setTheme(newTheme);
      set(newTheme);
    },
  };
}

export const themeStore = createThemeStore();
