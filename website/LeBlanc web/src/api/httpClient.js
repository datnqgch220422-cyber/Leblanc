import axios from "axios";
import { getStoredToken } from "@/services/session.service";

export const USE_GRAPHQL = import.meta.env.VITE_USE_GRAPHQL === "true";

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE,
});

apiClient.interceptors.request.use((config) => {
  const token = getStoredToken();
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default apiClient;
