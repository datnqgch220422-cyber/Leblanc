<script setup>
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { verifyToken } from "@/api";
import { persistSessionUser } from "@/composables/useSessionAuth";
import darkLogo from "@/assets/dark-logo.png";

const route = useRoute();
const router = useRouter();
const token = ref(route.query.token || "");
const loading = ref(false);
const message = ref("");
const error = ref("");

const handleVerify = async () => {
  error.value = "";
  message.value = "";
  if (!token.value) {
    error.value = "Missing token. Please use the link from your email.";
    return;
  }
  loading.value = true;
  try {
    const res = await verifyToken({ token: token.value });
    const user = res?.data?.user || res?.user;
    if (user) {
      persistSessionUser(user);
    }
    message.value = "Email verified! Redirecting...";
    setTimeout(() => router.push("/"), 800);
  } catch (err) {
    error.value =
      err?.response?.data?.error || "Verification failed. Please try again.";
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <section class="verify">
    <img :src="darkLogo" alt="Le'Blanc" class="logo" />
    <button class="btn" type="button" :disabled="loading" @click="handleVerify">
      <span v-if="loading">Verifying...</span>
      <span v-else>Verify your email</span>
    </button>
    <p v-if="message" class="status success">{{ message }}</p>
    <p v-if="error" class="status error">{{ error }}</p>
  </section>
</template>

<style scoped>
.verify {
  min-height: 60vh;
  display: grid;
  place-items: center;
  gap: 18px;
  padding: 32px 16px;
}

.logo {
  height: 56px;
  width: auto;
}

.btn {
  padding: 12px 18px;
  border-radius: 999px;
  border: 1px solid var(--dark);
  background: var(--dark);
  color: #fff;
  font-weight: 800;
  cursor: pointer;
}

.btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.status {
  margin: 0;
  font-weight: 700;
}

.status.success {
  color: #156f3d;
}

.status.error {
  color: #b00020;
}
</style>
