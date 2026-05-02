<script setup>
import { reactive, ref } from "vue";
import { RouterLink } from "vue-router";
import { registerUser } from "@/api";
import { sendVerificationEmail, isEmailReady } from "@/email";

const form = reactive({
  name: "",
  email: "",
  password: "",
});

const loading = ref(false);
const error = ref("");
const message = ref("");

const handleSubmit = async () => {
  error.value = "";
  message.value = "";

  if (!form.name || !form.email || !form.password) {
    error.value = "Please fill out your name, Gmail and password.";
    return;
  }

  loading.value = true;
  try {
    const res = await registerUser({
      name: form.name,
      email: form.email,
      password: form.password,
    });
    if (isEmailReady()) {
      try {
        await sendVerificationEmail(form.email, form.name, {
          token: res.token,
          verifyLink: res.verifyUrl,
          expiresAt: res.expiresAt,
        });
      } catch (mailErr) {
        console.warn("Could not send verification email", mailErr);
        error.value = "Invalid email. Try another email.";
        message.value = "";
        return;
      }
    }
    message.value =
      "Account created! Please verify your email before signing in.";
    // Stay on page; let user decide next step after verification.
  } catch (err) {
    const msg = err?.response?.data?.error || "";
    if (msg.toLowerCase().includes("invalid email")) {
      error.value = "Invalid email. Try another email.";
    } else {
      error.value = msg || "Could not create your account. Please try again.";
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
      <h1>Create your account</h1>
      <p class="lede">
        Join the ritual—save your preferences, bookings, and recommendations.
      </p>

      <form class="auth-form" @submit.prevent="handleSubmit">
        <label class="field">
          <span>Name</span>
          <input
            v-model="form.name"
            type="text"
            placeholder="E.g. Jackie Nguyen"
            autocomplete="name"
          />
        </label>

        <label class="field">
          <span>Email</span>
          <input
            v-model="form.email"
            type="email"
            placeholder="name@gmail.com"
            autocomplete="email"
          />
        </label>

        <label class="field">
          <span>Password</span>
          <input
            v-model="form.password"
            type="password"
            placeholder="Create a password"
            autocomplete="new-password"
          />
        </label>

        <div class="actions">
          <button class="btn-primary" type="submit" :disabled="loading">
            <span v-if="loading">Creating...</span>
            <span v-else>Create account</span>
          </button>
          <RouterLink to="/login" class="btn-link"
            >Already have an account? Sign in</RouterLink
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
    linear-gradient(200deg, rgba(15, 20, 36, 0.66), rgba(0, 0, 0, 0.28)),
    url("https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&w=1400&q=80")
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
  padding: 12px 14px;
  border-radius: 10px;
  font-weight: 700;
}

.banner.error {
  background: rgba(220, 90, 36, 0.12);
  color: #ffc7aa;
  border: 1px solid rgba(220, 90, 36, 0.35);
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
