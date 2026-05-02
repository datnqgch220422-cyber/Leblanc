import { defineStore } from "pinia";

const THEME_STORAGE_KEY = "theme";

const normalizeTheme = (value) => (value === "night" ? "night" : "day");

const readStoredTheme = () => {
  if (typeof window === "undefined") return "day";

  try {
    return normalizeTheme(localStorage.getItem(THEME_STORAGE_KEY));
  } catch (err) {
    console.warn("Could not read stored theme", err);
    return "day";
  }
};

const applyThemeToDocument = (value) => {
  if (typeof document === "undefined") return;

  const themeValue = normalizeTheme(value);
  const body = document.body;
  body.classList.remove("theme-day", "theme-night");
  body.classList.add(themeValue === "night" ? "theme-night" : "theme-day");
  document.documentElement.setAttribute("data-theme", themeValue);

  try {
    localStorage.setItem(THEME_STORAGE_KEY, themeValue);
  } catch (err) {
    console.warn("Could not persist theme", err);
  }
};

export const useThemeStore = defineStore("theme", {
  state: () => ({
    theme: readStoredTheme(),
  }),
  getters: {
    isNight: (state) => state.theme === "night",
  },
  actions: {
    initTheme() {
      this.theme = normalizeTheme(this.theme);
      applyThemeToDocument(this.theme);
    },
    setTheme(value) {
      this.theme = normalizeTheme(value);
      applyThemeToDocument(this.theme);
    },
    toggleTheme() {
      this.setTheme(this.isNight ? "day" : "night");
    },
  },
});
