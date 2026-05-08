<script setup>
import { computed, reactive, ref, onMounted } from "vue";
import { RouterLink } from "vue-router";
import { getBookingsByEmail, cancelBooking } from "@/api/bookings.api";
import { getSessionUser } from "@/composables/useSessionAuth";
import { getDrinks } from "@/api/drinks.api";
import ConfirmDialog from "@/components/ConfirmDialog.vue";

const user = ref(getSessionUser());
const bookings = ref([]);
const drinks = ref([]);
const loading = ref(false);
const error = ref("");
const selectedStatus = ref(null);
const message = reactive({
  type: "",
  text: "",
});
const confirmDialog = reactive({
  open: false,
  booking: null,
  pending: false,
});

const statusOptions = [
  { value: null, label: "All" },
  { value: "pending", label: "Pending" },
  { value: "confirmed", label: "Confirmed" },
  { value: "completed", label: "Completed" },
  { value: "cancelled", label: "Cancelled" },
];

const summarizeItems = (items = []) =>
  items.length
    ? items
        .map((item) => {
          const drink = drinks.value.find((d) => d._id === item.drinkId) || {};
          return `${drink.name || item.drinkId} x${item.qty || 1}`;
        })
        .join(", ")
    : "No items";

const bookingSummary = (booking) => {
  const items = Array.isArray(booking?.items) ? booking.items : [];
  const lineItems = items.map((item) => {
    const drink = drinks.value.find((d) => d._id === item.drinkId) || {};
    const price = Number(drink.price) || 0;
    const qty = Number(item.qty) || 0;
    return {
      name: drink.name || item.drinkId,
      qty,
      price,
      lineTotal: price * qty,
    };
  });

  return {
    count: lineItems.reduce((sum, item) => sum + item.qty, 0),
    total: lineItems.reduce((sum, item) => sum + item.lineTotal, 0),
  };
};

const formatCurrency = (value) =>
  `${(Number(value) || 0).toLocaleString("en-US")} VND`;

const formatStatus = (status) => {
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

const statusClass = (status) => {
  switch ((status || "pending").toLowerCase()) {
    case "confirmed":
      return "status status--confirmed";
    case "completed":
      return "status status--completed";
    case "cancelled":
      return "status status--cancelled";
    default:
      return "status status--pending";
  }
};

const bookingsWithSummary = computed(() =>
  bookings.value.map((booking) => ({
    ...booking,
    summary: bookingSummary(booking),
  })),
);

const filteredBookings = computed(() => {
  if (!selectedStatus.value) return bookingsWithSummary.value;
  return bookingsWithSummary.value.filter(
    (b) =>
      (b.status || "pending").toLowerCase() ===
      selectedStatus.value.toLowerCase(),
  );
});

const formatDate = (value) => {
  if (!value) return "—";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
};

const load = async () => {
  if (!user.value?.email) return;
  loading.value = true;
  error.value = "";
  try {
    const [dlist, blist] = await Promise.all([
      getDrinks(),
      getBookingsByEmail(user.value.email),
    ]);
    drinks.value = Array.isArray(dlist) ? dlist : [];
    bookings.value = Array.isArray(blist) ? blist : [];
  } catch (err) {
    error.value =
      err?.response?.data?.error || err?.message || "Could not load orders.";
  } finally {
    loading.value = false;
  }
};

const canCancel = (status) => {
  const s = (status || "").toLowerCase();
  return s === "pending" || s === "confirmed";
};

const handleCancel = async (booking) => {
  if (!user.value?.email) {
    message.type = "error";
    message.text = "You must be logged in to cancel.";
    return;
  }
  if (!canCancel(booking.status)) return;
  confirmDialog.booking = booking;
  confirmDialog.open = true;
};

const closeConfirmDialog = () => {
  if (confirmDialog.pending) return;
  confirmDialog.open = false;
  confirmDialog.booking = null;
};

const confirmCancel = async () => {
  const booking = confirmDialog.booking;
  if (!booking) return;
  confirmDialog.pending = true;
  try {
    await cancelBooking(booking._id, user.value.email);
    await load();
    message.type = "success";
    message.text = "Booking cancelled.";
    confirmDialog.open = false;
    confirmDialog.booking = null;
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
  <section class="orders">
    <header class="orders__header">
      <h1>My Orders</h1>
      <p class="muted">Review your past orders and their status.</p>
    </header>

    <div v-if="message.text" :class="['notice', `notice--${message.type}`]">
      {{ message.text }}
    </div>

    <div v-if="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="!bookings.length" class="empty-state">
      You have no orders yet.
    </div>
    <div v-else>
      <div class="filter-tabs">
        <button
          v-for="option in statusOptions"
          :key="option.value"
          @click="selectedStatus = option.value"
          :class="[
            'filter-tab',
            { 'filter-tab--active': selectedStatus === option.value },
          ]"
        >
          {{ option.label }}
        </button>
      </div>
      <div v-if="!filteredBookings.length" class="empty-state">
        No orders with this status.
      </div>
      <div v-else class="orders-list">
        <article v-for="b in filteredBookings" :key="b._id" class="order-item">
          <RouterLink
            class="order-link"
            :to="{ name: 'order-detail', params: { id: b._id } }"
          >
            <h3>{{ formatDate(b.time) }}</h3>
            <span class="view-more">View details</span>
          </RouterLink>
          <p>{{ summarizeItems(b.items) }}</p>
          <p>Total items: {{ b.summary.count }}</p>
          <p>Total drinks price: {{ formatCurrency(b.summary.total) }}</p>
          <p :class="statusClass(b.status)">{{ formatStatus(b.status) }}</p>
          <p>Channel: {{ b.channel || "web" }}</p>
          <div>
            <button
              v-if="canCancel(b.status)"
              class="btn btn--danger"
              @click.prevent="handleCancel(b)"
            >
              Cancel booking
            </button>
          </div>
        </article>
      </div>
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
.orders__header {
  margin-bottom: 12px;
}
.order-item {
  padding: 12px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 8px;
  margin-bottom: 10px;
  display: grid;
  gap: 6px;
}
.order-link {
  text-decoration: none;
  color: inherit;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
}
.order-link h3 {
  margin: 0;
}
.view-more {
  font-size: 0.85rem;
  font-weight: 800;
  color: #18592a;
}
.muted {
  color: rgba(0, 0, 0, 0.6);
}
.error {
  color: #9d3412;
}
.status {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  padding: 6px 10px;
  border-radius: 999px;
  font-weight: 800;
  letter-spacing: 0.02em;
}
.status--pending {
  background: rgba(125, 76, 18, 0.12);
  color: #7d4c12;
}
.status--confirmed {
  background: rgba(24, 89, 42, 0.12);
  color: #18592a;
}
.status--completed {
  background: rgba(31, 78, 121, 0.12);
  color: #1f4e79;
}
.status--cancelled {
  background: rgba(157, 52, 18, 0.12);
  color: #9d3412;
}
.filter-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.filter-tab {
  padding: 8px 14px;
  border: 1px solid rgba(0, 0, 0, 0.12);
  background: rgba(255, 255, 255, 0.6);
  border-radius: 20px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s;
}
.filter-tab:hover {
  border-color: rgba(0, 0, 0, 0.24);
  background: rgba(255, 255, 255, 0.9);
}
.filter-tab--active {
  background: #18592a;
  color: white;
  border-color: #18592a;
}
.empty-state {
  padding: 20px;
  text-align: center;
  color: rgba(0, 0, 0, 0.6);
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.02);
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
