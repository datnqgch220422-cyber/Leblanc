<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import {
  createAdminBooking,
  createAdminDrink,
  createAdminUser,
  deleteAdminBooking,
  deleteAdminDrink,
  deleteAdminUser,
  getAdminBookings,
  getAdminDrinksWithFilters,
  getAdminUsers,
  updateAdminBooking,
  updateAdminDrink,
  updateAdminUser,
} from "@/api";
import { useThemeState } from "@/composables/useThemeState";

const { theme } = useThemeState();

const adminCache = {
  users: null,
  bookings: null,
  drinks: new Map(),
};

const activeSection = ref("users");
const users = ref([]);
const bookings = ref([]);
const drinks = ref([]);
const initialLoading = ref(true);

const sectionLoaded = reactive({
  users: false,
  bookings: false,
  drinks: false,
});

const loading = reactive({
  users: false,
  bookings: false,
  drinks: false,
});

const saving = reactive({
  user: false,
  booking: false,
  drink: false,
});

// Controls whether the left "form" column is visible for each section.
const showUserForm = ref(false);
const showBookingForm = ref(false);
const showDrinkForm = ref(false);

const openUserForm = () => {
  resetUserForm();
  showUserForm.value = true;
  activeSection.value = "users";
};
const closeUserForm = () => {
  resetUserForm();
  showUserForm.value = false;
};

const openBookingForm = () => {
  resetBookingForm();
  showBookingForm.value = true;
  activeSection.value = "bookings";
};
const closeBookingForm = () => {
  resetBookingForm();
  showBookingForm.value = false;
};

const openDrinkForm = () => {
  resetDrinkForm();
  showDrinkForm.value = true;
  activeSection.value = "products";
};
const closeDrinkForm = () => {
  resetDrinkForm();
  showDrinkForm.value = false;
};

const sectionError = reactive({
  users: "",
  bookings: "",
  drinks: "",
});

const feedback = reactive({
  type: "",
  text: "",
});

const createEmptyUserForm = () => ({
  _id: "",
  name: "",
  email: "",
  password: "",
  role: "user",
  verified: true,
});

const createEmptyBookingItem = () => ({
  drinkId: "",
  qty: 1,
  optionsText: "{}",
});

const createEmptyBookingForm = () => ({
  _id: "",
  name: "",
  email: "",
  phone: "",
  time: "",
  guests: 2,
  channel: "admin",
  items: [createEmptyBookingItem()],
  status: "pending",
});

const createEmptyDrinkForm = () => ({
  _id: "",
  name: "",
  price: 0,
  stock: 0,
  available: false,
  tagsText: "",
  caffeine: "none",
  temp: "iced",
  sweetness: 5,
  colorTone: "neutral",
  calm: 0.5,
  happy: 0.5,
  stressed: 0.5,
  sad: 0.5,
  adventurous: 0.5,
  image: "",
  desc: "",
});

const userForm = reactive(createEmptyUserForm());
const bookingForm = reactive(createEmptyBookingForm());
const drinkForm = reactive(createEmptyDrinkForm());
const drinkListFilters = reactive({
  query: "",
  available: "all",
  sortBy: "name",
  order: "asc",
});

const isRefreshing = computed(
  () => loading.users || loading.bookings || loading.drinks,
);

const resetUserForm = () => Object.assign(userForm, createEmptyUserForm());
const resetDrinkForm = () => Object.assign(drinkForm, createEmptyDrinkForm());
const resetBookingForm = () => {
  const next = createEmptyBookingForm();
  bookingForm._id = next._id;
  bookingForm.name = next.name;
  bookingForm.email = next.email;
  bookingForm.phone = next.phone;
  bookingForm.time = next.time;
  bookingForm.guests = next.guests;
  bookingForm.channel = next.channel;
  bookingForm.status = next.status;
  bookingForm.items = next.items;
};

const formatBookingStatus = (status) => {
  switch ((status || "pending").toLowerCase()) {
    case "confirmed":
      return "Confirmed";
    case "completed":
      return "Completed";
    case "cancelled":
      return "Cancelled";
    default:
      return "Pending";
  }
};

const statusBadgeClass = (status) => ({
  ok: ["confirmed", "completed"].includes((status || "pending").toLowerCase()),
  danger: (status || "pending").toLowerCase() === "cancelled",
});

const setFeedback = (type, text) => {
  feedback.type = type;
  feedback.text = text;
};

const clearFeedback = () => {
  feedback.type = "";
  feedback.text = "";
};

const getErrorMessage = (err, fallback) =>
  err?.response?.data?.error || err?.message || fallback;

const pad = (value) => String(value).padStart(2, "0");

const toLocalDateTime = (iso) => {
  if (!iso) return "";
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return "";
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
    date.getDate(),
  )}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
};

const toIsoDateTime = (value) => {
  if (!value) return "";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "" : date.toISOString();
};

