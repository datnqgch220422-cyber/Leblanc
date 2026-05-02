<script setup>
import { reactive, ref } from "vue";
import { useRouter, RouterLink } from "vue-router";
import { loginUser } from "@/api";
import {
  ADMIN_HOME_PATH,
  isAdminUser,
  persistSession,
} from "@/composables/useSessionAuth";

const router = useRouter();
const form = reactive({
  name: "",
  password: "",
});

const loading = ref(false);
const error = ref("");
const message = ref("");

const getSafeRedirect = () => {
  const redirect = router.currentRoute.value.query.redirect;
  return typeof redirect === "string" && redirect.startsWith("/")
    ? redirect
    : "";
};

const handleSubmit = async () => {
  error.value = "";
  message.value = "";

  if (!form.name || !form.password) {
    error.value = "Please provide both your name and password.";
    return;
  }

  loading.value = true;
  try {
    const res = await loginUser({
      name: form.name,
      password: form.password,
    });
    persistSession({
      user: res.user,
      token: res.token || "",
    });
    const fallbackTarget = isAdminUser(res.user) ? ADMIN_HOME_PATH : "/";
    const target = getSafeRedirect() || fallbackTarget;
    message.value = `Welcome back, ${res.user.name}!`;
    setTimeout(() => router.push(target), 600);
  } catch (err) {
    // Fallback message first to guarantee something shows.
    error.value = "Unable to sign you in. Please try again.";
    const status = err?.response?.status;
    const raw = err?.response?.data?.error || "";
    const lower = raw.toLowerCase();
    if (
      lower.includes("user") ||
      lower.includes("name") ||
      lower.includes("found")
    ) {
      error.value = "Wrong User Name. Please correct it.";
    } else if (lower.includes("password")) {
      error.value = "Wrong Password. Please correct it.";
    } else if (status === 401) {
      // If only status is known, assume wrong password.
      error.value = "Wrong Password. Please correct it.";
    }
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <section class="auth-page">
    <div class="auth-visual" aria-hidden="true"></div>
    <div class="auth-card">
      <p class="eyebrow">Le'Blanc</p>
      <h1>Welcome back</h1>
      <p class="lede">Sign in to keep your rituals and bookings in sync.</p>

      <form class="auth-form" @submit.prevent="handleSubmit">
        <label class="field">
          <span>Name</span>
          <input
            v-model="form.name"
            type="text"
            placeholder="E.g. Jackie Nguyen"
            autocomplete="username"
          />
        </label>

        <label class="field">
          <span>Password</span>
          <input
            v-model="form.password"
            type="password"
            placeholder="Enter your password"
            autocomplete="current-password"
          />
        </label>

        <div class="actions">
          <button class="btn-primary" type="submit" :disabled="loading">
            <span v-if="loading">Signing in...</span>
            <span v-else>Log in</span>
          </button>
          <RouterLink to="/register" class="btn-link"
            >Create account</RouterLink
          >
        </div>
      </form>

      <div v-if="error" class="banner error">{{ error }}</div>
      <div v-if="message" class="banner success">{{ message }}</div>
    </div>
  </section>
</template>

<style scoped>
.auth-page {
  width: 100%;
  min-height: 78vh;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
  align-items: stretch;
  background: linear-gradient(135deg, #0f1424 0%, #1d1a16 100%);
  color: #f6efe6;
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 28px 70px rgba(0, 0, 0, 0.28);
}

.auth-visual {
  background:
    linear-gradient(160deg, rgba(15, 20, 36, 0.7), rgba(0, 0, 0, 0.2)),
    url("https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1400&q=80")
      center/cover no-repeat;
}

.auth-card {
  padding: clamp(24px, 5vw, 48px);
  display: grid;
  gap: 16px;
  background: linear-gradient(
    145deg,
    rgba(246, 239, 230, 0.06),
    rgba(246, 239, 230, 0.12)
  );
  backdrop-filter: blur(6px);
}

.eyebrow {
  margin: 0;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  font-size: 0.78rem;
  font-family: "Georgia", "Times New Roman", serif;
  color: #e9d7b6;
}

h1 {
  margin: 0;
  font-size: clamp(1.9rem, 3vw, 2.4rem);
  line-height: 1.2;
}

.lede {
  margin: 0 0 8px;
  color: #e5dfd6;
}

.auth-form {
  display: grid;
  gap: 14px;
  margin-top: 4px;
}

.field {
  display: grid;
  gap: 6px;
  font-weight: 700;
  color: #f6efe6;
}

.field span {
  font-size: 0.92rem;
  letter-spacing: 0.02em;
}

.field input {
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.24);
  background: rgba(255, 255, 255, 0.06);
  color: #f6efe6;
  font-size: 1rem;
}

.field input::placeholder {
  color: rgba(246, 239, 230, 0.6);
}

.actions {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
  margin-top: 6px;
}

.btn-primary {
  padding: 12px 20px;
  border-radius: 999px;
  border: none;
  background: #b88443;
  color: #0b0b0b;
  font-weight: 800;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition:
    transform 0.18s ease,
    box-shadow 0.18s ease,
    background 0.18s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.22);
  background: #c8954f;
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-link {
  color: #f6efe6;
  text-decoration: underline;
  font-weight: 700;
}

.banner {
  padding: 0;
  font-weight: 800;
  font-size: 1.05rem;
  color: #ffc7aa;
}

.banner.error {
  background: transparent;
  border: none;
}

.banner.success {
  background: rgba(105, 205, 145, 0.14);
  color: #d0ffe8;
  border: 1px solid rgba(105, 205, 145, 0.35);
}

@media (max-width: 720px) {
  .auth-page {
    grid-template-columns: 1fr;
  }

  .auth-visual {
    min-height: 180px;
  }
}
</style>
