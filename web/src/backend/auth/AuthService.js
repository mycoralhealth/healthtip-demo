import auth0 from 'auth0-js';
import {AUTH0_CONFIG} from './auth0-variables';

export default class AuthService {
  constructor() {
    this.login = this.login.bind(this);
    this.isAuthenticated = this.isAuthenticated.bind(this);
    this.handleAuthentication = this.handleAuthentication.bind(this);
  }

  auth0 = new auth0.WebAuth({
    domain: AUTH0_CONFIG.domain,
    clientID: AUTH0_CONFIG.clientID,
    redirectUri: AUTH0_CONFIG.redirectUri,
    audience: AUTH0_CONFIG.audience,
    responseType: AUTH0_CONFIG.responseType,
    scope: AUTH0_CONFIG.scope,
  });

  login() {
    this.auth0.authorize();
  }

  handleAuthentication() {
    this.auth0.parseHash((err, authResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        localStorage.setItem('authResult', JSON.stringify(authResult));
      } else if (err) {
        console.log(err);
        // alert(`Error: ${err.error}. Check the console for further details.`);
      }
    });
  }

  isAuthenticated() {
    // Check whether the current time is past the
    // Access Token's expiry time
    let expiresAt = JSON.parse(localStorage.getItem('expires_at'));
    return new Date().getTime() < expiresAt;
  }
}
