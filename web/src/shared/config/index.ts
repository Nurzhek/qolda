export const isDevEnv = import.meta.env.DEV;
export const isProdEnv = import.meta.env.PROD;

// One Go binary serves the SPA and the API, so the API is always on the same origin
export const API_BASE_URL = '/api';
