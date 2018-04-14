const CALLBACK_URL =
  process.env.CALLBACK_URL || 'https://tips.mycoralhealth.com/callback';

export const AUTH0_CONFIG = {
  domain: 'mycoralhealth.auth0.com',
  clientID: 'n0gKuFFCKyJIKeW4miwDyPdfqCJ9yt7q',
  redirectUri: CALLBACK_URL,
  audience: 'https://tips.mycoralhealth.com/api/',
  responseType: 'token id_token',
  scope: 'openid email profile',
};
