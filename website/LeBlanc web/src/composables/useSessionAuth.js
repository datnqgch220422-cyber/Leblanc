import {
  emitUserUpdated,
  getStoredSession,
  getStoredToken,
  getStoredUser,
  removeStoredUser,
  setStoredSession,
} from "@/services/session.service";

export const ADMIN_HOME_PATH = "/admin/users";

export const isAdminUser = (user) => (user?.role || "").toLowerCase() === "admin";

export const getSession = () => getStoredSession();

export const getSessionUser = () => getStoredUser();

export const getSessionToken = () => getStoredToken();

export const hasSessionToken = () => Boolean(getStoredToken());

export const persistSession = ({ user, token = "" }) => {
  setStoredSession({ user, token });
  emitUserUpdated({ user, token });
};

export const persistSessionUser = (user, token = "") => {
  persistSession({ user, token });
};

export const clearSessionUser = () => {
  removeStoredUser();
  emitUserUpdated(null);
};
