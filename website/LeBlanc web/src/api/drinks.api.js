import apiClient, { USE_GRAPHQL } from "./httpClient";
import { getDrinkGraphQL, getDrinksGraphQL } from "./graphql.client";

export const getDrinksREST = () =>
  apiClient.get("/drinks").then((res) => res.data);

export const getDrinkREST = (id) =>
  apiClient.get(`/drinks/${id}`).then((res) => res.data);

export const getDrinks = () =>
  USE_GRAPHQL ? getDrinksGraphQL() : getDrinksREST();

export const getDrink = (id) =>
  USE_GRAPHQL ? getDrinkGraphQL(id) : getDrinkREST(id);