const formatDate = (value) => {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
};

const formatCurrency = (value) =>
  `${(Number(value) || 0).toLocaleString("vi-VN")} VND`;

const resolveDrinkName = (id) => {
  const drink = drinks.value.find((item) => item._id === id);
  return drink?.name || id || "Drink";
};

const summarizeBookingItems = (items = []) => {
  if (!items.length) return "No pre-order items";
  return items
    .map((item) => `${resolveDrinkName(item.drinkId)} x${item.qty || 1}`)
    .join(", ");
};

const getDrinksCacheKey = () =>
  JSON.stringify({
    q: drinkListFilters.query.trim(),
    sort: drinkListFilters.sortBy,
    order: drinkListFilters.order,
    available: drinkListFilters.available,
  });

const ensureUsersLoaded = async () => {
  if (adminCache.users !== null) {
    users.value = adminCache.users;
    sectionLoaded.users = true;
    return;
  }

  await loadUsers();
};

const ensureBookingsLoaded = async () => {
  if (adminCache.bookings !== null) {
    bookings.value = adminCache.bookings;
    sectionLoaded.bookings = true;
    return;
  }

  await loadBookings();
};

const ensureDrinksLoaded = async () => {
  const cacheKey = getDrinksCacheKey();
  if (adminCache.drinks.has(cacheKey)) {
    drinks.value = adminCache.drinks.get(cacheKey);
    sectionLoaded.drinks = true;
    return;
  }

  await loadDrinks();
};

const ensureActiveSectionLoaded = async (section = activeSection.value) => {
  if (section === "users") {
    await ensureUsersLoaded();
    return;
  }

  if (section === "bookings") {
    await ensureBookingsLoaded();
    return;
  }

  await ensureDrinksLoaded();
};

const syncUsersCache = (nextUsers) => {
  users.value = nextUsers;
  adminCache.users = nextUsers;
  sectionLoaded.users = true;
};

const syncBookingsCache = (nextBookings) => {
  bookings.value = nextBookings;
  adminCache.bookings = nextBookings;
  sectionLoaded.bookings = true;
};

const syncDrinksCache = (nextDrinks, cacheKey = getDrinksCacheKey()) => {
  drinks.value = nextDrinks;
  adminCache.drinks.set(cacheKey, nextDrinks);
  sectionLoaded.drinks = true;
};

const loadUsers = async () => {
  loading.users = true;
  sectionError.users = "";
  try {
    const nextUsers = await getAdminUsers();
    syncUsersCache(nextUsers || []);
  } catch (err) {
    sectionError.users = getErrorMessage(err, "Could not load users.");
    sectionLoaded.users = false;
  } finally {
    loading.users = false;
  }
};

const loadBookings = async () => {
  loading.bookings = true;
  sectionError.bookings = "";
  try {
    const nextBookings = await getAdminBookings();
    syncBookingsCache(nextBookings || []);
  } catch (err) {
    sectionError.bookings = getErrorMessage(err, "Could not load bookings.");
    sectionLoaded.bookings = false;
  } finally {
    loading.bookings = false;
  }
};

const loadDrinks = async () => {
  loading.drinks = true;
  sectionError.drinks = "";
  try {
    const params = {
      q: drinkListFilters.query.trim(),
      sort: drinkListFilters.sortBy,
      order: drinkListFilters.order,
    };

    if (drinkListFilters.available === "true") {
      params.available = true;
    } else if (drinkListFilters.available === "false") {
      params.available = false;
    }

    const cacheKey = JSON.stringify(params);
    if (adminCache.drinks.has(cacheKey)) {
      drinks.value = adminCache.drinks.get(cacheKey);
      sectionLoaded.drinks = true;
      return;
    }

    const nextDrinks = await getAdminDrinksWithFilters(params);
    syncDrinksCache(nextDrinks || [], cacheKey);
  } catch (err) {
    sectionError.drinks = getErrorMessage(err, "Could not load products.");
    sectionLoaded.drinks = false;
  } finally {
    loading.drinks = false;
  }
};

const resetDrinkFilters = async () => {
  drinkListFilters.query = "";
  drinkListFilters.available = "all";
  drinkListFilters.sortBy = "name";
  drinkListFilters.order = "asc";
  await loadDrinks();
};

const refreshAll = async () => {
  clearFeedback();
  await Promise.allSettled([loadUsers(), loadBookings(), loadDrinks()]);
};

watch(activeSection, async (section) => {
  if (!sectionLoaded[section]) {
    await ensureActiveSectionLoaded(section);
  }
});

const editUser = (user) => {
  clearFeedback();
  activeSection.value = "users";
  Object.assign(userForm, {
    _id: user._id,
    name: user.name || "",
    email: user.email || "",
    password: "",
    role: user.role || "user",
    verified: Boolean(user.verified),
  });
  showUserForm.value = true;
};

