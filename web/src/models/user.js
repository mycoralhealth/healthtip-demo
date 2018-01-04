
export default class User {
  static from (result) {
    try {
      var user = JSON.parse(result)

      return new User({
        "api_user": user.Token.Api_user,
        "api_key": user.Token.Api_key,
        "email": user.Email,
        "firstName": user.First_name,
        "lastName": user.Last_name
      })
    } catch (x) {
      console.log(x)
      console.log("Exception")
      return null
    }
  }

  constructor ({ api_user, api_key, email, firstName, lastName }) {
    this.id = api_user
    this.key = api_key 
    this.email = email
    this.firstName = firstName
    this.lastName = lastName
  }
}