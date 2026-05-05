<script setup>
import { useBookingFlow } from "@/composables/useBookingFlow";
import { useThemeState } from "@/composables/useThemeState";

const { isNight } = useThemeState();

const {
  form,
  formDate,
  formClock,
  mood,
  caffeinePref,
  tempPref,
  sweetness,
  nightType,
  nightBase,
  reco,
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
} = useBookingFlow(isNight);
</script>

<template>
  <section class="booking">
    <div class="panel form">
      <p class="eyebrow">Reservation</p>
      <h1>Book your table</h1>
      <p class="lede">Giữ chỗ và thêm đồ uống trước nếu bạn muốn.</p>

      <form class="form-fields" @submit.prevent="book">
        <label>
          Name
          <input v-model="form.name" placeholder="Nguyễn Văn A" />
        </label>
        <label>
          Phone
          <input v-model="form.phone" placeholder="0123 456 789" />
        </label>
        <label>
          Email
          <input
            v-model="form.email"
            type="email"
            placeholder="name@gmail.com"
          />
        </label>
        <label>
          Arrival Date
          <input v-model="formDate" type="date" />
        </label>
        <label>
          Arrival Time
          <input v-model="formClock" type="time" step="900" />
        </label>
        <label>
          Guests
          <input v-model.number="form.guests" type="number" min="1" max="10" />
        </label>

        <div class="selected" v-if="selectedItems.length">
          <p class="mini-title">Pre-order drinks ({{ totalItems }} items)</p>
          <div class="chip-list">
            <div v-for="item in selectedItems" :key="item.drinkId" class="chip">
              <span>{{ item.drink?.name || "Drink" }}</span>
              <div class="qty">
                <button type="button" @click="updateQty(item.drinkId, -1)">
                  -
                </button>
                <span>{{ item.qty }}</span>
                <button type="button" @click="updateQty(item.drinkId, 1)">
                  +
                </button>
              </div>
            </div>
          </div>
        </div>

        <button type="submit" :disabled="!canSubmit || bookingLoading">
          <span v-if="bookingLoading">Processing...</span>
          <span v-else>Book table{{ totalItems ? " & drinks" : "" }}</span>
        </button>
        <p v-if="bookingOk" class="status success">
          Đặt bàn đã được tạo. Đang chuyển sang thanh toán VNPay...
        </p>
        <p v-if="bookingError && !bookingOk" class="status error">
          {{ bookingError }}
        </p>
      </form>
    </div>

    <div class="panel reco">
      <div>
        <p class="eyebrow">My Booking</p>
        <h2>Added drinks</h2>
      </div>

      <div v-if="selectedItems.length" class="selected-list">
        <p class="mini-title">Added {{ totalItems }} drinks</p>
        <div class="chip-list">
          <div v-for="item in selectedItems" :key="item.drinkId" class="chip">
            <div>
              <div class="name">{{ item.drink?.name || "Drink" }}</div>
              <div class="meta">
                {{
                  item.drink?.price
                    ? item.drink.price.toLocaleString("vi-VN") + " VND"
                    : "—"
                }}
              </div>
            </div>

            <div class="qty">
              <button type="button" @click="updateQty(item.drinkId, -1)">
                -
              </button>
              <span>{{ item.qty }}</span>
              <button type="button" @click="updateQty(item.drinkId, 1)">
                +
              </button>
              <button
                type="button"
                class="mini"
                @click="updateQty(item.drinkId, -item.qty)"
              >
                Xóa
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="status">
        No drinks added yet. <br />
        You can add some from the menu before booking.
      </div>
    </div>
  </section>
</template>

<style scoped>
.booking {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  gap: 24px;
  color: var(--ink);
}

.panel {
  background: var(--paper);
  padding: clamp(24px, 4vw, 32px);
  border-radius: 16px;
  display: grid;
  gap: 14px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.12);
  align-content: start;
}

.selected-list {
  display: grid;
  gap: 12px;
  overflow-y: auto;
  max-height: 60vh;
}

.form-fields {
  display: grid;
  gap: 12px;
}

