<template>
  <div class="login-overlay">
    <div class="login-wrapper border border-light">
      <form class="form-signin" @submit.prevent="login">
        <img class="logo" src="../assets/logo.png">
        <h2 class="form-signin-heading">Health Tip Sign-in</h2>
        <div class="alert alert-danger" v-if="error">{{ error }}</div>
        <label for="inputEmail" class="sr-only">Email address</label>
        <input v-model="email" type="email" id="inputEmail" class="form-control" placeholder="Email address" required autofocus>
        <label for="inputPassword" class="sr-only">Password</label>
        <input v-model="password" type="password" id="inputPassword" class="form-control" placeholder="Password" required>
        <button class="btn btn-lg btn-primary btn-block" :disabled="loading" type="submit"><i class="fa fa-circle-o-notch fa-spin" v-if="loading"></i><div v-else="loading">Sign-in</div></button>
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

      if (typeof(req.data.Token) === "undefined" || req.data.Token === null) {
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

<style lang="css">
body {
  background: #605B56;
}

.login-overlay {
  background: #605B56 !important;
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
}

.login-wrapper {
  background: #fff;
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

.form-signin {
  max-width: 330px;
  padding: 10% 15px;
  margin: 0 auto;
}

.form-signin .form-signin-heading,

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
</style>