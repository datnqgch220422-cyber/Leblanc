<script setup>
import { ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import { clearSessionUser, getSessionUser } from "@/composables/useSessionAuth";

const router = useRouter();
const user = ref(getSessionUser());

const logout = () => {
  clearSessionUser();
  router.push("/login");
};
</script>

<template>
  <section class="account">
    <header class="account__header">
      <div>
        <p class="eyebrow">Le'Blanc</p>
        <h1>Your account</h1>
        <p class="lede">
          View who is signed in and jump back to the experience.
        </p>
      </div>
      <RouterLink to="/" class="btn-ghost">Back to home</RouterLink>
    </header>

    <div class="card" v-if="user">
      <div class="avatar">{{ (user.name?.[0] || "A").toUpperCase() }}</div>
      <div class="meta">
        <p class="name">{{ user.name }}</p>
        <p class="email">{{ user.email }}</p>
      </div>
      <div class="actions">
        <button class="btn" type="button" @click="logout">Log out</button>
      </div>
    </div>

    <div class="card empty" v-else>
      <p>You are not signed in. Redirecting to login...</p>
    </div>

    <div class="history card-block" v-if="user">
      <div class="history__header">
        <div>
          <p class="eyebrow">Le'Blanc</p>
          <h2>Booking Hub</h2>
        </div>
        <div class="actions">
          <RouterLink to="/booking" class="btn btn-link">Booking</RouterLink>
          <RouterLink to="/orders" class="btn btn-link">My Orders</RouterLink>
        </div>
      </div>
      <p class="muted">
        Manage bookings and pre-order accompanying drinks directly on the
        Booking page, then track status in "My Orders".
      </p>
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

.history__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.history-list {
  display: grid;
  gap: 12px;
}

.history-item {
  display: grid;
  gap: 6px;
  padding: 14px 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.6);
  border: 1px solid rgba(0, 0, 0, 0.06);
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
