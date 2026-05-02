import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { createBooking, getDrinks, recoFromFeatures } from "@/api";
import { isBookingEmailReady, sendBookingEmail } from "@/email";
import { getSessionUser } from "@/composables/useSessionAuth";
import { subscribeUserUpdated } from "@/services/session.service";

export const useBookingFlow = (isNightRef) => {
  const form = ref({
    name: "",
    phone: "",
    email: "",
    time: "",
    guests: 2,
  });
  const formDate = ref("");
  const formClock = ref("");

  const mood = ref("happy");
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
  const bookingEmailSent = ref(false);
  const bookingEmailError = ref("");

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
  const bookingEmailReady = computed(() => isBookingEmailReady());

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

  const moodToEmotionFit = (val) => {
    switch (val) {
      case "calm":
        return {
          calm: 0.9,
          happy: 0.4,
          stressed: 0.2,
          sad: 0.3,
          adventurous: 0.3,
        };
      case "stressed":
        return {
          calm: 0.2,
          happy: 0.3,
          stressed: 0.9,
          sad: 0.2,
          adventurous: 0.3,
        };
      case "sad":
        return {
          calm: 0.3,
          happy: 0.2,
          stressed: 0.2,
          sad: 0.9,
          adventurous: 0.3,
        };
      case "adventurous":
        return {
          calm: 0.3,
          happy: 0.6,
          stressed: 0.3,
          sad: 0.2,
          adventurous: 0.9,
        };
      default:
        return {
          calm: 0.3,
          happy: 0.9,
          stressed: 0.2,
          sad: 0.2,
          adventurous: 0.4,
        };
    }
  };

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
      const result = await recoFromFeatures({
        emotionFit: moodToEmotionFit(mood.value),
        caffeine: caffeinePref.value || undefined,
        temp: tempPref.value || undefined,
        sweetness: sweetness.value,
      });

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
      recoError.value = err?.message || "Không thể gợi ý lúc này.";
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
    bookingEmailSent.value = false;
    bookingEmailError.value = "";
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
      };
      const res = await createBooking(payload);
      bookingOk.value = Boolean(res?.ok || res?._id);
      if (bookingOk.value) {
        if (form.value.email && bookingEmailReady.value) {
          const emailItems = selectedItems.value.map((item) => ({
            drinkId: item.drinkId,
            name: item.drink?.name || "Drink",
            qty: item.qty,
          }));
          const bookingForEmail = {
            ...payload,
            items: emailItems,
            bookingId: res?.id || res?._id || "",
          };
          sendBookingEmail(bookingForEmail)
            .then(() => {
              bookingEmailSent.value = true;
            })
            .catch((err) => {
              bookingEmailError.value =
                err?.message || "Gửi email xác nhận thất bại.";
              console.warn("Booking email failed", err);
            });
        }
        selection.value = {};
      }
    } catch (err) {
      bookingError.value = err?.message || "Không thể đặt lúc này.";
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
    mood,
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
    bookingEmailSent,
    bookingEmailError,
    recoLoading,
    recoError,
    canSubmit,
    bookingEmailReady,
    fetchReco,
    addDrink,
    updateQty,
    selectedItems,
    totalItems,
    book,
  };
};
