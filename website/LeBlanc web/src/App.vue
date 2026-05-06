<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink, RouterView, useRoute, useRouter } from "vue-router";
import {
  ADMIN_HOME_PATH,
  clearSessionUser,
  getSession,
  getSessionUser,
  isAdminUser,
} from "@/composables/useSessionAuth";
import { provideThemeState } from "@/composables/useThemeState";
import { USER_STORAGE_KEY } from "@/services/session.service";
import darkLogo from "@/assets/dark-logo.png";
import brightLogo from "@/assets/bright-logo.png";

const { theme, toggleTheme } = provideThemeState();
const router = useRouter();
const route = useRoute();
const session = ref(getSession());
const user = ref(getSessionUser());
const isAuthed = computed(() => !!user.value);
const isAdmin = computed(() => isAdminUser(user.value));
const isAdminOnlyShell = computed(
  () => isAdmin.value && route.path.startsWith("/admin"),
);
const isPlainLayout = computed(() => route.meta?.layout === "plain");
const brandTarget = computed(() =>
  isAdminOnlyShell.value ? ADMIN_HOME_PATH : "/",
);

const showAccountMenu = ref(false);
const accountRef = ref(null);

const handleUserUpdated = (event) => {
  session.value = event?.detail ?? getSession();
  user.value = session.value?.user ?? getSessionUser();
  showAccountMenu.value = false;
};

const handleStorage = (event) => {
  if (event.key === USER_STORAGE_KEY) {
    session.value = getSession();
    user.value = getSessionUser();
    showAccountMenu.value = false;
  }
};

onMounted(() => {
  window.addEventListener("leblanc-user-updated", handleUserUpdated);
  window.addEventListener("storage", handleStorage);
  document.addEventListener("click", handleClickOutside);
});

onBeforeUnmount(() => {
  window.removeEventListener("leblanc-user-updated", handleUserUpdated);
  window.removeEventListener("storage", handleStorage);
  document.removeEventListener("click", handleClickOutside);
});

const logoSrc = computed(() =>
  theme.value === "night" ? brightLogo : darkLogo,
);

const userInitial = computed(() =>
  (user.value?.name?.[0] || "A").toUpperCase(),
);

const logout = () => {
  clearSessionUser();
  user.value = null;
  router.push("/");
};

const toggleAccountMenu = () => {
  showAccountMenu.value = !showAccountMenu.value;
};

const handleClickOutside = (event) => {
  if (!accountRef.value) return;
  if (!accountRef.value.contains(event.target)) {
    showAccountMenu.value = false;
  }
};
</script>

<template>
  <div class="shell">
    <header v-if="!isPlainLayout" class="header">
      <RouterLink :to="brandTarget" class="brand">
        <img :src="logoSrc" alt="Le'Blanc logo" class="brand-logo" />
        <span class="brand-tag">Where every drink tells a story</span>
      </RouterLink>
      <nav class="nav">
        <template v-if="isAdminOnlyShell">
          <RouterLink to="/admin" class="nav-link" exact-active-class="active"
            >Admin</RouterLink
          >
        </template>
        <template v-else>
          <RouterLink to="/about" class="nav-link" exact-active-class="active"
            >About</RouterLink
          >
          <RouterLink to="/menu" class="nav-link" exact-active-class="active"
            >Menu</RouterLink
          >
          <RouterLink to="/booking" class="nav-link" exact-active-class="active"
            >Booking</RouterLink
          >
          <RouterLink
            v-if="isAdmin"
            :to="ADMIN_HOME_PATH"
            class="nav-link"
            exact-active-class="active"
            >Admin</RouterLink
          >
          <RouterLink
            v-if="!isAuthed"
            to="/login"
            class="nav-link"
            exact-active-class="active"
            >Login</RouterLink
          >
          <RouterLink
            v-if="!isAuthed"
            to="/register"
            class="nav-link"
            exact-active-class="active"
            >Sign up</RouterLink
          >
        </template>
        <div v-if="isAuthed" ref="accountRef" class="account-wrap">
          <button
            class="account-pill"
            type="button"
            :aria-label="`Account for ${user?.name || 'user'}`"
            @click.stop="toggleAccountMenu"
          >
            <div class="avatar" aria-hidden="true">{{ userInitial }}</div>
          </button>
          <div v-if="showAccountMenu" class="account-menu">
            <div class="account-meta">
              <span class="account-name">{{ user?.name }}</span>
              <span class="account-email">{{ user?.email }}</span>
            </div>
            <RouterLink
              v-if="isAdmin"
              :to="ADMIN_HOME_PATH"
              class="account-link"
              exact-active-class="active"
              @click="showAccountMenu = false"
            >
              Admin dashboard
            </RouterLink>
            <RouterLink
              to="/account"
              class="account-link"
              exact-active-class="active"
              @click="showAccountMenu = false"
            >
              Account
            </RouterLink>
            <RouterLink
              to="/orders"
              class="account-link"
              exact-active-class="active"
              @click="showAccountMenu = false"
            >
              My Orders
            </RouterLink>
            <button class="logout" type="button" @click="logout">
              Log out
            </button>
          </div>
        </div>
        <button
          class="theme-toggle"
          :class="{ 'is-night': theme === 'night' }"
          type="button"
          :aria-label="
            theme === 'night' ? 'Switch to day mode' : 'Switch to night mode'
          "
          @click="toggleTheme"
        >
          <span class="toggle-track" :class="{ 'is-night': theme === 'night' }">
            <span class="toggle-thumb" aria-hidden="true"></span>
          </span>
        </button>
      </nav>
    </header>

    <main class="page-frame">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.shell {
  min-height: 100vh;
  display: grid;
  gap: 12px;
}