const submitUser = async () => {
  clearFeedback();
  sectionError.users = "";

  const payload = {
    name: userForm.name.trim(),
    email: userForm.email.trim(),
    password: userForm.password.trim(),
    role: userForm.role,
    verified: Boolean(userForm.verified),
  };

  if (!payload.name || !payload.email || (!userForm._id && !payload.password)) {
    sectionError.users = "Name, email and password are required for new users.";
    return;
  }

  saving.user = true;
  try {
    if (userForm._id) {
      await updateAdminUser(userForm._id, payload);
      setFeedback("success", "User updated.");
    } else {
      await createAdminUser(payload);
      setFeedback("success", "User created.");
    }
    resetUserForm();
    showUserForm.value = false;
    await loadUsers();
  } catch (err) {
    sectionError.users = getErrorMessage(err, "Could not save user.");
  } finally {
    saving.user = false;
  }
};

const removeUser = async (user) => {
  if (!window.confirm(`Delete user "${user.name}"?`)) return;
  clearFeedback();
  sectionError.users = "";
  try {
    await deleteAdminUser(user._id);
    if (userForm._id === user._id) {
      resetUserForm();
    }
    syncUsersCache(
      users.value.filter((currentUser) => currentUser._id !== user._id),
    );
    setFeedback("success", "User deleted.");
  } catch (err) {
    sectionError.users = getErrorMessage(err, "Could not delete user.");
  }
};

const normalizeOptionsText = (options) => {
  if (typeof options === "string") return options || "{}";
  return JSON.stringify(options || {}, null, 2);
};

const editBooking = (booking) => {
  clearFeedback();
  activeSection.value = "bookings";
  bookingForm._id = booking._id;
  bookingForm.name = booking.name || "";
  bookingForm.email = booking.email || "";
  bookingForm.phone = booking.phone || "";
  bookingForm.time = toLocalDateTime(booking.time);
  bookingForm.guests = booking.guests || 1;
  bookingForm.channel = booking.channel || "admin";
  bookingForm.status = booking.status || "pending";
  bookingForm.items =
    booking.items?.length > 0
      ? booking.items.map((item) => ({
          drinkId: item.drinkId || "",
          qty: item.qty || 1,
          optionsText: normalizeOptionsText(item.options),
        }))
      : [createEmptyBookingItem()];
  showBookingForm.value = true;
};

const addBookingItem = () => {
  bookingForm.items.push(createEmptyBookingItem());
};

const removeBookingItem = (index) => {
  if (bookingForm.items.length === 1) {
    bookingForm.items = [createEmptyBookingItem()];
    return;
  }
  bookingForm.items.splice(index, 1);
};

const buildBookingPayload = () => {
  const time = toIsoDateTime(bookingForm.time);
  if (!time) {
    throw new Error("A valid booking time is required.");
  }

  const items = bookingForm.items
    .filter((item) => item.drinkId)
    .map((item) => {
      let options = {};
      const raw = item.optionsText?.trim();
      if (raw) {
        const parsed = JSON.parse(raw);
        options =
          parsed && typeof parsed === "object" && !Array.isArray(parsed)
            ? parsed
            : {};
      }

      return {
        drinkId: item.drinkId,
        qty: Math.max(1, Number(item.qty) || 1),
        options,
      };
    });

  return {
    name: bookingForm.name.trim(),
    email: bookingForm.email.trim(),
    phone: bookingForm.phone.trim(),
    time,
    guests: Math.max(1, Number(bookingForm.guests) || 1),
    channel: bookingForm.channel.trim() || "admin",
    status: bookingForm.status || "pending",
    items,
  };
};

const submitBooking = async () => {
  clearFeedback();
  sectionError.bookings = "";

  let payload;
  try {
    payload = buildBookingPayload();
  } catch (err) {
    sectionError.bookings = err.message || "Could not parse booking.";
    return;
  }

  if (!payload.name || !payload.email || !payload.phone) {
    sectionError.bookings = "Name, email and phone are required.";
    return;
  }

  saving.booking = true;
  try {
    if (bookingForm._id) {
      await updateAdminBooking(bookingForm._id, payload);
      setFeedback("success", "Booking updated.");
    } else {
      await createAdminBooking(payload);
      setFeedback("success", "Booking created.");
    }
    resetBookingForm();
    showBookingForm.value = false;
    await loadBookings();
  } catch (err) {
    sectionError.bookings = getErrorMessage(err, "Could not save booking.");
  } finally {
    saving.booking = false;
  }
};

