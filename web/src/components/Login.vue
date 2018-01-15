<template>
  <div class="login-overlay">
    <div class="login-wrapper border border-light">
      <form class="form-signin" @submit.prevent="login">
        <img class="logo" src="../assets/logo.png">
        <h2 class="form-signin-heading">Health Tip</h2>
        <div class="alert alert-danger" v-if="error">{{ error }}</div>
        <label for="inputEmail" class="sr-only">Email address</label>
        <input v-model="email" type="email" id="inputEmail" class="form-control" placeholder="Email address" required autofocus>
        <label for="inputPassword" class="sr-only">Password</label>
        <input v-model="password" type="password" id="inputPassword" class="form-control" placeholder="Password" required>
        <button class="btn btn-lg btn-primary btn-block" :disabled="loading" type="submit"><i class="fa fa-refresh fa-spin" v-if="loading"></i><div v-else="loading">Sign-in</div></button>
        <small><a href="/forgot">I forgot my password</a></small>

        <p class="sign-up">
          Not a member yet?
          <br>
          <a href="/signup">Sign-up</a>
        </p>

        <p class="text-muted copy"><small>Copyright &copy; 2018 <a href="https://mycoralhealth.com">Coral Health</a></small></p>
      </form>
    </div>
  </div>
</template>

<script>

import { mapGetters } from 'vuex'

export default {
  name: 'Login',
  data () {
    return {
      email: '',
      password: '',
      error: false,
      loading: false
    }
  },
  computed: {
    ...mapGetters({ currentUser: 'currentUser' })
  },
  created () {
    this.checkCurrentLogin()
  },
  updated () {
    this.checkCurrentLogin()
  },
  methods: {
    checkCurrentLogin () {
      if (this.currentUser) {
        this.$router.replace(this.$route.query.redirect || '/records')
      }
    },
    login () {
      this.loading = true
      this.$http.post('/login', {"email" : this.email, "password" : this.password}, {headers: {'Authorization': 'Basic ' + btoa(this.email + ':' + this.password)}})
        .then(request => this.loginSuccessful(request))
        .catch(() => this.loginFailed())
    },
    loginSuccessful (req) {
      this.loading = false

      if (typeof(req.data.token) === "undefined" || req.data.token === null) {
        this.loginFailed()
        return
      }

      this.error = false

      localStorage.result = JSON.stringify(req.data);
      this.$store.dispatch('login')
      this.$router.replace(this.$route.query.redirect || '/records')
    },
    loginFailed () {
      this.loading = false
      this.error = 'Invalid email address or password'
      this.$store.dispatch('logout')
      delete localStorage.result
    }
  }
}
</script>

<style lang="css" scoped>
body {
  background-color: white;
}

.login-overlay {
  background-image: url("../assets/background.jpg");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  position: absolute;
  width: 100%;
  height: 120%;
  top: 0;
  left: 0;
}

.login-wrapper {
  width: 400px;
  min-width: 260px;
  max-width: 400px;  
  background: white;
  background-color: white;
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

.form-signin {
  max-width: 330px;
  padding: 10% 15px;
  margin: 0 auto;
}

.form-signin-heading {
  margin-bottom: 30px;
  text-align: center;
  width: 100%;
}

.form-signin, 
.form-signin .checkbox {
  margin-bottom: 10px;
}

.form-signin .checkbox {
  font-weight: normal;
}

.form-signin .form-control {
  position: relative;
  height: auto;
  -webkit-box-sizing: border-box;
          box-sizing: border-box;
  padding: 10px;
  font-size: 16px;
}

.form-signin .form-control:focus {
  z-index: 2;
}

.form-signin input[type="email"] {
  margin-bottom: -1px;
  border-bottom-right-radius: 0;
  border-bottom-left-radius: 0;
}

.form-signin input[type="password"] {
  margin-bottom: 10px;
  border-top-left-radius: 0;
  border-top-right-radius: 0;
}

.sign-up {
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