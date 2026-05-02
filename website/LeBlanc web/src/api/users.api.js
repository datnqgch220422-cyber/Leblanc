import apiClient, { USE_GRAPHQL } from "./httpClient";
import { getUsersGraphQL } from "./graphql.client";

export const getUsersREST = () =>
  apiClient.get("/users").then((res) => res.data);

export const getUsers = () =>
  USE_GRAPHQL ? getUsersGraphQL() : getUsersREST();