const removeBooking = async (booking) => {
  if (!window.confirm(`Delete booking for "${booking.name}"?`)) return;
  clearFeedback();
  sectionError.bookings = "";
  try {
    await deleteAdminBooking(booking._id);
    if (bookingForm._id === booking._id) {
      resetBookingForm();
    }
    syncBookingsCache(
      bookings.value.filter(
        (currentBooking) => currentBooking._id !== booking._id,
      ),
    );
    setFeedback("success", "Booking deleted.");
  } catch (err) {
    sectionError.bookings = getErrorMessage(err, "Could not delete booking.");
  }
};

const editDrink = (drink) => {
  clearFeedback();
  activeSection.value = "products";
  Object.assign(drinkForm, {
    _id: drink._id,
    name: drink.name || "",
    price: drink.price || 0,
    stock: drink.stock ?? 0,
    available: Boolean(drink.available ?? (drink.stock ?? 0) > 0),
    tagsText: (drink.tags || []).join(", "),
    caffeine: drink.caffeine || "none",
    temp: drink.temp || "iced",
    sweetness: drink.sweetness ?? 5,
    colorTone: drink.colorTone || "neutral",
    calm: drink.emotionFit?.calm ?? 0,
    happy: drink.emotionFit?.happy ?? 0,
    stressed: drink.emotionFit?.stressed ?? 0,
    sad: drink.emotionFit?.sad ?? 0,
    adventurous: drink.emotionFit?.adventurous ?? 0,
    image: drink.image || "",
    desc: drink.desc || "",
  });
  showDrinkForm.value = true;
};

const buildDrinkPayload = () => ({
  name: drinkForm.name.trim(),
  price: Math.max(0, Number(drinkForm.price) || 0),
  stock: Math.max(0, Number(drinkForm.stock) || 0),
  available:
    Boolean(drinkForm.available) ||
    Math.max(0, Number(drinkForm.stock) || 0) > 0,
  tags: drinkForm.tagsText
    .split(",")
    .map((tag) => tag.trim())
    .filter(Boolean),
  caffeine: drinkForm.caffeine,
  temp: drinkForm.temp,
  sweetness: Math.max(0, Number(drinkForm.sweetness) || 0),
  colorTone: drinkForm.colorTone,
  emotionFit: {
    calm: Number(drinkForm.calm) || 0,
    happy: Number(drinkForm.happy) || 0,
    stressed: Number(drinkForm.stressed) || 0,
    sad: Number(drinkForm.sad) || 0,
    adventurous: Number(drinkForm.adventurous) || 0,
  },
  image: drinkForm.image.trim(),
  desc: drinkForm.desc.trim(),
});

const submitDrink = async () => {
  clearFeedback();
  sectionError.drinks = "";

  const payload = buildDrinkPayload();
  if (!payload.name) {
    sectionError.drinks = "Product name is required.";
    return;
  }

  saving.drink = true;
  try {
    if (drinkForm._id) {
      await updateAdminDrink(drinkForm._id, payload);
      setFeedback("success", "Product updated.");
    } else {
      await createAdminDrink(payload);
      setFeedback("success", "Product created.");
    }
    resetDrinkForm();
    showDrinkForm.value = false;
    await loadDrinks();
  } catch (err) {
    sectionError.drinks = getErrorMessage(err, "Could not save product.");
  } finally {
    saving.drink = false;
  }
};

const removeDrink = async (drink) => {
  if (!window.confirm(`Delete product "${drink.name}"?`)) return;
  clearFeedback();
  sectionError.drinks = "";
  try {
    await deleteAdminDrink(drink._id);
    if (drinkForm._id === drink._id) {
      resetDrinkForm();
    }
    syncDrinksCache(
      drinks.value.filter((currentDrink) => currentDrink._id !== drink._id),
    );
    setFeedback("success", "Product deleted.");
  } catch (err) {
    sectionError.drinks = getErrorMessage(err, "Could not delete product.");
  }
};

onMounted(async () => {
  await ensureActiveSectionLoaded();
  initialLoading.value = false;
});
</script>

