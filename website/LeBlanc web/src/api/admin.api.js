import apiClient from "./httpClient";
import { ADMIN_LIST_TIMEOUT_MS } from "@/config/adminConfig";

const withQuery = (params = {}) => ({
  params,
});

const unwrapAdminData = (res) => res.data?.data;

export const getAdminUsers = () =>
  apiClient
    .get("/admin/users", { timeout: ADMIN_LIST_TIMEOUT_MS })
    .then(unwrapAdminData);

export const createAdminUser = (payload) =>
  apiClient.post("/admin/users", payload).then(unwrapAdminData);

export const updateAdminUser = (id, payload) =>
  apiClient.put(`/admin/users/${id}`, payload).then(unwrapAdminData);

export const deleteAdminUser = (id) =>
  apiClient.delete(`/admin/users/${id}`).then((res) => res.data);

export const getAdminBookings = () =>
  apiClient
    .get("/admin/bookings", { timeout: ADMIN_LIST_TIMEOUT_MS })
    .then(unwrapAdminData);

export const createAdminBooking = (payload) =>
  apiClient.post("/admin/bookings", payload).then(unwrapAdminData);

export const updateAdminBooking = (id, payload) =>
  apiClient.put(`/admin/bookings/${id}`, payload).then(unwrapAdminData);

export const deleteAdminBooking = (id) =>
  apiClient.delete(`/admin/bookings/${id}`).then((res) => res.data);

export const getAdminDrinks = () =>
  apiClient
    .get("/admin/drinks", {
      ...withQuery(),
      timeout: ADMIN_LIST_TIMEOUT_MS,
    })
    .then(unwrapAdminData);

export const getAdminDrinksWithFilters = (params = {}) =>
  apiClient
    .get("/admin/drinks", {
      ...withQuery(params),
      timeout: ADMIN_LIST_TIMEOUT_MS,
    })
    .then(unwrapAdminData);

export const createAdminDrink = (payload) =>
  apiClient.post("/admin/drinks", payload).then(unwrapAdminData);

export const updateAdminDrink = (id, payload) =>
  apiClient.put(`/admin/drinks/${id}`, payload).then(unwrapAdminData);

export const deleteAdminDrink = (id) =>
  apiClient.delete(`/admin/drinks/${id}`).then((res) => res.data);
