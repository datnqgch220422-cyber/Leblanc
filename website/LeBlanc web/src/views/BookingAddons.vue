<script setup>
import { computed, onMounted, ref } from "vue";
import { getSessionUser } from "@/composables/useSessionAuth";
import { useBookingAddons } from "@/composables/useBookingAddons";
import { getDrinks } from "@/api/drinks.api";

const { bookingAddons, load, update, remove, loading, doCheckout } =
  useBookingAddons();

const drinks = ref([]);
const name = ref("");
const phone = ref("");
const email = ref("");
const paymentMethod = ref("vnpay");
const checkoutLoading = ref(false);
const checkoutError = ref("");
const toast = ref("");

const formatCurrency = (value) =>
  `${(Number(value) || 0).toLocaleString("vi-VN")} VND`;

const drinkMap = computed(() => {
  const map = new Map();
  for (const d of drinks.value) map.set(d._id, d);
  return map;
});

const detailItems = computed(() =>
  (bookingAddons.value?.items || []).map((item) => {
    const drink = drinkMap.value.get(item.drinkId) || {};
    const price = Number(drink.price) || 0;
    const qty = Number(item.qty) || 0;
    return {
      ...item,
      name: drink.name || item.drinkId,
      price,
      lineTotal: price * qty,
    };
  }),
);

const subTotal = computed(() =>
  detailItems.value.reduce((sum, item) => sum + item.lineTotal, 0),
);
const totalItems = computed(() =>
  detailItems.value.reduce((sum, item) => sum + (Number(item.qty) || 0), 0),
);

const showToast = (message) => {
  toast.value = message;
  setTimeout(() => {
    toast.value = "";
  }, 1800);
};

const changeQty = async (drinkId, qty) => {
  if (qty <= 0) {
    await remove(drinkId);
    showToast("Đã bỏ đồ uống khỏi booking");
    return;
  }
  await update(drinkId, qty);
};

const removeItem = async (drinkId) => {
  await remove(drinkId);
  showToast("Đã bỏ đồ uống khỏi booking");
};

const onCreateBooking = async () => {
  checkoutError.value = "";
  if (!name.value || !phone.value || !email.value) {
    checkoutError.value = "Vui lòng nhập đủ tên, điện thoại và email";
    return;
  }

  checkoutLoading.value = true;
  try {
    const res = await doCheckout({
      name: name.value,
      phone: phone.value,
      email: email.value,
      channel: "web-booking",
      paymentMethod: paymentMethod.value,
    });
    const payUrl = String(res?.payUrl || "").trim();
    if (!payUrl) {
      throw new Error("Không lấy được link thanh toán VNPay");
    }

    window.location.href = payUrl;
  } catch (err) {
    checkoutError.value =
      err?.response?.data?.error || err?.message || "Không thể tạo booking";
  } finally {
    checkoutLoading.value = false;
  }
};

onMounted(async () => {
  const user = getSessionUser();
  if (user?.name) name.value = user.name;
  if (user?.email) email.value = user.email;
  drinks.value = (await getDrinks()) || [];
  await load();
});
</script>

<template>
  <section class="booking-addons">
    <header class="head">
      <h2>Đồ uống đi kèm đặt bàn</h2>
      <p class="muted">
        Chọn trước đồ uống để quầy chuẩn bị đúng theo booking của bạn.
      </p>
    </header>

    <p v-if="toast" class="toast">{{ toast }}</p>

    <div v-if="loading">Đang tải...</div>
    <div v-else-if="detailItems.length === 0" class="empty">
      Bạn chưa chọn đồ uống kèm nào.
    </div>

    <ul v-else class="item-list">
      <li v-for="item in detailItems" :key="item.drinkId" class="item">
        <div>
          <p class="name">{{ item.name }}</p>
          <p class="muted">{{ formatCurrency(item.price) }} x {{ item.qty }}</p>
        </div>
        <div class="qty">
          <button @click="changeQty(item.drinkId, item.qty - 1)">-</button>
          <span>{{ item.qty }}</span>
          <button @click="changeQty(item.drinkId, item.qty + 1)">+</button>
        </div>
        <div class="line-total">{{ formatCurrency(item.lineTotal) }}</div>
        <button class="remove" @click="removeItem(item.drinkId)">Bỏ</button>
      </li>
    </ul>

    <div class="summary" v-if="detailItems.length">
      <p>
        Tổng món: <strong>{{ totalItems }}</strong>
      </p>
      <p>
        Tạm tính đồ uống: <strong>{{ formatCurrency(subTotal) }}</strong>
      </p>
    </div>

    <div class="checkout" v-if="detailItems.length">
      <h3>Thông tin booking</h3>
      <label>Tên <input v-model="name" /></label>
      <label>Phone <input v-model="phone" /></label>
      <label>Email <input v-model="email" /></label>

      <label>
        Phương thức
        <select v-model="paymentMethod">
          <option value="vnpay">VNPay Sandbox</option>
        </select>
      </label>
      <p class="muted">
        Sau khi bấm xác nhận, bạn sẽ được chuyển tới cổng thanh toán VNPay
        sandbox.
      </p>

      <button
        class="btn-primary"
        type="button"
        :disabled="checkoutLoading"
        @click="onCreateBooking"
      >
        {{
          checkoutLoading
            ? "Đang tạo link thanh toán..."
            : "Thanh toán với VNPay"
        }}
      </button>

      <p v-if="checkoutError" class="error">{{ checkoutError }}</p>
    </div>
  </section>
</template>

<style scoped>
.booking-addons {
  display: grid;
  gap: 14px;
}
.head h2 {
  margin: 0;
}
.muted {
  margin: 0;
  color: rgba(0, 0, 0, 0.66);
}
.toast {
  margin: 0;
  padding: 10px 12px;
  background: #edf9ee;
  border: 1px solid #b6e0ba;
  color: #17622a;
  border-radius: 10px;
}
.empty {
  padding: 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
}
.item-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}
.item {
  display: grid;
  grid-template-columns: 1fr auto auto auto;
  gap: 12px;
  align-items: center;
  padding: 12px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.75);
}
.name {
  margin: 0;
  font-weight: 800;
}
.qty {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.qty button {
  width: 28px;
  height: 28px;
}
.line-total {
  font-weight: 800;
}
.remove {
  border: 1px solid rgba(0, 0, 0, 0.2);
  background: transparent;
  border-radius: 8px;
  padding: 6px 10px;
}
.summary {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
}
.checkout {
  display: grid;
  gap: 10px;
  padding: 14px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.75);
}
.checkout label {
  display: grid;
  gap: 6px;
  font-weight: 700;
}
.checkout input,
.checkout select {
  padding: 10px;
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.16);
}
.btn-primary {
  cursor: pointer;
  border: none;
  border-radius: 10px;
  padding: 10px 12px;
  font-weight: 800;
  background: #0f1424;
  color: #fff;
}
.error {
  margin: 0;
  color: #9d3412;
}

@media (max-width: 760px) {
  .item {
    grid-template-columns: 1fr auto;
  }
}
</style>