<template>
  <section
    class="admin-page"
    :class="{
      'theme-dark': theme === 'dark',
      'theme-night': theme === 'night',
    }"
  >
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <h3>Admin</h3>
      </div>
      <nav class="sidebar-menu">
        <button
          class="menu-item"
          :class="{ active: activeSection === 'users' }"
          @click="activeSection = 'users'"
        >
          <span>Users</span>
          <span class="badge-count">{{ users.length }}</span>
        </button>
        <button
          class="menu-item"
          :class="{ active: activeSection === 'bookings' }"
          @click="activeSection = 'bookings'"
        >
          <span>Bookings</span>
          <span class="badge-count">{{ bookings.length }}</span>
        </button>
        <button
          class="menu-item"
          :class="{ active: activeSection === 'products' }"
          @click="activeSection = 'products'"
        >
          <span>Products</span>
          <span class="badge-count">{{ drinks.length }}</span>
        </button>
      </nav>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
      <div v-if="initialLoading" class="loading">
        <p>Loading admin data...</p>
      </div>
      <template v-else>
        <div class="content-header">
          <h1>
            <span v-if="activeSection === 'users'">User Management</span>
            <span v-else-if="activeSection === 'bookings'"
              >Booking Management</span
            >
            <span v-else>Product Management</span>
          </h1>
          <button
            class="btn btn-default"
            :disabled="isRefreshing"
            @click="refreshAll"
          >
            {{ isRefreshing ? "Refreshing..." : "Refresh" }}
          </button>
        </div>

        <div
          v-if="feedback.text"
          class="alert"
          :class="`alert-${feedback.type}`"
        >
          {{ feedback.text }}
        </div>

        <!-- USERS -->
        <div v-if="activeSection === 'users'" class="section">
          <div class="section-layout">
            <div v-if="showUserForm" class="modal-overlay" @click.self="closeUserForm">
              <div class="form-card modal-content">
                <h2>{{ userForm._id ? "Edit User" : "New User" }}</h2>
                <div v-if="sectionError.users" class="alert alert-error">
                  {{ sectionError.users }}
                </div>
                <form @submit.prevent="submitUser">
                  <div class="form-field">
                    <label>Name</label>
                    <input v-model="userForm.name" type="text" placeholder="Full name" />
                  </div>
                  <div class="form-field">
                    <label>Email</label>
                    <input v-model="userForm.email" type="email" placeholder="Email" />
                  </div>
                  <div class="form-field">
                    <label>Password</label>
                    <input v-model="userForm.password" type="password" :placeholder="userForm._id ? 'Leave blank to keep' : 'Password'" />
                  </div>
                  <div class="form-field">
                    <label>Role</label>
                    <select v-model="userForm.role">
                      <option value="user">User</option>
                      <option value="admin">Admin</option>
                    </select>
                  </div>
                  <div class="form-field checkbox">
                    <input v-model="userForm.verified" type="checkbox" id="user-verified" />
                    <label for="user-verified">Verified</label>
                  </div>
                  <div class="form-actions">
                    <button class="btn btn-primary" :disabled="saving.user">
                      {{ saving.user ? "Saving..." : userForm._id ? "Update" : "Create" }}
                    </button>
                    <button class="btn btn-default" type="button" @click="closeUserForm">
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            </div>

            <div class="list-card">
              <div
                style="
                  display: flex;
                  align-items: center;
                  justify-content: space-between;
                "
              >
                <h2>Users</h2>
                <div style="display: flex; gap: 8px; align-items: center">
                  <button class="btn btn-sm" @click="openUserForm">
                    Add New</button
                  ><button class="btn btn-sm" @click="loadUsers">
                    Refresh
                  </button>
                </div>
              </div>
              <div v-if="!sectionLoaded.users && loading.users" class="empty">
                Loading...
              </div>
              <div
                v-else-if="sectionLoaded.users && !users.length"
                class="empty"
              >
                No users
              </div>
              <div v-else class="list">
                <div v-for="user in users" :key="user._id" class="list-item">
                  <div class="item-main">
                    <p class="item-title">{{ user.name }}</p>
                    <p class="item-text email-text">{{ user.email }}</p>
                  </div>
                  <div class="item-meta">
                    <span class="tag">{{ user.role }}</span
                    ><span
                      class="tag"
                      :class="{ 'tag-verified': user.verified }"
                      >{{ user.verified ? "✓ Verified" : "Pending" }}</span
                    >
                  </div>
                  <div class="item-actions">
                    <button 
                    class="btn btn-sm" 
                    @click="editUser(user)">
                      Edit</button
                    ><button
                      class="btn btn-sm btn-danger"
                      @click="removeUser(user)"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- BOOKINGS -->
        <div v-if="activeSection === 'bookings'" class="section">
          <div class="section-layout">
            <div v-if="showBookingForm" class="modal-overlay" @click.self="closeBookingForm">
              <div class="form-card modal-content">
                <h2>{{ bookingForm._id ? "Edit Booking" : "New Booking" }}</h2>
                <div v-if="sectionError.bookings" class="alert alert-error">
                  {{ sectionError.bookings }}
                </div>
                <form @submit.prevent="submitBooking">
                  <div class="form-row">
                    <div class="form-field">
                      <label>Name</label>
                      <input v-model="bookingForm.name" type="text" />
                    </div>
                    <div class="form-field">
                      <label>Email</label>
                      <input v-model="bookingForm.email" type="email" />
                    </div>
                  </div>
                  <div class="form-row">
                    <div class="form-field">
                      <label>Phone</label>
                      <input v-model="bookingForm.phone" type="text" />
                    </div>
                    <div class="form-field">
                      <label>Date & Time</label>
                      <input v-model="bookingForm.time" type="datetime-local" />
                    </div>
                  </div>
                  <div class="form-row">
                    <div class="form-field">
                      <label>Guests</label>
                      <input v-model.number="bookingForm.guests" type="number" min="1" />
                    </div>
                    <div class="form-field">
                      <label>Status</label>
                      <select v-model="bookingForm.status">
                        <option value="pending">Pending</option>
                        <option value="confirmed">Confirmed</option>
                        <option value="completed">Completed</option>
                        <option value="cancelled">Cancelled</option>
                      </select>
                    </div>
                  </div>
                  <div class="form-actions">
                    <button class="btn btn-primary" :disabled="saving.booking">
                      {{ saving.booking ? "Saving..." : bookingForm._id ? "Update" : "Create" }}
                    </button>
                    <button class="btn btn-default" type="button" @click="closeBookingForm">
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            </div>

            <div class="list-card">
              <div
                style="
                  display: flex;
                  align-items: center;
                  justify-content: space-between;
                "
              >
                <h2>Bookings</h2>
                <div style="display: flex; gap: 8px">
                  <button class="btn btn-sm" @click="openBookingForm">
                    Add New</button
                  ><button class="btn btn-sm" @click="loadBookings">
                    Refresh
                  </button>
                </div>
              </div>
              <div
                v-if="!sectionLoaded.bookings && loading.bookings"
                class="empty"
              >
                Loading...
              </div>
              <div
                v-else-if="sectionLoaded.bookings && !bookings.length"
                class="empty"
              >
                No bookings
              </div>
              <div v-else class="list">
                <div
                  v-for="booking in bookings"
                  :key="booking._id"
                  class="list-item"
                >
                  <div class="item-main">
                    <p class="item-title">{{ booking.name }}</p>
                    <p class="item-text email-text">
                      {{ booking.email }} • {{ booking.phone }}
                    </p>
                  </div>
                  <div class="item-meta">
                    <span class="tag">{{ formatDate(booking.time) }}</span
                    ><span
                      class="tag"
                      :class="statusBadgeClass(booking.status)"
                      >{{ formatBookingStatus(booking.status) }}</span
                    >
                  </div>
                  <div class="item-actions">
                    <button class="btn btn-sm" @click="editBooking(booking)">
                      Edit</button
                    ><button
                      class="btn btn-sm btn-danger"
                      @click="removeBooking(booking)"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- PRODUCTS -->
        <div v-if="activeSection === 'products'" class="section">
          <div class="section-layout">
            <div v-if="showDrinkForm" class="modal-overlay" @click.self="closeDrinkForm">
              <div class="form-card modal-content">
                <h2>{{ drinkForm._id ? "Edit Product" : "New Product" }}</h2>
                <div v-if="sectionError.drinks" class="alert alert-error">
                  {{ sectionError.drinks }}
                </div>
                <form @submit.prevent="submitDrink">
                  <div class="form-row">
                    <div class="form-field">
                      <label>Name</label>
                      <input v-model="drinkForm.name" type="text" />
                    </div>
                    <div class="form-field">
                      <label>Price</label>
                      <input v-model.number="drinkForm.price" type="number" min="0" />
                    </div>
                  </div>
                  <div class="form-row">
                    <div class="form-field">
                      <label>Stock</label>
                      <input v-model.number="drinkForm.stock" type="number" min="0" />
                    </div>
                    <div class="form-field checkbox">
                      <input v-model="drinkForm.available" type="checkbox" id="prod-available" />
                      <label for="prod-available">Available</label>
                    </div>
                  </div>
                  <div class="form-actions">
                    <button class="btn btn-primary" :disabled="saving.drink">
                      {{ saving.drink ? "Saving..." : drinkForm._id ? "Update" : "Create" }}
                    </button>
                    <button class="btn btn-default" type="button" @click="closeDrinkForm">
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            </div>

            <div class="list-card">
              <div
                style="
                  display: flex;
                  align-items: center;
                  justify-content: space-between;
                "
              >
                <h2>Products</h2>
                <div style="display: flex; gap: 8px">
                  <button class="btn btn-sm" @click="openDrinkForm">
                    Add New</button
                  ><button class="btn btn-sm" @click="loadDrinks">
                    Refresh
                  </button>
                </div>
              </div>
              <div class="search-bar">
                <input
                  v-model="drinkListFilters.query"
                  type="text"
                  placeholder="Search products"
                  @keyup.enter="loadDrinks"
                />
              </div>
              <div v-if="!sectionLoaded.drinks && loading.drinks" class="empty">
                Loading...
              </div>
              <div
                v-else-if="sectionLoaded.drinks && !drinks.length"
                class="empty"
              >
                No products
              </div>
              <div v-else class="list">
                <div v-for="drink in drinks" :key="drink._id" class="list-item">
                  <div class="item-main product-main">
                    <img
                      v-if="drink.image"
                      :src="drink.image"
                      :alt="drink.name"
                      class="product-thumb"
                    />
                    <div v-else class="product-thumb product-thumb-fallback">
                      {{ (drink.name || "P").charAt(0).toUpperCase() }}
                    </div>
                    <div class="product-info">
                      <p class="item-title">{{ drink.name }}</p>
                      <p class="item-text product-price">
                        {{ formatCurrency(drink.price) }}
                      </p>
                    </div>
                  </div>
                  <div class="item-meta">
                    <span class="tag">{{ drink.caffeine }}</span
                    ><span class="tag">{{ drink.temp }}</span
                    ><span
                      class="tag"
                      :class="{ 'tag-success': drink.stock > 0 }"
                      >Stock: {{ drink.stock }}</span
                    >
                  </div>
                  <div class="item-actions">
                    <button class="btn btn-sm" @click="editDrink(drink)">
                      Edit</button
                    ><button
                      class="btn btn-sm btn-danger"
                      @click="removeDrink(drink)"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </main>
  </section>
