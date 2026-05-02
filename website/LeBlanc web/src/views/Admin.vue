<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import {
  createAdminBooking,
  createAdminDrink,
  createAdminUser,
  deleteAdminBooking,
  deleteAdminDrink,
  deleteAdminUser,
  getAdminBookings,
  getAdminDrinks,
  getAdminUsers,
  updateAdminBooking,
  updateAdminDrink,
  updateAdminUser,
} from "@/api";
import { useThemeState } from "@/composables/useThemeState";

const { theme } = useThemeState();

const activeSection = ref("users");
const users = ref([]);
const bookings = ref([]);
const drinks = ref([]);
const initialLoading = ref(true);

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
});

const createEmptyDrinkForm = () => ({
  _id: "",
  name: "",
  price: 0,
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
  bookingForm.items = next.items;
};

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

const loadUsers = async () => {
  loading.users = true;
  sectionError.users = "";
  try {
    users.value = await getAdminUsers();
  } catch (err) {
    sectionError.users = getErrorMessage(err, "Could not load users.");
  } finally {
    loading.users = false;
  }
};

const loadBookings = async () => {
  loading.bookings = true;
  sectionError.bookings = "";
  try {
    bookings.value = await getAdminBookings();
  } catch (err) {
    sectionError.bookings = getErrorMessage(err, "Could not load bookings.");
  } finally {
    loading.bookings = false;
  }
};

const loadDrinks = async () => {
  loading.drinks = true;
  sectionError.drinks = "";
  try {
    drinks.value = await getAdminDrinks();
  } catch (err) {
    sectionError.drinks = getErrorMessage(err, "Could not load products.");
  } finally {
    loading.drinks = false;
  }
};

const refreshAll = async () => {
  clearFeedback();
  await Promise.allSettled([loadUsers(), loadBookings(), loadDrinks()]);
};

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
    setFeedback("success", "User deleted.");
    await loadUsers();
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
  bookingForm.items =
    booking.items?.length > 0
      ? booking.items.map((item) => ({
          drinkId: item.drinkId || "",
          qty: item.qty || 1,
          optionsText: normalizeOptionsText(item.options),
        }))
      : [createEmptyBookingItem()];
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
    setFeedback("success", "Booking deleted.");
    await loadBookings();
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
};

const buildDrinkPayload = () => ({
  name: drinkForm.name.trim(),
  price: Math.max(0, Number(drinkForm.price) || 0),
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
    setFeedback("success", "Product deleted.");
    await loadDrinks();
  } catch (err) {
    sectionError.drinks = getErrorMessage(err, "Could not delete product.");
  }
};

onMounted(async () => {
  await refreshAll();
  initialLoading.value = false;
});
</script>

