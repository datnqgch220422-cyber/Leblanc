<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getBookingsByEmail } from "@/api/bookings.api";
import { getDrinks } from "@/api/drinks.api";
import { getSessionUser } from "@/composables/useSessionAuth";

const route = useRoute();
const router = useRouter();
const bookingId = route.params.id;

const user = ref(getSessionUser());
const bookings = ref([]);
const drinks = ref([]);
const loading = ref(false);
const error = ref("");

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

const formatCurrency = (v) => `${(Number(v) || 0).toLocaleString("vi-VN")} VND`;
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
          <p>Tổng món: {{ totals.count }}</p>
          <p>Tổng tiền đồ uống: {{ formatCurrency(totals.total) }}</p>
        </div>
      </section>
    </div>
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
</style>
