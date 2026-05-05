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
      path: "/booking/addons",
      name: "booking-addons",
      component: () => import("@/views/BookingAddons.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/cart",
      redirect: { path: "/booking/addons" },
    },
    {
      path: "/orders",
      name: "orders",
      component: () => import("@/views/MyOrders.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/orders/:id",
      name: "order-detail",
      component: () => import("@/views/BookingDetail.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/booking/confirmation",
      name: "booking-confirmation",
      component: () => import("@/views/BookingConfirmation.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/booking/payment-result",
      name: "booking-payment-result",
      component: () => import("@/views/PaymentResult.vue"),
      meta: { requiresAuth: true },
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
      component: () => import("@/views/Admin.vue"),
      meta: { requiresAdmin: true },
      children: [
        {
          path: "users",
          redirect: { path: "/admin" },
        },
        {
          path: "bookings",
          redirect: { path: "/admin" },
        },
        {
          path: "drinks",
          redirect: { path: "/admin" },
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
