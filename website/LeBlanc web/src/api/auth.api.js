import apiClient from "./httpClient";

export const registerUserREST = (payload) =>
  apiClient.post("/auth/register", payload).then((res) => res.data);

export const loginUserREST = (payload) =>
  apiClient.post("/auth/login", payload).then((res) => res.data);

export const requestVerifyREST = (payload) =>
  apiClient.post("/auth/request-verify", payload).then((res) => res.data);

export const verifyTokenREST = (payload) =>
  apiClient.post("/auth/verify", payload).then((res) => res.data);

export const registerUser = (payload) => registerUserREST(payload);

export const loginUser = (payload) => loginUserREST(payload);

export const requestVerify = (payload) => requestVerifyREST(payload);

export const verifyToken = (payload) => verifyTokenREST(payload);
