import auth0 from 'auth0-js';
import router from '@/router';
import store from '@/store';
import {AUTH0_CONFIG} from './auth0-variables';

export default class AuthService {
  constructor() {
    this.login = this.login.bind(this);
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
      let error = null;
      if (authResult && authResult.accessToken && authResult.idToken) {
        localStorage.setItem('authResult', JSON.stringify(authResult));
      } else if (err) {
        console.log(err);
        error = err;
      }
      this.loginHandler(err);
    });
  }

  loginHandler(err) {
    if (err) {
      router.push({name: 'Login', params: {error: error}});
    }
    store.dispatch('login');
    router.push('/');
  }
}
