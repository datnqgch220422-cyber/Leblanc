import apiClient, { USE_GRAPHQL } from "./httpClient";
import { createBookingGraphQL, getBookingsGraphQL } from "./graphql.client";

const normalizeTime = (time) => {
  if (!time) return time;
  const parsed = new Date(time);
  return Number.isNaN(parsed.getTime()) ? time : parsed.toISOString();
};

const mapBookingToGraphQLInput = (booking, normalizedTime) => ({
  email: booking.email,
  name: booking.name,
  phone: booking.phone,
  time: normalizedTime,
  guests: booking.guests ?? booking.guest ?? null,
  items: (booking.items || []).map((item) => ({
    drinkId: item.drinkId,
    qty: item.qty,
    options: item.options ? JSON.stringify(item.options) : "",
  })),
  channel: booking.channel || "web",
});

export const getBookingsREST = () =>
  apiClient.get("/bookings").then((res) => res.data);

export const getBookingsByEmailREST = (email) =>
  apiClient.get("/bookings", { params: { email } }).then((res) => res.data);

export const getBookingPaymentStatusREST = (paymentOrderId) =>
  apiClient
    .get("/bookings/payment-status", { params: { paymentOrderId } })
    .then((res) => res.data);

export const createBookingREST = (booking) =>
  apiClient.post("/bookings", booking).then((res) => res.data);

export const getBookings = () =>
  USE_GRAPHQL ? getBookingsGraphQL() : getBookingsREST();

export const getBookingsByEmail = (email) =>
  USE_GRAPHQL ? getBookingsGraphQL() : getBookingsByEmailREST(email);

export const getBookingPaymentStatus = (paymentOrderId) =>
  getBookingPaymentStatusREST(paymentOrderId);

export const createBooking = (booking) => {
  const normalizedTime = normalizeTime(booking.time);

  if (USE_GRAPHQL) {
    return createBookingGraphQL(
      mapBookingToGraphQLInput(booking, normalizedTime),
    );
  }

  return createBookingREST({ ...booking, time: normalizedTime });
};