<template>
  <section class="admin-page">
    <header class="hero">
      <div>
        <p class="eyebrow">ADMIN PAGE</p>
        <h1>Manage users, bookings, and menu products.</h1>
        <p class="lede">
          Full CRUD for the three main collections behind Le'Blanc.
        </p>
      </div>
      <button
        class="btn"
        :class="{ 'is-night': theme === 'night' }"
        type="button"
        :disabled="isRefreshing"
        @click="refreshAll"
      >
        {{ isRefreshing ? "Refreshing..." : "Refresh all" }}
      </button>
    </header>

    <div v-if="feedback.text" class="banner" :class="feedback.type">
      {{ feedback.text }}
    </div>

    <div class="tabs">
      <button
        class="tab"
        :class="{ active: activeSection === 'users' }"
        type="button"
        @click="activeSection = 'users'"
      >
        Users
        <span class="count">{{ users.length }}</span>
      </button>
      <button
        class="tab"
        :class="{ active: activeSection === 'bookings' }"
        type="button"
        @click="activeSection = 'bookings'"
      >
        Bookings
        <span class="count">{{ bookings.length }}</span>
      </button>
      <button
        class="tab"
        :class="{ active: activeSection === 'products' }"
        type="button"
        @click="activeSection = 'products'"
      >
        Products
        <span class="count">{{ drinks.length }}</span>
      </button>
    </div>

    <div v-if="initialLoading" class="panel state">
      Loading admin data...
    </div>

    <template v-else>
      <div v-if="activeSection === 'users'" class="panel-grid">
        <article class="panel">
          <div class="panel-head">
            <div>
              <p class="panel-kicker">Users</p>
              <h2>{{ userForm._id ? "Edit user" : "Create user" }}</h2>
            </div>
            <button class="btn subtle" type="button" @click="resetUserForm">
              New
            </button>
          </div>

          <p v-if="sectionError.users" class="banner error">
            {{ sectionError.users }}
          </p>

          <form class="form-grid" @submit.prevent="submitUser">
            <label>
              <span>Name</span>
              <input v-model="userForm.name" type="text" />
            </label>
            <label>
              <span>Email</span>
              <input v-model="userForm.email" type="email" />
            </label>
            <label>
              <span>Password</span>
              <input
                v-model="userForm.password"
                type="password"
                :placeholder="userForm._id ? 'Leave blank to keep current' : ''"
              />
            </label>
            <label>
              <span>Role</span>
              <select v-model="userForm.role">
                <option value="user">User</option>
                <option value="admin">Admin</option>
              </select>
            </label>
            <label class="checkbox">
              <input v-model="userForm.verified" type="checkbox" />
              <span>Verified account</span>
            </label>

            <div class="form-actions">
              <button class="btn primary" type="submit" :disabled="saving.user">
                {{
                  saving.user
                    ? "Saving..."
                    : userForm._id
                      ? "Update user"
                      : "Create user"
                }}
              </button>
              <button
                v-if="userForm._id"
                class="btn subtle"
                type="button"
                @click="resetUserForm"
              >
                Cancel
              </button>
            </div>
          </form>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div>
              <p class="panel-kicker">Directory</p>
              <h2>Registered users</h2>
            </div>
            <button class="btn subtle" type="button" @click="loadUsers">
              Refresh
            </button>
          </div>

          <div v-if="loading.users" class="state">Loading users...</div>
          <div v-else-if="!users.length" class="state">No users found.</div>
          <div v-else class="record-list">
            <article v-for="user in users" :key="user._id" class="record-card">
              <div class="record-main">
                <div>
                  <p class="record-title">{{ user.name }}</p>
                  <p class="record-sub">{{ user.email }}</p>
                </div>
                <div class="badges">
                  <span class="badge">{{ user.role || "user" }}</span>
                  <span class="badge" :class="{ ok: user.verified }">
                    {{ user.verified ? "Verified" : "Pending" }}
                  </span>
                </div>
              </div>
              <p class="record-meta">Joined {{ formatDate(user.createdAt) }}</p>
              <div class="record-actions">
                <button class="btn subtle" type="button" @click="editUser(user)">
                  Edit
                </button>
                <button class="btn danger" type="button" @click="removeUser(user)">
                  Delete
                </button>
              </div>
            </article>
          </div>
        </article>
      </div>

      <div v-else-if="activeSection === 'bookings'" class="panel-grid">
        <article class="panel">
          <div class="panel-head">
            <div>
              <p class="panel-kicker">Bookings</p>
              <h2>{{ bookingForm._id ? "Edit booking" : "Create booking" }}</h2>
            </div>
            <button class="btn subtle" type="button" @click="resetBookingForm">
              New
            </button>
          </div>

          <p v-if="sectionError.bookings" class="banner error">
            {{ sectionError.bookings }}
          </p>

          <form class="form-grid" @submit.prevent="submitBooking">
            <label>
              <span>Name</span>
              <input v-model="bookingForm.name" type="text" />
            </label>
            <label>
              <span>Email</span>
              <input v-model="bookingForm.email" type="email" />
            </label>
            <label>
              <span>Phone</span>
              <input v-model="bookingForm.phone" type="text" />
            </label>
            <label>
              <span>Arrival</span>
              <input v-model="bookingForm.time" type="datetime-local" />
            </label>
            <label>
              <span>Guests</span>
              <input v-model.number="bookingForm.guests" type="number" min="1" />
            </label>
            <label>
              <span>Channel</span>
              <input v-model="bookingForm.channel" type="text" />
            </label>

            <div class="item-editor wide">
              <div class="item-editor-head">
                <p>Order items</p>
                <button class="btn subtle" type="button" @click="addBookingItem">
                  Add item
                </button>
              </div>

              <div
                v-for="(item, index) in bookingForm.items"
                :key="`${index}-${item.drinkId}`"
                class="item-card"
              >
                <label>
                  <span>Drink</span>
                  <select v-model="item.drinkId">
                    <option value="">Select a drink</option>
                    <option
                      v-for="drink in drinks"
                      :key="drink._id"
                      :value="drink._id"
                    >
                      {{ drink.name }}
                    </option>
                  </select>
                </label>
                <label>
                  <span>Qty</span>
                  <input v-model.number="item.qty" type="number" min="1" />
                </label>
                <label class="wide">
                  <span>Options JSON</span>
                  <textarea v-model="item.optionsText" rows="4"></textarea>
                </label>
                <button
                  class="btn danger"
                  type="button"
                  @click="removeBookingItem(index)"
                >
                  Remove item
                </button>
              </div>
            </div>

            <div class="form-actions wide">
              <button
                class="btn primary"
                type="submit"
                :disabled="saving.booking"
              >
                {{
                  saving.booking
                    ? "Saving..."
                    : bookingForm._id
                      ? "Update booking"
                      : "Create booking"
                }}
              </button>
              <button
                v-if="bookingForm._id"
                class="btn subtle"
                type="button"
                @click="resetBookingForm"
              >
                Cancel
              </button>
            </div>
          </form>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div>
              <p class="panel-kicker">Orders</p>
              <h2>All bookings</h2>
            </div>
            <button class="btn subtle" type="button" @click="loadBookings">
              Refresh
            </button>
          </div>

          <div v-if="loading.bookings" class="state">Loading bookings...</div>
          <div v-else-if="!bookings.length" class="state">No bookings found.</div>
          <div v-else class="record-list">
            <article
              v-for="booking in bookings"
              :key="booking._id"
              class="record-card"
            >
              <div class="record-main">
                <div>
                  <p class="record-title">{{ booking.name }}</p>
                  <p class="record-sub">
                    {{ booking.email }} · {{ booking.phone }}
                  </p>
                </div>
                <span class="badge">{{ booking.channel || "web" }}</span>
              </div>
              <p class="record-meta">
                {{ formatDate(booking.time) }} · {{ booking.guests || 1 }} guests
              </p>
              <p class="record-meta">{{ summarizeBookingItems(booking.items) }}</p>
              <div class="record-actions">
                <button
                  class="btn subtle"
                  type="button"
                  @click="editBooking(booking)"
                >
                  Edit
                </button>
                <button
                  class="btn danger"
                  type="button"
                  @click="removeBooking(booking)"
                >
                  Delete
                </button>
              </div>
            </article>
          </div>
        </article>
      </div>

      <div v-else class="panel-grid">
        <article class="panel">
          <div class="panel-head">
            <div>
              <p class="panel-kicker">Products</p>
              <h2>{{ drinkForm._id ? "Edit drink" : "Create drink" }}</h2>
            </div>
            <button class="btn subtle" type="button" @click="resetDrinkForm">
              New
            </button>
          </div>

          <p v-if="sectionError.drinks" class="banner error">
            {{ sectionError.drinks }}
          </p>

          <form class="form-grid" @submit.prevent="submitDrink">
            <label>
              <span>Name</span>
              <input v-model="drinkForm.name" type="text" />
            </label>
            <label>
              <span>Price</span>
              <input v-model.number="drinkForm.price" type="number" min="0" />
            </label>
            <label class="wide">
              <span>Tags</span>
              <input
                v-model="drinkForm.tagsText"
                type="text"
                placeholder="day, coffee, iced"
              />
            </label>
            <label>
              <span>Caffeine</span>
              <select v-model="drinkForm.caffeine">
                <option value="none">None</option>
                <option value="low">Low</option>
                <option value="med">Medium</option>
                <option value="high">High</option>
              </select>
            </label>
            <label>
              <span>Temperature</span>
              <select v-model="drinkForm.temp">
                <option value="hot">Hot</option>
                <option value="iced">Iced</option>
                <option value="cold">Cold</option>
                <option value="either">Either</option>
                <option value="room">Room</option>
              </select>
            </label>
            <label>
              <span>Sweetness</span>
              <input
                v-model.number="drinkForm.sweetness"
                type="number"
                min="0"
                max="10"
              />
            </label>
            <label>
              <span>Color tone</span>
              <select v-model="drinkForm.colorTone">
                <option value="warm">Warm</option>
                <option value="cool">Cool</option>
                <option value="neutral">Neutral</option>
              </select>
            </label>
            <label class="wide">
              <span>Image URL</span>
              <input v-model="drinkForm.image" type="text" />
            </label>
            <label class="wide">
              <span>Description</span>
              <textarea v-model="drinkForm.desc" rows="4"></textarea>
            </label>

            <div class="emotion-grid wide">
              <label>
                <span>Calm</span>
                <input
                  v-model.number="drinkForm.calm"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                />
              </label>
              <label>
                <span>Happy</span>
                <input
                  v-model.number="drinkForm.happy"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                />
              </label>
              <label>
                <span>Stressed</span>
                <input
                  v-model.number="drinkForm.stressed"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                />
              </label>
              <label>
                <span>Sad</span>
                <input
                  v-model.number="drinkForm.sad"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                />
              </label>
              <label>
                <span>Adventurous</span>
                <input
                  v-model.number="drinkForm.adventurous"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                />
              </label>
            </div>

            <div class="form-actions wide">
              <button
                class="btn primary"
                type="submit"
                :disabled="saving.drink"
              >
                {{
                  saving.drink
                    ? "Saving..."
                    : drinkForm._id
                      ? "Update product"
                      : "Create product"
                }}
              </button>
              <button
                v-if="drinkForm._id"
                class="btn subtle"
                type="button"
                @click="resetDrinkForm"
              >
                Cancel
              </button>
            </div>
          </form>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div>
              <p class="panel-kicker">Menu</p>
              <h2>Current drinks</h2>
            </div>
            <button class="btn subtle" type="button" @click="loadDrinks">
              Refresh
            </button>
          </div>

          <div v-if="loading.drinks" class="state">Loading products...</div>
          <div v-else-if="!drinks.length" class="state">No products found.</div>
          <div v-else class="record-list">
            <article
              v-for="drink in drinks"
              :key="drink._id"
              class="record-card"
            >
              <div class="record-main">
                <div>
                  <p class="record-title">{{ drink.name }}</p>
                  <p class="record-sub">{{ formatCurrency(drink.price) }}</p>
                </div>
                <div class="badges">
                  <span class="badge">{{ drink.caffeine || "none" }}</span>
                  <span class="badge">{{ drink.temp || "n/a" }}</span>
                </div>
              </div>
              <p class="record-meta">{{ (drink.tags || []).join(", ") || "No tags" }}</p>
              <p class="record-meta">{{ drink.desc || "No description" }}</p>
              <div class="record-actions">
                <button
                  class="btn subtle"
                  type="button"
                  @click="editDrink(drink)"
                >
                  Edit
                </button>
                <button
                  class="btn danger"
                  type="button"
                  @click="removeDrink(drink)"
                >
                  Delete
                </button>
              </div>
            </article>
          </div>
        </article>
      </div>
    </template>
  </section>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 18px;
}

.hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
}

.eyebrow {
  margin: 0;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  font-size: 0.86rem;
  color: var(--tan);
}

h1 {
  margin: 8px 0 10px;
  font-size: clamp(2rem, 4vw, 3rem);
  line-height: 1.05;
}

.lede {
  margin: 0;
  max-width: 62ch;
  color: var(--ink);
  opacity: 0.82;
}

.tabs {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.tab {
  padding: 12px 16px;
  border-radius: 999px;
  border: 1px solid rgba(0, 0, 0, 0.16);
  background: rgba(255, 255, 255, 0.66);
  color: var(--ink);
  cursor: pointer;
  font-weight: 800;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.tab.active {
  background: #0f1424;
  color: #f6efe6;
}

.count {
  min-width: 24px;
  height: 24px;
  border-radius: 999px;
  display: inline-grid;
  place-items: center;
  font-size: 0.82rem;
  background: rgba(0, 0, 0, 0.08);
}

.tab.active .count {
  background: rgba(255, 255, 255, 0.16);
}

.panel-grid {
  display: grid;
  grid-template-columns: minmax(320px, 0.95fr) minmax(420px, 1.25fr);
  gap: 18px;
  align-items: start;
}

.panel {
  display: grid;
  gap: 16px;
  padding: 22px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.08);
}

.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.panel-kicker {
  margin: 0;
  font-size: 0.8rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--tan);
}

