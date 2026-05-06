<script setup>
import { computed } from "vue";
import { useRoute, RouterLink } from "vue-router";

const route = useRoute();
const bookingId = computed(() => route.query.id || "");
const bookingTotal = computed(() => Number(route.query.total || 0));
const bookingItems = computed(() => Number(route.query.items || 0));

const formatCurrency = (value) =>
  `${(Number(value) || 0).toLocaleString("en-US")} VND`;
</script>

<template>
  <section class="confirm">
    <h1>Booking created</h1>
    <p>
      Thank you for your booking and selecting accompanying drinks.
      <span v-if="bookingId"
        >Booking ID: <strong>{{ bookingId }}</strong
        >.</span
      >
    </p>
    <div class="summary">
      <p>
        Total items: <strong>{{ bookingItems }}</strong>
      </p>
      <p>
        Total drinks price: <strong>{{ formatCurrency(bookingTotal) }}</strong>
      </p>
    </div>
    <div class="actions">
      <RouterLink class="btn" to="/orders">View My Orders</RouterLink>
      <RouterLink class="btn ghost" to="/menu">Back to Menu</RouterLink>
    </div>
  </section>
</template>

<style scoped>
.confirm {
  display: grid;
  gap: 12px;
  padding: 14px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.8);
}
.confirm h1 {
  margin: 0;
}
.actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
.summary {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  padding: 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
}
.btn {
  text-decoration: none;
  padding: 10px 12px;
  border-radius: 10px;
  background: #0f1424;
  color: #fff;
  font-weight: 800;
}
.btn.ghost {
  background: #fff;
  color: #0f1424;
  border: 1px solid rgba(0, 0, 0, 0.16);
}
</style>