label {
  display: grid;
  gap: 6px;
  font-weight: 700;
}

input,
select {
  border: 1px solid var(--cream-strong);
  padding: 12px 14px;
  border-radius: 10px;
  font-family: inherit;
  background: var(--paper);
  color: var(--ink);
}

button {
  border: 1px solid var(--dark);
  background: var(--dark);
  color: #fff;
  padding: 14px 16px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 700;
  font-size: 1rem;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.ghost {
  background: transparent;
  color: var(--ink);
  border: 1px solid rgba(0, 0, 0, 0.12);
}

.mini {
  padding: 8px 12px;
  font-size: 0.95rem;
}

.eyebrow {
  margin: 0;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  font-size: 0.8rem;
}

.lede {
  margin: 0 0 4px;
  color: rgba(0, 0, 0, 0.7);
}

.selected {
  border: 1px dashed rgba(0, 0, 0, 0.12);
  padding: 10px;
  border-radius: 12px;
}

.chip-list {
  display: grid;
  gap: 8px;
  max-height: 50vh;
  overflow-y: auto;
}

.chip {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
}

.chip .qty {
  display: flex;
  align-items: center;
  gap: 8px;
}

.chip .qty button {
  padding: 4px 10px;
  border-radius: 8px;
}

.controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 10px;
}

.reco-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.reco-list {
  display: grid;
  gap: 12px;
}

.card {
  padding: 12px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.06);
  display: grid;
  gap: 6px;
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.name {
  margin: 0;
  font-weight: 800;
}

.meta {
  margin: 0;
  color: rgba(0, 0, 0, 0.7);
}

.score {
  margin-left: 8px;
  font-weight: 700;
}

.desc {
  margin: 0;
  color: rgba(0, 0, 0, 0.7);
}

.status {
  margin: 0;
  padding: 10px 12px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.04);
}

.status.success {
  background: #e6f5ed;
  color: #156f3d;
}

.status.error {
  background: #fde9ea;
  color: #b00020;
}

.mini-title {
  margin: 0 0 6px;
  font-weight: 800;
}

:global(.theme-night) .panel {
  background: rgba(15, 20, 36, 0.7);
  border-color: rgba(255, 255, 255, 0.08);
  color: #f6efe6;
}

:global(.theme-night) input,
:global(.theme-night) select {
  background: rgba(255, 255, 255, 0.06);
  color: #f6efe6;
  border-color: rgba(255, 255, 255, 0.2);
}

:global(.theme-night) label,
:global(.theme-night) .kicker,
:global(.theme-night) .hero-title,
:global(.theme-night) .hero-copy,
:global(.theme-night) .booking {
  color: #f6efe6;
}

:global(.theme-night) .lede,
:global(.theme-night) .meta,
:global(.theme-night) .desc {
  color: rgba(245, 241, 232, 0.82);
}

:global(.theme-night) .lede {
  color: rgba(245, 241, 232, 0.82);
}

:global(.theme-night) .card .name {
  color: #f6efe6;
}

:global(.theme-night) .chip {
  background: rgba(255, 255, 255, 0.08);
  color: #f6efe6;
}

:global(.theme-night) .status {
  background: rgba(255, 255, 255, 0.08);
  color: #f6efe6;
}

:global(.theme-night) .card {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.06);
}

:global(.theme-night) .status {
  background: rgba(255, 255, 255, 0.06);
  color: #f6efe6;
}

:global(.theme-night) .status.success {
  background: rgba(62, 146, 98, 0.18);
  color: #a7f3c7;
}

:global(.theme-night) .status.error {
  background: rgba(255, 94, 94, 0.18);
  color: #ffc7c7;
}

:global(.theme-night) .mini-title {
  color: #f6efe6;
}

:global(.theme-night) .chip {
  background: rgba(255, 255, 255, 0.08);
}

:global(.theme-night) .chip .name {
  color: #f6efe6;
}

:global(.theme-night) .chip .meta {
  color: rgba(245, 241, 232, 0.82);
}

:global(.theme-night) .selected-list {
  color: #f6efe6;
}
</style>
