export const USER_STORAGE_KEY = "leblancUser";
export const USER_UPDATED_EVENT = "leblanc-user-updated";

const isBrowser = () => typeof window !== "undefined";

const isLegacyUserShape = (value) =>
  value &&
  typeof value === "object" &&
  !Array.isArray(value) &&
  !("user" in value) &&
  ("_id" in value || "email" in value || "name" in value);

const normalizeStoredSession = (value) => {
  if (!value) return null;

  if (value.user && typeof value.user === "object") {
    return {
      user: value.user,
      token: typeof value.token === "string" ? value.token : "",
    };
  }

  if (isLegacyUserShape(value)) {
    return {
      user: value,
      token: "",
    };
  }

  return null;
};

export const getStoredSession = () => {
  if (!isBrowser()) return null;

  try {
    const raw = localStorage.getItem(USER_STORAGE_KEY);
    if (!raw) return null;
    return normalizeStoredSession(JSON.parse(raw));
  } catch (err) {
    console.warn("Could not parse stored session", err);
    return null;
  }
};

export const getStoredUser = () => getStoredSession()?.user || null;

export const getStoredToken = () => getStoredSession()?.token || "";

export const setStoredSession = (session) => {
  if (!isBrowser()) return;

  const normalized = normalizeStoredSession(session);
  if (!normalized) {
    removeStoredUser();
    return;
  }

  try {
    localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(normalized));
  } catch (err) {
    console.warn("Could not persist session", err);
  }
};

export const setStoredUser = (user) => {
  const existingToken = getStoredToken();
  setStoredSession({ user, token: existingToken });
};

export const removeStoredUser = () => {
  if (!isBrowser()) return;

  try {
    localStorage.removeItem(USER_STORAGE_KEY);
  } catch (err) {
    console.warn("Could not clear stored session", err);
  }
};

export const emitUserUpdated = (session) => {
  if (!isBrowser()) return;
  window.dispatchEvent(new CustomEvent(USER_UPDATED_EVENT, { detail: session }));
};

export const subscribeUserUpdated = (handler) => {
  if (!isBrowser() || typeof handler !== "function") {
    return () => {};
  }

  window.addEventListener(USER_UPDATED_EVENT, handler);
  return () => window.removeEventListener(USER_UPDATED_EVENT, handler);
};
