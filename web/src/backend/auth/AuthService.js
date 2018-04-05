import auth0 from 'auth0-js';

export default class AuthService {
  constructor() {
    this.login = this.login.bind(this);
    this.isAuthenticated = this.isAuthenticated.bind(this);
    this.handleAuthentication = this.handleAuthentication.bind(this);
  }

  auth0 = new auth0.WebAuth({
    domain: 'mycoralhealth.auth0.com',
    clientID: 'n0gKuFFCKyJIKeW4miwDyPdfqCJ9yt7q',
    redirectUri: 'http://localhost:8080/callback',
    audience: 'https://tips.mycoralhealth.com/api/',
    responseType: 'token id_token',
    scope: 'openid',
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
