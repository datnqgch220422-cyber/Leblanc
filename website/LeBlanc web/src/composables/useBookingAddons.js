import { ref } from "vue";
import { getSessionUser } from "./useSessionAuth";
import {
  getBookingAddons,
  addBookingAddonItem,
  updateBookingAddonItem,
  removeBookingAddonItem,
  checkoutBookingAddons,
} from "@/api/bookingAddons.api";

export const useBookingAddons = () => {
  const bookingAddons = ref({ items: [] });
  const loading = ref(false);
  const error = ref("");

  const load = async (email) => {
    loading.value = true;
    error.value = "";
    try {
      const userEmail = email || getSessionUser()?.email || "";
      if (!userEmail) {
        bookingAddons.value = { items: [] };
        return;
      }
      const res = await getBookingAddons(userEmail);
      bookingAddons.value = res || { items: [] };
    } catch (err) {
      error.value = err?.message || "Could not load booking add-ons";
    } finally {
      loading.value = false;
    }
  };

  const add = async (drinkId, qty = 1, email) => {
    loading.value = true;
    try {
      const userEmail = email || getSessionUser()?.email || "";
      if (!userEmail) throw new Error("email required");
      await addBookingAddonItem({ email: userEmail, drinkId, qty });
      await load(userEmail);
    } finally {
      loading.value = false;
    }
  };

  const update = async (drinkId, qty, email) => {
    loading.value = true;
    try {
      const userEmail = email || getSessionUser()?.email || "";
      if (!userEmail) throw new Error("email required");
      await updateBookingAddonItem({ email: userEmail, drinkId, qty });
      await load(userEmail);
    } finally {
      loading.value = false;
    }
  };

  const remove = async (drinkId, email) => {
    loading.value = true;
    try {
      const userEmail = email || getSessionUser()?.email || "";
      if (!userEmail) throw new Error("email required");
      await removeBookingAddonItem({ email: userEmail, drinkId });
      await load(userEmail);
    } finally {
      loading.value = false;
    }
  };

  const doCheckout = async (payload) => {
    loading.value = true;
    try {
      const user = getSessionUser();
      const email = payload?.email || user?.email;
      if (!email) throw new Error("email required");
      const body = { ...payload, email };
      const res = await checkoutBookingAddons(body);
      await load(email);
      return res;
    } finally {
      loading.value = false;
    }
  };

  return {
    bookingAddons,
    loading,
    error,
    load,
    add,
    update,
    remove,
    doCheckout,
  };
};