.header {
  width: 100%;
  margin: 0;
  padding: 22px 28px 8px;
  display: flex;
  justify-content: space-between;
  gap: 18px;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--paper);
}

.brand {
  display: grid;
  gap: 6px;
  align-items: center;
  text-decoration: none;
  color: inherit;
}

.brand-logo {
  height: 96px;
  width: auto;
  display: block;
}

.brand-tag {
  font-size: 0.8rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  opacity: 0.8;
}

.nav {
  display: flex;
  gap: 20px;
  align-items: center;
  flex-wrap: wrap;
}

.account-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding: 0;
  border-radius: 999px;
  background: transparent;
  border: none;
}

.account-wrap {
  position: relative;
}

.account-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border-radius: 999px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-weight: 800;
}

.avatar {
  height: 40px;
  width: 40px;
  border-radius: 50%;
  background: linear-gradient(145deg, #b88443, #e1c58d);
  color: #0b0b0b;
  display: grid;
  place-items: center;
  font-weight: 900;
  letter-spacing: 0.02em;
}

.account-name {
  font-weight: 800;
  color: var(--ink);
}

.account-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 6px);
  min-width: 220px;
  display: grid;
  gap: 10px;
  padding: 12px;
  border-radius: 14px;
  background: var(--paper);
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 18px 36px rgba(0, 0, 0, 0.14);
  z-index: 20;
}

.account-meta {
  display: grid;
  line-height: 1.2;
}

.account-email {
  font-size: 0.9rem;
  color: rgba(0, 0, 0, 0.65);
}

.account-link {
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  text-decoration: none;
  color: var(--ink);
  font-weight: 800;
}

.logout {
  border: none;
  background: rgba(0, 0, 0, 0.08);
  padding: 10px 10px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 800;
}

.nav-link {
  text-decoration: none;
  color: var(--ink);
  font-weight: 800;
  letter-spacing: 0.02em;
  padding: 12px 0;
  font-size: 1.15rem;
}

.nav-link.active {
  color: var(--tan);
}

.theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: none;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
}

.theme-toggle.is-night {
  color: var(--tan);
}

.toggle-track {
  position: relative;
  width: 52px;
  height: 30px;
  border-radius: 999px;
  border: 1.5px solid currentColor;
  background: rgba(0, 0, 0, 0.06);
  display: inline-flex;
  align-items: center;
  padding: 2px;
  box-sizing: border-box;
  transition:
    background 0.18s ease,
    border-color 0.18s ease;
}

.toggle-track.is-night {
  background: rgba(0, 0, 0, 0.14);
  border-color: currentColor;
}

.toggle-thumb {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: currentColor;
  transition: transform 0.18s ease;
  transform: translateX(0);
}

.toggle-track.is-night .toggle-thumb {
  transform: translateX(22px);
}

.theme-toggle:hover {
  background: transparent;
}

.cta {
  padding: 10px 14px;
  background: var(--dark);
  color: #fff;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 900;
  letter-spacing: 0.04em;
}

.page-frame {
  max-width: min(1600px, 98vw);
  margin: 0 auto;
  width: 100%;
  padding: 0 24px 64px;
  box-sizing: border-box;
}

@media (max-width: 720px) {
  .brand-word {
    font-size: 1.2rem;
  }
}
</style>
