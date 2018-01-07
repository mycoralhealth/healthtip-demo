<template>
  <div class="signup-overlay">
    <div class="signup-wrapper border border-light">
      <form class="form-signup" @submit.prevent="signup">
        <img class="logo" src="../assets/logo.png">
        <h3 class="form-signup-heading">Get started - it's free.</h3>
        <h5 class="form-signup-subheading text-muted">Registration takes few seconds.</h5>
        <div class="alert alert-danger" v-if="error">{{ error }}</div>

        <div class="row">
          <div class="col-md-6">
            <div class="form-group">
              <label for="inputFirstName" class="sr-only">First name</label>
              <input v-model="firstName" type="text" id="inputFirstName" class="form-control" placeholder="First Name" required autofocus>
            </div>
          </div>
          <div class="col-md-6">
            <div class="form-group">
              <label for="inputLastName" class="sr-only">Last name</label>
              <input v-model="lastName" type="text" id="inputLastName" class="form-control" placeholder="Last Name" required>
            </div>
          </div>
        </div>

        <label for="inputEmail" class="sr-only">Email address</label>
        <input v-model="email" type="email" id="inputEmail" class="form-control" placeholder="Email address" required>

        <label for="inputPassword" class="sr-only">Password</label>
        <input v-model="password" type="password" id="inputPassword" class="form-control" placeholder="Password" required>
        <label for="inputConfirmPassword" class="sr-only">Confirm passwordk</label>
        <input v-model="confirmPassword" type="password" id="inputConfirmPassword" class="form-control" placeholder="Confirm Password" required>

        <button class="btn btn-lg btn-primary btn-block" :disabled="loading" type="submit"><i class="fa fa-refresh fa-spin" v-if="loading"></i><div v-else="loading">Sign-up for Health Tip</div></button>

        <p class="sign-in">
          Already a member?
          <br>
          <a href="/">Sign-in</a>
        </p>

        <p class="text-muted copy"><small>Copyright &copy; 2018 <a href="https://mycoralhealth.com">Coral Health</a></small></p>
      </form>
    </div>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'

export default {
  name: 'Signup',
  data () {
    return {
      firstName: '',
      lastName: '',
      email: '',
      password: '',
      confirmPassword: '',
      error: false,
      loading: false
    }
  },
  computed: {
    ...mapGetters({ currentUser: 'currentUser' })
  },
  created () {
    this.checkCurrentSignup()
  },
  updated () {
    this.checkCurrentSignup()
  },
  methods: {
    checkCurrentSignup () {
      if (this.currentUser) {
        this.$router.replace(this.$route.query.redirect || '/records')
      }
    },
    signup () {
      if (this.password !== this.confirmPassword) {
        this.signupFailed('The password doesn\'t match the confirmation')
        return
      }

      this.loading = true
      this.$http.post('/users', 
        {
          "email" : this.email, 
          "password" : this.password,
          "firstName" : this.firstName,
          "lastName" : this.lastName
        })
        .then(request => this.signupSuccessful(request))
        .catch(err => this.signupFailed(err.response.data))
    },
    signupSuccessful (req) {
      this.loading = false

      if (typeof(req.data.token) === "undefined" || req.data.token === null) {
        this.signupFailed()
        return
      }

      this.error = false

      localStorage.result = JSON.stringify(req.data);
      this.$store.dispatch('login')
      this.$router.replace(this.$route.query.redirect || '/records')
    },
    signupFailed (message) {
      this.loading = false
      this.error = message || 'Error processing your request'
      delete localStorage.result
    }
  }
}
</script>

<style lang="css" scoped>
body {
  background-color: white;
}

.signup-overlay {
  background-image: url("../assets/background.jpg");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  position: absolute;
  width: 100%;
  height: 140%;
  top: 0;
  left: 0;
}

.signup-wrapper {
  min-width: 400px;
  max-width: 510px;  
  background: white;
  background-color: white;
  width: 70%;
  margin: 12% auto;
  animation: fadein 0.6s;
}

@keyframes fadein {
    from { opacity: 0; }
    to   { opacity: 1; }
}

.logo {
  text-align: center;
  width: 100%;
  height: 100%;
  margin-bottom: 40px;
}

.form-signup {
  max-width: 450px;
  padding: 10% 15px;
  margin: 0 auto;
}

.form-signup-heading {
  margin-bottom: 10px;
  text-align: left;
  width: 100%;
}

.form-signup-subheading {
  margin-bottom: 30px;
  text-align: left;
  width: 100%;
}

.form-signup, 
.form-signup .checkbox {
  margin-bottom: 10px;
}

.form-signup .checkbox {
  font-weight: normal;
}

.form-signup .form-control {
  position: relative;
  height: auto;
  -webkit-box-sizing: border-box;
          box-sizing: border-box;
  padding: 10px;
  font-size: 16px;
}

.form-signup .form-control:focus {
  z-index: 2;
}

.form-signup input {
  margin-bottom: 10px;
}

.sign-in {
  margin-top: 40px;
  margin-bottom: 20px;
}

.copy {
  margin-top: 40px;
  margin-bottom: 0px;
  width: 100%;
  text-align: center;
}

</style>