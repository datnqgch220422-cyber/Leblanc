<script setup>
import { computed, reactive, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import {
  getSession,
  getSessionUser,
  persistSession,
} from "@/composables/useSessionAuth";

const router = useRouter();
const user = ref(getSessionUser());
const form = reactive({
  name: user.value?.name ?? "",
  email: user.value?.email ?? "",
  password: "",
  confirmPassword: "",
});
const isSaving = ref(false);
const successMessage = ref("");
const errorMessage = ref("");
const showPassword = ref(false);
const showConfirm = ref(false);

const hasChanges = computed(() => {
  const nameChanged = form.name.trim() !== (user.value?.name ?? "");
  const emailChanged = form.email.trim() !== (user.value?.email ?? "");
  const passwordChanged =
    form.password.trim() !== "" || form.confirmPassword.trim() !== "";
  return nameChanged || emailChanged || passwordChanged;
});

const isSaveDisabled = computed(() => {
  return !hasChanges.value || isSaving.value;
});

const saveAccount = () => {
  if (!user.value) {
    return;
  }

  const trimmedName = form.name.trim();
  const trimmedEmail = form.email.trim();

  successMessage.value = "";
  errorMessage.value = "";

  if (!trimmedName) {
    errorMessage.value = "Username cannot be empty.";
    return;
  }

  if (!trimmedEmail) {
    errorMessage.value = "Email cannot be empty.";
    return;
  }

  const passwordTrimmed = form.password.trim();
  const confirmTrimmed = form.confirmPassword.trim();

  if (passwordTrimmed && !confirmTrimmed) {
    errorMessage.value = "Please confirm your new password.";
    return;
  }

  if (!passwordTrimmed && confirmTrimmed) {
    errorMessage.value = "Please enter your new password.";
    return;
  }

  if (passwordTrimmed) {
    if (passwordTrimmed.length < 6) {
      errorMessage.value = "New password must be at least 6 characters.";
      return;
    }

    if (passwordTrimmed !== confirmTrimmed) {
      errorMessage.value = "Password confirmation does not match.";
      return;
    }
  }

  isSaving.value = true;

  const updatedUser = {
    ...user.value,
    name: trimmedName,
    email: trimmedEmail,
    ...(passwordTrimmed ? { password: passwordTrimmed } : {}),
  };

  persistSession({
    user: updatedUser,
    token: getSession()?.token ?? "",
  });

  user.value = updatedUser;
  form.password = "";
  form.confirmPassword = "";
  successMessage.value = "Account information has been updated.";
  isSaving.value = false;
};
</script>

<template>
  <section class="account">
    <header class="account__header">
      <div>
        <p class="eyebrow">Le'Blanc</p>
        <h1>Your account</h1>
      </div>
      <RouterLink to="/" class="btn-ghost">Back to home</RouterLink>
    </header>

    <div class="card-block editor" v-if="user">
      <div class="editor__header">
        <div>
          <p class="eyebrow">Le'Blanc</p>
          <h2>Edit account</h2>
        </div>
      </div>

      <form class="editor__form" @submit.prevent="saveAccount">
        <label class="field">
          <span>Username</span>
          <input v-model="form.name" type="text" autocomplete="username" />
        </label>

        <label class="field">
          <span>Email</span>
          <input v-model="form.email" type="email" autocomplete="email" />
        </label>

        <label class="field">
          <div class="field__header">
            <span>New Password</span>
            <button
              class="toggle-pwd"
              type="button"
              @click="showPassword = !showPassword"
              :title="showPassword ? 'Hide password' : 'Show password'"
            >
              {{ showPassword ? "Hide" : "Show" }}
            </button>
          </div>
          <input
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="new-password"
            placeholder="Enter new password"
          />
        </label>

        <label class="field">
          <div class="field__header">
            <span>Confirm Password</span>
            <button
              class="toggle-pwd"
              type="button"
              @click="showConfirm = !showConfirm"
              :title="showConfirm ? 'Hide password' : 'Show password'"
            >
              {{ showConfirm ? "Hide" : "Show" }}
            </button>
          </div>
          <input
            v-model="form.confirmPassword"
            :type="showConfirm ? 'text' : 'password'"
            autocomplete="new-password"
            placeholder="Confirm new password"
          />
        </label>

        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
        <p v-if="successMessage" class="success">{{ successMessage }}</p>

        <div class="editor__actions">
          <button class="btn" type="submit" :disabled="isSaveDisabled">
            {{ isSaving ? "Saving..." : "Save changes" }}
          </button>
        </div>
      </form>
    </div>
  </section>
</template>

<style scoped>
.account {
  display: grid;
  gap: 18px;
}

.account__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.eyebrow {
  margin: 0;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  font-size: 0.82rem;
  color: var(--tan);
}

h1 {
  margin: 2px 0;
  font-size: clamp(1.8rem, 3vw, 2.4rem);
}

.lede {
  margin: 0;
  color: rgba(0, 0, 0, 0.7);
}

.card {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  border-radius: 14px;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.9),
    rgba(246, 239, 230, 0.8)
  );
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.08);
}

.card.empty {
  grid-template-columns: 1fr;
}

.avatar {
  height: 56px;
  width: 56px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: linear-gradient(145deg, #b88443, #e1c58d);
  font-weight: 900;
  color: #0b0b0b;
  letter-spacing: 0.02em;
}

.meta {
  display: grid;
  gap: 4px;
}

.name {
  margin: 0;
  font-weight: 900;
  font-size: 1.2rem;
}

.email,
.id {
  margin: 0;
  color: rgba(0, 0, 0, 0.7);
  word-break: break-all;
}

.actions {
  display: flex;
  gap: 10px;
}

.btn-link {
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.card-block {
  display: grid;
  gap: 14px;
  padding: 18px 20px;
  border-radius: 14px;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.9),
    rgba(246, 239, 230, 0.8)
  );
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.08);
}

.editor {
  gap: 18px;
}

.editor__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.editor__form {
  display: grid;
  gap: 14px;
}

.field {
  display: grid;
  gap: 8px;
  font-weight: 700;
}

.field span {
  color: rgba(0, 0, 0, 0.8);
}

.field input {
  width: 100%;
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.12);
  background: rgba(255, 255, 255, 0.75);
  color: #0b0b0b;
  font: inherit;
  outline: none;
}

.field input:focus {
  border-color: rgba(184, 132, 67, 0.7);
  box-shadow: 0 0 0 3px rgba(184, 132, 67, 0.16);
}

.field__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.toggle-pwd {
  padding: 0;
  border: none;
  background: none;
  color: #b88443;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  transition: color 0.18s ease;
  text-decoration: underline;
}

.toggle-pwd:hover {
  color: #c8954f;
}

.editor__actions {
  display: flex;
  justify-content: flex-end;
}

.success {
  margin: 0;
  color: #166534;
  font-weight: 700;
}

.empty-state,
.muted,
.error {
  margin: 0;
}

.muted {
  color: rgba(0, 0, 0, 0.68);
}

.status {
  margin: 0;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #7d4c12;
}

.error {
  color: #9d3412;
}

.btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn {
  padding: 10px 14px;
  border-radius: 10px;
  border: none;
  background: #0f1424;
  color: #f6efe6;
  font-weight: 800;
  cursor: pointer;
}

.btn-ghost {
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid rgba(0, 0, 0, 0.2);
  color: inherit;
  text-decoration: none;
  font-weight: 800;
}

@media (max-width: 720px) {
  .card {
    grid-template-columns: auto 1fr;
    grid-template-areas:
      "avatar meta"
      "actions actions";
  }
  .actions {
    grid-area: actions;
  }
}
</style>