.panel h2 {
  margin: 6px 0 0;
  font-size: 1.55rem;
  line-height: 1.15;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.form-grid label,
.emotion-grid label {
  display: grid;
  gap: 6px;
}

.form-grid span,
.emotion-grid span {
  font-size: 0.9rem;
  font-weight: 800;
}

.form-grid input,
.form-grid select,
.form-grid textarea,
.emotion-grid input {
  width: 100%;
  padding: 11px 12px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.14);
  background: rgba(255, 255, 255, 0.84);
  color: var(--ink);
  font: inherit;
  box-sizing: border-box;
}

.wide {
  grid-column: 1 / -1;
}

.checkbox {
  display: flex !important;
  align-items: center;
  gap: 10px;
  font-weight: 800;
}

.checkbox input {
  width: 18px;
  height: 18px;
}

.form-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  padding: 11px 15px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.16);
  background: #0f1424;
  color: #f6efe6;
  font-weight: 800;
  cursor: pointer;
}

.btn.is-night,
.btn.primary {
  background: var(--tan);
  color: #0f1424;
}

.btn.subtle {
  background: rgba(0, 0, 0, 0.06);
  color: var(--ink);
}

.btn.danger {
  background: rgba(176, 0, 32, 0.1);
  color: #8e1026;
  border-color: rgba(176, 0, 32, 0.22);
}

