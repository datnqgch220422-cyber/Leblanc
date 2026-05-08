<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getBookingsByEmail, cancelBooking } from "@/api/bookings.api";
import { getDrinks } from "@/api/drinks.api";
import { getSessionUser } from "@/composables/useSessionAuth";
import ConfirmDialog from "@/components/ConfirmDialog.vue";

const route = useRoute();
const router = useRouter();
const bookingId = route.params.id;

const user = ref(getSessionUser());
const bookings = ref([]);
const drinks = ref([]);
const loading = ref(false);
const error = ref("");
const message = reactive({
  type: "",
  text: "",
});
const confirmDialog = reactive({
  open: false,
  pending: false,
});

const booking = computed(
  () => bookings.value.find((b) => String(b._id) === String(bookingId)) || null,
);

const detailItems = computed(() => {
  const items = Array.isArray(booking.value?.items) ? booking.value.items : [];
  return items.map((item) => {
    const drink = drinks.value.find((d) => d._id === item.drinkId) || {};
    const qty = Number(item.qty) || 0;
    const price = Number(drink.price) || 0;
    return {
      ...item,
      name: drink.name || item.drinkId,
      price,
      lineTotal: qty * price,
    };
  });
});

const totals = computed(() => ({
  count: detailItems.value.reduce((s, it) => s + (Number(it.qty) || 0), 0),
  total: detailItems.value.reduce((s, it) => s + it.lineTotal, 0),
}));

const formatCurrency = (v) => `${(Number(v) || 0).toLocaleString("en-US")} VND`;
const formatDate = (v) => {
  if (!v) return "—";
  const d = new Date(v);
  return Number.isNaN(d.getTime()) ? v : d.toLocaleString();
};

const load = async () => {
  if (!user.value?.email) return;
  loading.value = true;
  error.value = "";
  try {
    const [drinkList, bookingList] = await Promise.all([
      getDrinks(),
      getBookingsByEmail(user.value.email),
    ]);
    drinks.value = Array.isArray(drinkList) ? drinkList : [];
    bookings.value = Array.isArray(bookingList) ? bookingList : [];
  } catch (err) {
    error.value =
      err?.response?.data?.error || err?.message || "Could not load booking.";
  } finally {
    loading.value = false;
  }
};

const goBack = () => router.push({ name: "orders" });

const canCancel = (status) => {
  const s = (status || "").toLowerCase();
  return s === "pending" || s === "confirmed";
};

const handleCancel = async () => {
  if (!user.value?.email) {
    message.type = "error";
    message.text = "You must be logged in to cancel.";
    return;
  }
  const b = booking.value;
  if (!b || !canCancel(b.status)) return;
  confirmDialog.open = true;
};

const closeConfirmDialog = () => {
  if (confirmDialog.pending) return;
  confirmDialog.open = false;
};

const confirmCancel = async () => {
  const b = booking.value;
  if (!b) return;
  confirmDialog.pending = true;
  try {
    await cancelBooking(b._id, user.value.email);
    await load();
    message.type = "success";
    message.text = "Booking cancelled.";
    confirmDialog.open = false;
    goBack();
  } catch (err) {
    message.type = "error";
    message.text =
      err?.response?.data?.error || err?.message || "Could not cancel booking.";
  } finally {
    confirmDialog.pending = false;
  }
};

onMounted(load);
</script>

<template>
  <section class="booking-detail">
    <header class="booking-detail__header">
      <div>
        <p class="eyebrow">Le'Blanc</p>
        <h1>Booking detail</h1>
      </div>
      <button @click="goBack" class="btn btn--link">← Back to orders</button>
    </header>

    <div v-if="message.text" :class="['notice', `notice--${message.type}`]">
      {{ message.text }}
    </div>

    <div v-if="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="!booking" class="empty-state">Booking not found.</div>

    <div v-else class="booking-card">
      <h2>{{ formatDate(booking.time) }}</h2>
      <p><strong>Name:</strong> {{ booking.name || "—" }}</p>
      <p><strong>Phone:</strong> {{ booking.phone || "—" }}</p>
      <p><strong>Email:</strong> {{ booking.email || "—" }}</p>
      <p>
        <strong>Guests:</strong> {{ booking.guests || booking.guest || "—" }}
      </p>
      <p><strong>Channel:</strong> {{ booking.channel || "web" }}</p>

      <section class="items">
        <h3>Items</h3>
        <ul>
          <li v-for="item in detailItems" :key="item.drinkId">
            <span>{{ item.name }}</span>
            <span>x{{ item.qty }}</span>
            <span>{{ formatCurrency(item.price) }}</span>
            <span>{{ formatCurrency(item.lineTotal) }}</span>
          </li>
        </ul>
        <div class="summary">
          <p>Total items: {{ totals.count }}</p>
          <p>Total drinks price: {{ formatCurrency(totals.total) }}</p>
        </div>
        <div style="margin-top: 12px">
          <button
            v-if="canCancel(booking.status)"
            class="btn btn--danger"
            @click="handleCancel"
          >
            Cancel booking
          </button>
        </div>
      </section>
    </div>

    <ConfirmDialog
      v-model="confirmDialog.open"
      title="Cancel booking"
      message="Are you sure you want to cancel this booking?"
      :confirm-text="confirmDialog.pending ? 'Cancelling...' : 'Cancel booking'"
      cancel-text="Keep booking"
      :pending="confirmDialog.pending"
      :danger="true"
      @confirm="confirmCancel"
      @cancel="closeConfirmDialog"
    />
  </section>
</template>

<style scoped>
.booking-detail__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.booking-card {
  padding: 12px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 8px;
}
.items ul {
  list-style: none;
  padding: 0;
}
.items li {
  display: grid;
  grid-template-columns: 1fr auto auto auto;
  gap: 8px;
  padding: 6px 0;
}
.summary {
  margin-top: 8px;
  font-weight: 800;
}
.error {
  color: #9d3412;
}
.notice {
  margin: 10px 0 14px;
  padding: 10px 12px;
  border-radius: 8px;
  font-weight: 600;
}
.notice--success {
  color: #18592a;
  background: rgba(24, 89, 42, 0.12);
}
.notice--error {
  color: #9d3412;
  background: rgba(157, 52, 18, 0.12);
}
.btn {
  border: 1px solid rgba(0, 0, 0, 0.15);
  background: #fff;
  border-radius: 8px;
  padding: 8px 12px;
  font-weight: 700;
  cursor: pointer;
}
.btn--danger {
  border-color: #9d3412;
  color: #9d3412;
}
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
