import apiClient from "./httpClient";

export const getBookingAddons = (email) =>
  apiClient.get("/booking-addons", { params: { email } }).then((r) => r.data);

export const addBookingAddonItem = (payload) =>
  apiClient.post("/booking-addons/items", payload).then((r) => r.data);

export const updateBookingAddonItem = (payload) =>
  apiClient.put("/booking-addons/items", payload).then((r) => r.data);

export const removeBookingAddonItem = (payload) =>
  apiClient
    .delete("/booking-addons/items", { data: payload })
    .then((r) => r.data);

export const checkoutBookingAddons = (payload) =>
  apiClient.post("/booking-addons/checkout", payload).then((r) => r.data);

export default {
  getBookingAddons,
  addBookingAddonItem,
  updateBookingAddonItem,
  removeBookingAddonItem,
  checkoutBookingAddons,
};
