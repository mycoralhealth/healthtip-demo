export default class User {
  static from(result) {
    try {
      if (!result) return null;

      var user = JSON.parse(result);

      return new User({
        api_user: user.token.apiUser,
        api_key: user.token.apiKey,
        email: user.email,
        firstName: user.firstName,
        lastName: user.lastName,
      });
    } catch (x) {
      console.log(x);
      return null;
    }
  }

  constructor({api_user, api_key, email, firstName, lastName}) {
    this.id = api_user;
    this.key = api_key;
    this.email = email;
    this.firstName = firstName;
    this.lastName = lastName;
  }

  getAuth() {
    return 'Basic ' + btoa(this.id + ':' + this.key);
  }
}
