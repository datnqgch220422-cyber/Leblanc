import apiClient, { USE_GRAPHQL } from "./httpClient";
import { recoFromFeaturesGraphQL } from "./graphql.client";

const mapEmotionFitPayload = (payload) =>
  payload.emotionFit || {
    calm: payload.calm || 0,
    happy: payload.happy || 0,
    stressed: payload.stressed || 0,
    sad: payload.sad || 0,
    adventurous: payload.adventurous || 0,
  };

export const recoFromFeaturesREST = (payload) =>
  apiClient.post("/reco/from-features", payload).then((res) => res.data);

export const recoFromFeatures = (payload) => {
  if (USE_GRAPHQL) {
    return recoFromFeaturesGraphQL(
      mapEmotionFitPayload(payload),
      payload.caffeine,
      payload.temp,
      payload.sweetness,
    );
  }

  return recoFromFeaturesREST(payload);
};
