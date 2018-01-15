<template>
  <div class="forgot-overlay">
    <div class="forgot-wrapper border border-light">
      <form class="form-forgot" @submit.prevent="forgotPass">
        <img class="logo" src="../assets/logo.png">
        <h5 class="form-forgot-heading">Reset your password.</h5>
        <h6 class="form-forgot-subheading text-muted">Please enter your email address and we'll send you instructions on how to reset your password.</h6>
        <div class="alert alert-danger" v-if="error">{{ error }}</div>
        <div class="alert alert-info" v-if="emailSent">
          <button type="button" class="close" data-dismiss="alert" aria-label="Close" v-on:click="dismissSentMessage()">
            <span aria-hidden="true">&times;</span>
          </button>
          We've sent you instructions on how to reset your password. Please check your email.
        </div>
        <label for="inputEmail" class="sr-only">Email address</label>
        <input v-model="email" type="email" id="inputEmail" class="form-control" placeholder="Email address" required>

        <button class="btn btn-lg btn-primary btn-block" :disabled="loading" type="submit"><i class="fa fa-refresh fa-spin" v-if="loading"></i><div v-else="loading">Reset Password</div></button>

        <p class="sign-in">
          Go back to Login?
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
  name: 'Forgot',
  data () {
    return {
      email: '',
      error: false,
      loading: false,
      emailSent: false
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
    forgotPass () {
      this.loading = true
      this.$http.post('/resetPassword', 
        {
          "email" : this.email
        })
        .then(request => this.forgotSuccessful(request))
        .catch(err => this.forgotFailed(err.response.data))
    },
    forgotSuccessful (req) {
      this.loading = false
      this.error = false
      this.emailSent = true
    },
    forgotFailed (message) {
      this.emailSent = false;
      this.loading = false
      this.error = message || 'Error processing your request'
    },
    dismissSentMessage() {
      this.emailSent = false      
    }
  }
}
</script>

<style lang="css" scoped>
body {
  background-color: white;
}

.forgot-overlay {
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

.forgot-wrapper {
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

.form-forgot {
  max-width: 450px;
  padding: 10% 15px;
  margin: 0 auto;
}

.form-forgot-heading {
  margin-bottom: 10px;
  text-align: left;
  width: 100%;
}

.form-forgot-subheading {
  margin-bottom: 30px;
  text-align: left;
  width: 100%;
}

.form-forgot, 
.form-forgot .checkbox {
  margin-bottom: 10px;
}

.form-forgot .checkbox {
  font-weight: normal;
}

.form-forgot .form-control {
  position: relative;
  height: auto;
  -webkit-box-sizing: border-box;
          box-sizing: border-box;
  padding: 10px;
  font-size: 16px;
}

.form-forgot .form-control:focus {
  z-index: 2;
}

.form-forgot input {
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