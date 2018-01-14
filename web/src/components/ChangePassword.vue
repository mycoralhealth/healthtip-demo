<template>
  <div class="forgot-overlay">
    <div class="forgot-wrapper border border-light">
      <form class="form-forgot" @submit.prevent="changePassAttempt">
        <img class="logo" src="../assets/logo.png">
        <h5 class="form-forgot-heading">Change your password.</h5>
        <h6 class="form-forgot-subheading text-muted">Please enter a new password of 6 characters or more.</h6>
        <div class="alert alert-danger" v-if="error">{{ error }}</div>
        <div class="alert alert-info" v-if="success">
          <button type="button" class="close" data-dismiss="alert" aria-label="Close" v-on:click="dismissSuccessMessage()">
            <span aria-hidden="true">&times;</span>
          </button>
          Password change successful. Goto <strong><a href="/">Sign-in</a></strong>
        </div>

        <label for="inputPassword" class="sr-only">Password</label>
        <input v-model="password" type="password" id="inputPassword" class="form-control" placeholder="Password (minimum length 6 characters)" required>
        <label for="inputConfirmPassword" class="sr-only">Confirm password</label>
        <input v-model="confirmPassword" type="password" id="inputConfirmPassword" class="form-control" placeholder="Confirm Password" required>

        <button class="btn btn-lg btn-primary btn-block" :disabled="loading" type="submit"><i class="fa fa-refresh fa-spin" v-if="loading"></i><div v-else="loading">Change Password</div></button>

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

export default {
  name: 'ChangePassword',
  data () {
    return {
      token: '',
      password: '',
      confirmPassword: '',
      error: false,
      loading: false,
      success: false
    }
  },
  created () {
    this.checkToken()
  },
  methods: {
    checkToken () {
      this.token = this.$route.query.token
      if (!this.token || this.token.length == 0) {
        this.$router.replace('/')
      }
    },
    changePassAttempt () {
      if (!this.password || this.password.length < 6) {
        this.changeFailed('Required password length of 6 or more characters')
        return        
      }

      if (this.password !== this.confirmPassword) {
        this.changeFailed('The password doesn\'t match the confirmation')
        return
      }

      this.loading = true
      this.$http.post('/claimToken', 
        {
          "apiUser" : '-1',
          "apiKey" : this.token
        })
        .then(request => this.changePass(request))
        .catch(err => this.changeFailed(err.response.data))
    },
    changePass (req) {
      if (!req.data.token.apiUser || !req.data.token.apiKey) {
        changeFailed(nil)
        return
      }

      this.loading = true
      this.$http.post('/changePassword', 
        {
          "password" : this.password
        },
        {headers: {'Authorization': 'Basic ' + btoa(req.data.token.apiUser + ':' + req.data.token.apiKey)}}
        )
        .then(request => this.changeSuccessful(request))
        .catch(err => this.changeFailed(err.response.data))
    },
    changeSuccessful (req) {
      this.loading = false
      this.error = false
      this.success = true
    },
    changeFailed (message) {
      this.success = false;
      this.loading = false
      this.error = message || 'Error processing your request'
    },
    dismissSuccessMessage() {
      this.success = false      
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