.btn:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.banner {
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.banner.error {
  color: #9f2a13;
  background: rgba(220, 90, 36, 0.12);
  border-color: rgba(220, 90, 36, 0.22);
}

.banner.success {
  color: #155f35;
  background: rgba(30, 128, 70, 0.12);
  border-color: rgba(30, 128, 70, 0.22);
}

.state {
  min-height: 120px;
  display: grid;
  place-items: center;
  text-align: center;
  color: var(--ink);
  opacity: 0.8;
}

.record-list {
  display: grid;
  gap: 12px;
}

.record-card {
  display: grid;
  gap: 10px;
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.76);
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.record-main {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.record-title {
  margin: 0;
  font-size: 1.06rem;
  font-weight: 900;
}

.record-sub,
.record-meta {
  margin: 0;
  color: var(--ink);
  opacity: 0.8;
  line-height: 1.45;
}

.record-actions,
.badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.badge {
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.06);
  font-size: 0.82rem;
  font-weight: 800;
}

.badge.ok {
  background: rgba(20, 110, 58, 0.14);
  color: #146e3a;
}

.item-editor {
  display: grid;
  gap: 12px;
}

.item-editor-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.item-editor-head p {
  margin: 0;
  font-weight: 900;
}

.item-card {
  display: grid;
  grid-template-columns: 1.2fr 0.5fr;
  gap: 12px;
  padding: 14px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.emotion-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
}

@media (max-width: 1080px) {
  .panel-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .hero,
  .panel-head,
  .record-main,
  .item-editor-head {
    flex-direction: column;
  }

  .form-grid,
  .emotion-grid,
  .item-card {
    grid-template-columns: 1fr;
  }
}
</style>
