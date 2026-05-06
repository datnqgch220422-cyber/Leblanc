import {
  computed,
  onBeforeUnmount,
  onMounted,
  nextTick,
  ref,
  watch,
} from "vue";
import { createBooking, getDrinks } from "@/api";
import { getSessionUser } from "@/composables/useSessionAuth";
import { subscribeUserUpdated } from "@/services/session.service";
import { useBookingAddons } from "@/composables/useBookingAddons";

export const useBookingFlow = (isNightRef) => {
  const form = ref({
    name: "",
    phone: "",
    email: "",
    time: "",
    guests: 1,
  });
  const formDate = ref("");
  const formClock = ref("");

  const caffeinePref = ref("");
  const tempPref = ref("");
  const sweetness = ref(5);
  const nightType = ref("");
  const nightBase = ref("");

  const drinks = ref([]);
  const reco = ref([]);
  const selection = ref({});

  const bookingLoading = ref(false);
  const bookingOk = ref(false);
  const bookingError = ref("");
  const hydratingSelectionFromAddons = ref(false);

  const recoLoading = ref(false);
  const recoError = ref("");

  const canSubmit = computed(
    () =>
      form.value.name &&
      form.value.phone &&
      form.value.email &&
      formDate.value &&
      formClock.value &&
      form.value.time,
  );

  const isNightMode = () => Boolean(isNightRef?.value);

  watch([formDate, formClock], ([date, clock]) => {
    if (date && clock) {
      const local = new Date(`${date}T${clock}`);
      if (!Number.isNaN(local.getTime())) {
        form.value.time = local.toISOString();
        return;
      }
    }
    form.value.time = "";
  });

  const fetchDrinks = async () => {
    try {
      drinks.value = await getDrinks();
    } catch (err) {
      console.warn("Could not load drinks", err);
    }
  };

  const resolveDrink = (id) => drinks.value.find((d) => d._id === id);

  const fetchReco = async () => {
    recoLoading.value = true;
    recoError.value = "";
    try {
      const result = [];

      let mapped = (result || []).map((item) => {
        const drink = resolveDrink(item.drinkId) || {};
        return {
          ...drink,
          drinkId: item.drinkId || drink._id,
          score: item.score,
        };
      });

      if (isNightMode()) {
        const typeTag = nightType.value?.toLowerCase();
        const baseTag = nightBase.value?.toLowerCase();
        mapped = mapped.filter((item) => {
          const tags = (item.tags || item.Tags || []).map((t) =>
            (t || "").toLowerCase(),
          );
          const okType = !typeTag || tags.includes(typeTag);
          const okBase = !baseTag || tags.includes(baseTag);
          return okType && okBase;
        });
      }

      reco.value = mapped;
    } catch (err) {
      recoError.value =
        err?.message || "Unable to get recommendations right now.";
    } finally {
      recoLoading.value = false;
    }
  };

  const addDrink = (drink) => {
    if (!drink?.drinkId && !drink?._id) return;
    const id = drink.drinkId || drink._id;
    const current = selection.value[id]?.qty || 0;
    selection.value = {
      ...selection.value,
      [id]: { drink, qty: current + 1 },
    };
  };

  const updateQty = (id, delta) => {
    const current = selection.value[id]?.qty || 0;
    const next = Math.max(0, current + delta);
    if (next === 0) {
      const { [id]: _, ...rest } = selection.value;
      selection.value = rest;
      return;
    }
    selection.value = {
      ...selection.value,
      [id]: { drink: selection.value[id]?.drink, qty: next },
    };
  };

  // booking addons sync: initialize and propagate to server-side booking addons when user is signed in
  const {
    bookingAddons,
    load: loadAddons,
    add: addAddon,
    update: updateAddon,
    remove: removeAddon,
  } = useBookingAddons();

  const syncSelectionFromAddons = () => {
    try {
      const items = Array.isArray(bookingAddons.value?.items)
        ? bookingAddons.value.items
        : [];
      const map = {};
      for (const it of items) {
        const id = it.drinkId;
        const qty = Number(it.qty) || 0;
        if (!id || qty <= 0) continue;
        map[id] = { drink: resolveDrink(id) || {}, qty };
      }
      hydratingSelectionFromAddons.value = true;
      selection.value = map;
      void nextTick(() => {
        hydratingSelectionFromAddons.value = false;
      });
    } catch (err) {
      hydratingSelectionFromAddons.value = false;
      // ignore
    }
  };

  // watch selection and push diffs to bookingAddons for logged-in user
  watch(
    selection,
    async (val, oldVal) => {
      try {
        if (hydratingSelectionFromAddons.value) return;

        const user = getSessionUser();
        const email = user?.email;
        if (!email) return;

        const newIds = Object.keys(val || {});
        const oldIds = Object.keys(oldVal || {});

        for (const id of newIds) {
          const newQty = Number(val[id]?.qty || 0);
          const oldQty = Number(oldVal?.[id]?.qty || 0);
          if (newQty <= 0) continue;
          if (!oldIds.includes(id)) {
            await addAddon(id, newQty, email);
          } else if (newQty !== oldQty) {
            await updateAddon(id, newQty, email);
          }
        }

        for (const id of oldIds) {
          if (!newIds.includes(id)) {
            await removeAddon(id, email);
          }
        }
      } catch (err) {
        // non-fatal
      }
    },
    { deep: true },
  );

  const selectedItems = computed(() =>
    Object.entries(selection.value).map(([drinkId, entry]) => ({
      drinkId,
      drink: entry.drink || resolveDrink(drinkId) || {},
      qty: entry.qty,
    })),
  );

  const totalItems = computed(() =>
    selectedItems.value.reduce((sum, item) => sum + (item.qty || 0), 0),
  );

  const book = async () => {
    if (!canSubmit.value || bookingLoading.value) return;
    bookingLoading.value = true;
    bookingError.value = "";
    bookingOk.value = false;
    try {
      const items = selectedItems.value.map((item) => ({
        drinkId: item.drinkId,
        qty: item.qty,
        options: {},
      }));
      const payload = {
        ...form.value,
        items,
        channel: "web",
        paymentMethod: "vnpay",
      };
      const res = await createBooking(payload);
      const payUrl = String(res?.payUrl || "").trim();
      if (payUrl) {
        bookingOk.value = true;
        selection.value = {};
        window.location.href = payUrl;
        return;
      }

      bookingOk.value = Boolean(res?.ok || res?._id);
      if (bookingOk.value) {
        selection.value = {};
      }
    } catch (err) {
      bookingError.value =
        err?.message || "Unable to place booking at this time.";
    } finally {
      bookingLoading.value = false;
    }
  };

  const userPrefilled = ref(false);
  const applyUser = (user) => {
    if (!user) return;
    if (!userPrefilled.value || !form.value.name)
      form.value.name = user.name || form.value.name;
    if (!userPrefilled.value || !form.value.email)
      form.value.email = user.email || form.value.email;
    userPrefilled.value = true;
  };

  let unsubscribeUserUpdated = () => {};

  onMounted(() => {
    fetchDrinks();
    fetchReco();
    applyUser(getSessionUser());

    // load booking addons for current user and sync selection
    try {
      loadAddons();
      watch(bookingAddons, syncSelectionFromAddons);
    } catch (err) {
      // ignore
    }

    unsubscribeUserUpdated = subscribeUserUpdated((event) => {
      applyUser(event?.detail);
    });
  });

  onBeforeUnmount(() => {
    unsubscribeUserUpdated();
  });

  return {
    form,
    formDate,
    formClock,
    caffeinePref,
    tempPref,
    sweetness,
    nightType,
    nightBase,
    drinks,
    reco,
    selection,
    bookingLoading,
    bookingOk,
    bookingError,
    recoLoading,
    recoError,
    canSubmit,
    fetchReco,
    addDrink,
    updateQty,
    selectedItems,
    totalItems,
    book,
  };
};
