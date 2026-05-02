import apiClient from "./httpClient";

const unwrapAdminData = (res) => res.data?.data;

export const getAdminUsers = () =>
  apiClient.get("/admin/users").then(unwrapAdminData);

export const createAdminUser = (payload) =>
  apiClient.post("/admin/users", payload).then(unwrapAdminData);

export const updateAdminUser = (id, payload) =>
  apiClient.put(`/admin/users/${id}`, payload).then(unwrapAdminData);

export const deleteAdminUser = (id) =>
  apiClient.delete(`/admin/users/${id}`).then((res) => res.data);

export const getAdminBookings = () =>
  apiClient.get("/admin/bookings").then(unwrapAdminData);

export const createAdminBooking = (payload) =>
  apiClient.post("/admin/bookings", payload).then(unwrapAdminData);

export const updateAdminBooking = (id, payload) =>
  apiClient.put(`/admin/bookings/${id}`, payload).then(unwrapAdminData);

export const deleteAdminBooking = (id) =>
  apiClient.delete(`/admin/bookings/${id}`).then((res) => res.data);

export const getAdminDrinks = () =>
  apiClient.get("/admin/drinks").then(unwrapAdminData);

export const createAdminDrink = (payload) =>
  apiClient.post("/admin/drinks", payload).then(unwrapAdminData);

export const updateAdminDrink = (id, payload) =>
  apiClient.put(`/admin/drinks/${id}`, payload).then(unwrapAdminData);

export const deleteAdminDrink = (id) =>
  apiClient.delete(`/admin/drinks/${id}`).then((res) => res.data);