</template>

<style scoped>
* {
  box-sizing: border-box;
}

.admin-page {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 0;
  height: 100vh;
  background: var(--paper);
  color: var(--ink);
}

/* SIDEBAR */
.sidebar {
  background: transparent;
  border-right: 1px solid rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.admin-page.theme-dark .sidebar {
  background: transparent;
  border-right-color: rgba(255, 255, 255, 0.1);
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.admin-page.theme-dark .sidebar-header {
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--tan);
}

.sidebar-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 8px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 600;
  transition: all 0.2s;
  border-left: 3px solid transparent;
}

.menu-item:hover {
  background: transparent;
}

.admin-page.theme-dark .menu-item:hover {
  background: transparent;
}

.menu-item.active {
  background: transparent;
  border-left-color: var(--tan);
}

.admin-page.theme-dark .menu-item.active {
  background: transparent;
}

.badge-count {
  display: inline-block;
  min-width: 24px;
  height: 24px;
  padding: 0 6px;
  background: var(--tan);
  color: white;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 24px;
  text-align: center;
}

/* MAIN CONTENT */
.main-content {
  overflow-y: auto;
  padding: 28px;
  background: var(--paper);
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: rgba(0, 0, 0, 0.5);
}

.admin-page.theme-dark .loading {
  color: rgba(255, 255, 255, 0.5);
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.content-header h1 {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--tan);
}

