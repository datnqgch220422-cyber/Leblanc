import { storeToRefs } from "pinia";
import { useThemeStore } from "@/stores/theme.store";

let initialized = false;

const ensureThemeInitialized = (store) => {
  if (initialized) return;
  store.initTheme();
  initialized = true;
};

export const useThemeState = () => {
  const themeStore = useThemeStore();
  ensureThemeInitialized(themeStore);
  const { theme, isNight } = storeToRefs(themeStore);

  return {
    theme,
    isNight,
    toggleTheme: () => themeStore.toggleTheme(),
    setTheme: (value) => themeStore.setTheme(value),
  };
};

export const provideThemeState = () => useThemeState();
