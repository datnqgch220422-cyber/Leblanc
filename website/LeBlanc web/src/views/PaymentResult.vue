<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute, RouterLink } from "vue-router";
import { getBookingPaymentStatus } from "@/api/bookings.api";
import { isBookingEmailReady, sendBookingEmail } from "@/email";

const route = useRoute();

const responseCode = computed(() => String(route.query.vnp_ResponseCode || ""));
const isSuccess = computed(() => responseCode.value === "00");
const message = computed(() => String(route.query.vnp_OrderInfo || ""));
const orderId = computed(() => String(route.query.vnp_TxnRef || ""));
const transactionNo = computed(() =>
  String(route.query.vnp_TransactionNo || ""),
);
const rawAmount = computed(() => Number(route.query.vnp_Amount || 0));
const amount = computed(() =>
  rawAmount.value > 0 ? rawAmount.value / 100 : 0,
);

const syncState = ref({
  paymentStatus: "pending",
  status: "pending",
  paymentMessage: "",
  booking: null,
  bookingId: "",
});
const syncing = ref(false);
const syncError = ref("");
const emailSending = ref(false);
const emailSent = ref(false);
const emailError = ref("");
const pollCount = ref(0);
const maxPollCount = 15;
let pollTimer = null;
const emailStorageKey = computed(() =>
  orderId.value ? `leblanc-booking-email-sent:${orderId.value}` : "",
);

const syncedStatus = computed(() =>
  String(syncState.value.paymentStatus || "pending"),
);
const isPaidOnServer = computed(() => syncedStatus.value === "paid");
const isFailedOnServer = computed(() => syncedStatus.value === "failed");
const displayIsSuccess = computed(
  () => isPaidOnServer.value || isSuccess.value,
);
const displayMessage = computed(() => {
  if (syncState.value.paymentMessage) return syncState.value.paymentMessage;
  if (isPaidOnServer.value) return "VNPay has confirmed the transaction.";
  if (isFailedOnServer.value) return "Transaction incomplete.";
  return (
    message.value ||
    (isSuccess.value
      ? "VNPay has confirmed the transaction."
      : "Transaction incomplete.")
  );
});

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer);
    pollTimer = null;
  }
};

const refreshNow = async () => {
  await syncPaymentStatus();
};

const maybeSendEmail = async () => {
  if (!isBookingEmailReady()) return;
  if (emailSent.value || emailSending.value) return;
  if (!isPaidOnServer.value || !syncState.value.booking) return;

  if (
    emailStorageKey.value &&
    localStorage.getItem(emailStorageKey.value) === "1"
  ) {
    emailSent.value = true;
    return;
  }

  emailSending.value = true;
  emailError.value = "";
  try {
    await sendBookingEmail(syncState.value.booking, {
      bookingId:
        syncState.value.booking.bookingId ||
        syncState.value.bookingId ||
        orderId.value,
    });
    emailSent.value = true;
    if (emailStorageKey.value) {
      localStorage.setItem(emailStorageKey.value, "1");
    }
  } catch (err) {
    emailError.value =
      err?.message || "Unable to send confirmation email after payment";
  } finally {
    emailSending.value = false;
  }
};

const syncPaymentStatus = async () => {
  const paymentOrderId = orderId.value.trim();
  if (!paymentOrderId) return;

  syncing.value = true;
  syncError.value = "";
  try {
    const res = await getBookingPaymentStatus(paymentOrderId);
    syncState.value = {
      paymentStatus: String(res?.paymentStatus || "pending"),
      status: String(res?.status || "pending"),
      paymentMessage: String(res?.paymentMessage || ""),
      booking: res?.booking || null,
      bookingId: String(res?.bookingId || paymentOrderId),
    };

    if (syncState.value.paymentStatus !== "pending") {
      stopPolling();
      pollCount.value = 0;
    } else {
      pollCount.value += 1;
      if (pollCount.value >= maxPollCount) {
        stopPolling();
      }
    }

    if (syncState.value.paymentStatus === "paid") {
      await maybeSendEmail();
    }
  } catch (err) {
    syncError.value =
      err?.response?.data?.error ||
      err?.message ||
      "Unable to sync payment status";
  } finally {
    syncing.value = false;
  }
};

const formatCurrency = (value) =>
  `${(Number(value) || 0).toLocaleString("en-US")} VND`;

onMounted(async () => {
  await syncPaymentStatus();
  if (orderId.value) {
    pollTimer = setInterval(syncPaymentStatus, 2000);
  }
});

onBeforeUnmount(() => {
  stopPolling();
});
</script>

<template>
  <section class="payment-result">
    <h1>
      {{ displayIsSuccess ? "Payment successful" : "Payment not successful" }}
    </h1>
    <p
      class="status"
      :class="{ ok: displayIsSuccess, fail: !displayIsSuccess }"
    >
      {{
        displayMessage ||
        (displayIsSuccess
          ? "VNPay has confirmed the transaction."
          : "Transaction incomplete.")
      }}
    </p>

    <p v-if="syncing" class="syncing">Syncing payment status...</p>
    <p v-else-if="syncError" class="syncing error-text">{{ syncError }}</p>
    <p v-else-if="syncState.paymentStatus === 'pending'" class="syncing">
      The page will auto-check the status while waiting for VNPay confirmation.
    </p>
    <p v-if="emailSending" class="syncing">Sending confirmation email...</p>
    <p v-else-if="emailSent" class="syncing success-text">
      Confirmation email has been sent after successful payment.
    </p>
    <p v-else-if="emailError" class="syncing error-text">{{ emailError }}</p>

    <div class="info">
      <p><strong>Payment order ID:</strong> {{ orderId || "—" }}</p>
      <p><strong>Transaction ID:</strong> {{ transactionNo || "—" }}</p>
      <p>
        <strong>Amount:</strong> {{ amount ? formatCurrency(amount) : "—" }}
      </p>
      <p><strong>Response code:</strong> {{ responseCode || "—" }}</p>
      <p><strong>System status:</strong> {{ syncState.paymentStatus }}</p>
    </div>

    <div class="actions">
      <button class="btn" type="button" @click="refreshNow">Refresh now</button>
      <RouterLink class="btn" to="/orders">View my orders</RouterLink>
      <RouterLink class="btn ghost" to="/booking">Back to booking</RouterLink>
    </div>
  </section>
</template>

<style scoped>
.payment-result {
  display: grid;
  gap: 12px;
  padding: 14px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.82);
}
h1 {
  margin: 0;
}
.status {
  margin: 0;
  padding: 10px 12px;
  border-radius: 10px;
}
.status.ok {
  background: #e6f5ed;
  color: #156f3d;
}
.status.fail {
  background: #fde9ea;
  color: #b00020;
}
.syncing {
  margin: 0;
  color: rgba(0, 0, 0, 0.7);
}
.error-text {
  color: #9d3412;
}
.success-text {
  color: #156f3d;
}
.info {
  display: grid;
  gap: 6px;
  padding: 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.04);
}
.info p {
  margin: 0;
}
.actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
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