.alert {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 0.95rem;
  backdrop-filter: blur(5px);
}

.alert-success {
  background: rgba(214, 163, 90, 0.15);
  color: var(--tan);
  border: 1px solid rgba(214, 163, 90, 0.3);
}

.alert-error {
  background: rgba(240, 102, 44, 0.15);
  color: var(--orange-strong);
  border: 1px solid rgba(240, 102, 44, 0.3);
}

.section {
  margin-bottom: 24px;
}

.section-layout {
  display: block;
  
}

/* CARD */
.form-card,
.list-card {
  background: transparent;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 20px;
}

.admin-page.theme-dark .form-card,
.admin-page.theme-dark .list-card {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.1);
}

.form-card h2,
.list-card h2 {
  margin: 0 0 16px;
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--tan);
}

.list-card {
  display: flex;
  flex-direction: column;
}

.list-card h2 {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* FORM */
.form-field {
  margin-bottom: 14px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-field label {
  display: block;
  margin-bottom: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--ink);
}

.form-field input,
.form-field select,
.form-field textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
  background: transparent;
  color: var(--ink);
}

.admin-page.theme-dark .form-field input,
.admin-page.theme-dark .form-field select,
.admin-page.theme-dark .form-field textarea {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.15);
}

.form-field input:focus,
.form-field select:focus,
.form-field textarea:focus {
  outline: none;
  border-color: var(--tan);
  box-shadow: 0 0 0 2px rgba(184, 132, 67, 0.1);
}

.form-field.checkbox {
  display: flex;
  align-items: center;
  margin-bottom: 0;
}

.form-field.checkbox input {
  width: auto;
  margin-right: 8px;
}

.form-field.checkbox label {
  margin: 0;
}

.form-actions {
  display: flex;
  gap: 8px;
  margin-top: 20px;
}

/* LIST */
.empty {
  text-align: center;
  padding: 40px 20px;
  color: rgba(0, 0, 0, 0.4);
}

.admin-page.theme-dark .empty {
  color: rgba(255, 255, 255, 0.4);
}

.list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.list-item {
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 12px;
  align-items: center;
  padding: 12px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  background: transparent;
  transition: all 0.2s;
}

.admin-page.theme-dark .list-item {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.1);
}

.list-item:hover {
  background: rgba(255, 255, 255, 0.6);
  border-color: rgba(184, 132, 67, 0.2);
}

.admin-page.theme-dark .list-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(214, 163, 90, 0.3);
}

.item-main {
  min-width: 0;
}

.product-main {
  display: flex;
  align-items: center;
  gap: 10px;
}

