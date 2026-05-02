import { createRouter, createWebHistory } from "vue-router";
import {
  ADMIN_HOME_PATH,
  getSession,
  isAdminUser,
} from "@/composables/useSessionAuth";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("@/views/Home.vue"),
    },
    {
      path: "/about",
      name: "about",
      component: () => import("@/views/About.vue"),
    },
    {
      path: "/menu",
      name: "menu",
      component: () => import("@/views/Menu.vue"),
    },
    {
      path: "/booking",
      name: "booking",
      component: () => import("@/views/Booking.vue"),
    },
    {
      path: "/verify",
      name: "verify",
      component: () => import("@/views/Verify.vue"),
      meta: { layout: "plain" },
    },
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/Login.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/register",
      name: "register",
      component: () => import("@/views/Register.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/account",
      name: "account",
      component: () => import("@/views/Account.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/admin",
      component: () => import("@/components/admin/AdminLayout.vue"),
      meta: { requiresAdmin: true },
      children: [
        {
          path: "",
          redirect: { name: "admin-users" },
        },
        {
          path: "users",
          name: "admin-users",
          component: () => import("@/views/admin/AdminUsers.vue"),
        },
        {
          path: "bookings",
          name: "admin-bookings",
          component: () => import("@/views/admin/AdminBookings.vue"),
        },
        {
          path: "drinks",
          name: "admin-drinks",
          component: () => import("@/views/admin/AdminDrinks.vue"),
        },
      ],
    },
  ],
});

router.beforeEach((to) => {
  const session = getSession();
  const user = session?.user ?? null;
  const token = session?.token || "";
  const isAuthed = Boolean(user);
  const hasAdminSession = Boolean(token) && isAdminUser(user);

  if (to.meta?.requiresAdmin) {
    if (!isAuthed || !token) {
      return { path: "/login", query: { redirect: to.fullPath } };
    }

    if (!isAdminUser(user)) {
      return { path: "/" };
    }
  }

  if (to.meta?.requiresAuth && !isAuthed) {
    return { path: "/login", query: { redirect: to.fullPath } };
  }

  if (to.meta?.guestOnly && isAuthed) {
    return { path: hasAdminSession ? ADMIN_HOME_PATH : "/" };
  }

  return true;
});

export default router;
