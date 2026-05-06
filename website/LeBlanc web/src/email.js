// Template helper for sending verification emails using EmailJS.
// Prereqs:
// 1) Install: npm install @emailjs/browser
// 2) Set env in .env:
//    VITE_EMAILJS_SERVICE_ID=...
//    VITE_EMAILJS_TEMPLATE_ID=...
//    VITE_EMAILJS_PUBLIC_KEY=...
// 3) Configure your EmailJS template to accept variables: to_email, verify_link, token, expires_at.

import emailjs from "@emailjs/browser";
import { requestVerify, verifyToken } from "@/api";

const SERVICE_ID = import.meta.env.VITE_EMAILJS_SERVICE_ID;
const TEMPLATE_ID = import.meta.env.VITE_EMAILJS_TEMPLATE_ID;
const BOOKING_TEMPLATE_ID = import.meta.env.VITE_EMAILJS_BOOKING_TEMPLATE_ID;
const PUBLIC_KEY = import.meta.env.VITE_EMAILJS_PUBLIC_KEY;

const emailConfigured = SERVICE_ID && TEMPLATE_ID && PUBLIC_KEY;
const bookingEmailConfigured = SERVICE_ID && BOOKING_TEMPLATE_ID && PUBLIC_KEY;

// If token/verifyLink are provided (from backend register response), reuse them.
export const sendVerificationEmail = async (email, name = "", opts = {}) => {
  if (!emailConfigured)
    throw new Error("EmailJS is not configured (missing env vars).");
  let token = opts.token;
  let expiresAt = opts.expiresAt;
  let verifyLink = opts.verifyLink;

  if (!token) {
    const res = await requestVerify({ email });
    token = res.token;
    expiresAt = res.expiresAt;
  }
  if (!verifyLink) {
    verifyLink = `${window.location.origin}/verify?token=${encodeURIComponent(token)}`;
  }
  await emailjs.send(
    SERVICE_ID,
    TEMPLATE_ID,
    {
      to_email: email,
      to: email,
      name,
      to_name: name,
      user_name: name,
      username: name,
      verify_link: verifyLink,
      verifyUrl: verifyLink,
      token,
      expires_at: expiresAt,
    },
    PUBLIC_KEY,
  );
  return { token, expiresAt, verifyLink };
};

export const checkVerificationToken = async (token) => {
  return verifyToken({ token });
};

export const isEmailReady = () => Boolean(emailConfigured);

export const isBookingEmailReady = () => Boolean(bookingEmailConfigured);

export const sendBookingEmail = async (booking = {}, opts = {}) => {
  if (!bookingEmailConfigured) {
    throw new Error(
      "EmailJS booking template is not configured (missing env vars).",
    );
  }
  if (!booking.email) {
    throw new Error("Booking email is missing.");
  }

  const timeText = booking.time ? new Date(booking.time).toLocaleString() : "";
  const items = booking.items || [];
  const itemsText = items.length
    ? items
        .map((item, idx) => {
          const name =
            item.name ||
            item.drink?.name ||
            item.drinkName ||
            item.drinkId ||
            `Item ${idx + 1}`;
          const qty = item.qty || 1;
          return `${idx + 1}. ${name} x${qty}`;
        })
        .join("\n")
    : "No drinks pre-ordered.";

  await emailjs.send(
    SERVICE_ID,
    BOOKING_TEMPLATE_ID,
    {
      to_email: booking.email,
      customer_email: booking.email,
      customer_name: booking.name,
      customer_phone: booking.phone,
      time: timeText,
      guests: booking.guests || booking.guest || "",
      items_text: itemsText,
      notes: booking.note || booking.notes || "",
      booking_id: opts.bookingId || booking.bookingId || "",
      channel: booking.channel || "web",
    },
    PUBLIC_KEY,
  );
};
