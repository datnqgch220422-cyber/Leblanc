import { GraphQLClient, gql } from "graphql-request";

const graphqlEndpoint =
  (import.meta.env.VITE_API_BASE || "http://localhost:4000") + "/graphql";

const client = new GraphQLClient(graphqlEndpoint, {
  headers: {},
});

export const GET_USERS_QUERY = gql`
  query GetUsers {
    users {
      _id
      name
      email
      createdAt
    }
  }
`;

export const GET_DRINKS_QUERY = gql`
  query GetDrinks {
    drinks {
      _id
      name
      price
      tags
      caffeine
      temp
      sweetness
      colorTone
      emotionFit {
        calm
        happy
        stressed
        sad
        adventurous
      }
      image
      desc
    }
  }
`;

export const GET_DRINK_QUERY = gql`
  query GetDrink($id: ID!) {
    drink(id: $id) {
      _id
      name
      price
      tags
      caffeine
      temp
      sweetness
      colorTone
      emotionFit {
        calm
        happy
        stressed
        sad
        adventurous
      }
      image
      desc
    }
  }
`;

export const GET_BOOKINGS_QUERY = gql`
  query GetBookings {
    bookings {
      _id
      email
      name
      phone
      time
      guests
      items {
        drinkId
        qty
        options
      }
      status
      channel
    }
  }
`;

export const CREATE_BOOKING_MUTATION = gql`
  mutation CreateBooking($input: CreateBookingInput!) {
    createBooking(input: $input) {
      _id
      email
      name
      phone
      time
      guests
      items {
        drinkId
        qty
        options
      }
      status
      channel
    }
  }
`;

export const REGISTER_MUTATION = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      ok
      user {
        _id
        name
        email
        createdAt
      }
    }
  }
`;

export const LOGIN_MUTATION = gql`
  mutation Login($input: LoginInput!) {
    login(input: $input) {
      ok
      user {
        _id
        name
        email
        createdAt
      }
    }
  }
`;

export const RECOMMEND_FROM_FEATURES_MUTATION = gql`
  mutation RecommendFromFeatures(
    $emotionFit: EmotionFitInput!
    $caffeine: String
    $temp: String
    $sweetness: Int
  ) {
    recommendFromFeatures(
      emotionFit: $emotionFit
      caffeine: $caffeine
      temp: $temp
      sweetness: $sweetness
    ) {
      drinkId
      score
    }
  }
`;

export const getDrinksGraphQL = async () => {
  const data = await client.request(GET_DRINKS_QUERY);
  return data.drinks;
};

export const getUsersGraphQL = async () => {
  const data = await client.request(GET_USERS_QUERY);
  return data.users;
};

export const getDrinkGraphQL = async (id) => {
  const data = await client.request(GET_DRINK_QUERY, { id });
  return data.drink;
};

export const getBookingsGraphQL = async () => {
  const data = await client.request(GET_BOOKINGS_QUERY);
  return data.bookings;
};

export const createBookingGraphQL = async (input) => {
  const data = await client.request(CREATE_BOOKING_MUTATION, { input });
  return data.createBooking;
};

export const registerUserGraphQL = async (input) => {
  const data = await client.request(REGISTER_MUTATION, { input });
  return data.register;
};

export const loginUserGraphQL = async (input) => {
  const data = await client.request(LOGIN_MUTATION, { input });
  return data.login;
};

export const recoFromFeaturesGraphQL = async (
  emotionFit,
  caffeine,
  temp,
  sweetness,
) => {
  const variables = { emotionFit };
  if (caffeine) variables.caffeine = caffeine;
  if (temp) variables.temp = temp;
  if (sweetness) variables.sweetness = sweetness;

  const data = await client.request(
    RECOMMEND_FROM_FEATURES_MUTATION,
    variables,
  );
  return data.recommendFromFeatures;
};

export default client;