.product-thumb {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(0, 0, 0, 0.12);
  flex-shrink: 0;
}

.admin-page.theme-dark .product-thumb {
  border-color: rgba(255, 255, 255, 0.16);
}

.product-thumb-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--tan);
  background: rgba(184, 132, 67, 0.12);
  border-radius: 50%;
}

.product-info {
  min-width: 0;
}

.item-title {
  margin: 0;
  font-size: 0.95rem;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--ink);
}

.item-text {
  margin: 4px 0 0;
  font-size: 0.85rem;
  color: rgba(0, 0, 0, 0.5);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.admin-page.theme-dark .item-text {
  color: rgba(255, 255, 255, 0.5);
}

.email-text {
  color: #111111;
}

.admin-page.theme-dark .email-text {
  color: #ffffff;
}

.admin-page.theme-night .email-text {
  color: #ffffff;
}

.product-price {
  color: #111111;
}

.admin-page.theme-dark .product-price {
  color: #ffffff;
}

.admin-page.theme-night .product-price {
  color: #ffffff;
}

.item-meta {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.tag {
  display: inline-block;
  padding: 4px 8px;
  background: rgba(184, 132, 67, 0.15);
  color: var(--tan);
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
}

.admin-page.theme-dark .tag {
  background: rgba(214, 163, 90, 0.2);
}

.tag-success {
  background: rgba(240, 102, 44, 0.15);
  color: var(--orange);
}

.admin-page.theme-dark .tag-success {
  background: rgba(240, 102, 44, 0.2);
  color: var(--orange);
}

.tag-verified {
  background: rgba(76, 175, 80, 0.15);
  color: #4caf50;
}

.admin-page.theme-dark .tag-verified {
  background: rgba(76, 175, 80, 0.2);
  color: #81c784;
}

.item-actions {
  display: flex;
  gap: 6px;
}

.search-bar {
  margin-top: 12px;
  margin-bottom: 12px;
}

.search-bar input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  font-size: 0.9rem;
  background: transparent;
  color: var(--ink);
}

.admin-page.theme-dark .search-bar input {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.15);
}

.search-bar input:focus {
  outline: none;
  border-color: var(--tan);
}

/* BUTTONS */
.btn {
  padding: 8px 16px;
  border: 1px solid rgba(184, 132, 67, 0.4);
  border-radius: 6px;
  background: transparent;
  color: var(--tan);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.9rem;
}

.admin-page.theme-dark .btn {
  background: transparent;
  border-color: rgba(214, 163, 90, 0.4);
}

.btn:hover:not(:disabled) {
  border-color: var(--tan);
  background: rgba(184, 132, 67, 0.1);
}

.admin-page.theme-dark .btn:hover:not(:disabled) {
  background: rgba(214, 163, 90, 0.15);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-default {
  background: transparent;
}

.btn-default:hover:not(:disabled) {
  background: rgba(184, 132, 67, 0.1);
}

.admin-page.theme-dark .btn-default {
  background: transparent;
}

.admin-page.theme-dark .btn-default:hover:not(:disabled) {
  background: rgba(214, 163, 90, 0.1);
}

.btn-primary {
  background: var(--tan);
  color: #fff;
  border-color: var(--tan);
}

.btn-primary:hover:not(:disabled) {
  background: var(--tan);
  opacity: 0.9;
}

.btn-danger {
  color: var(--orange);
  border-color: var(--orange);
  background: rgba(240, 102, 44, 0.05);
}

.admin-page.theme-dark .btn-danger {
  background: rgba(240, 102, 44, 0.1);
}

.btn-danger:hover:not(:disabled) {
  background: rgba(240, 102, 44, 0.15);
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.85rem;
}

/* =========================================
   MODAL / POPUP STYLES
   ========================================= */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.6); /* Nền đen mờ */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  backdrop-filter: blur(3px); /* Làm mờ nhẹ background phía sau */
}

/* Đè lên class form-card cũ để căn giữa và thêm màu nền solid */
.modal-content.form-card {
  background: var(--paper); /* Bắt buộc để không bị nhìn xuyên thấu */
  width: 90%;
  max-width: 450px;
  max-height: 90vh;
  overflow-y: auto; /* Tạo thanh cuộn nếu form quá dài */
  margin: 0; 
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.25);
  animation: modalFadeIn 0.25s ease-out;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

/* Tinh chỉnh đổ bóng và viền cho chế độ Dark Mode */
.admin-page.theme-dark .modal-content.form-card {
  border-color: rgba(255, 255, 255, 0.15);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.6);
}

/* Hiệu ứng trượt và mờ dần khi mở Popup */
@keyframes modalFadeIn {
  from { 
    opacity: 0; 
    transform: translateY(-20px); 
  }
  to { 
    opacity: 1; 
    transform: translateY(0); 
  }
}
</style>
