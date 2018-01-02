<template>
  <div class="login-wrapper border border-light">
    <form class="form-signin" @submit.prevent="login">
      <img class="logo" src="../assets/logo.png">
      <h2 class="form-signin-heading">Health Tip Sign-in</h2>
      <div class="alert alert-danger" v-if="error">{{ error }}</div>
      <label for="inputEmail" class="sr-only">Email address</label>
      <input v-model="email" type="email" id="inputEmail" class="form-control" placeholder="Email address" required autofocus>
      <label for="inputPassword" class="sr-only">Password</label>
      <input v-model="password" type="password" id="inputPassword" class="form-control" placeholder="Password" required>
      <button class="btn btn-lg btn-primary btn-block" type="submit"><i class="fa fa-circle-o-notch fa-spin" v-if="loading">&nbsp;</i><div v-else="loading">Sign-in</div></button>
    </form>
  </div>
</template>

<script>
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
  updated () {
    if (localStorage.token) {
      this.$router.replace(this.$route.query.redirect || '/records')
    }
  },
  methods: {
    login () {
      this.loading = true
      this.$http.post('/login', {"email" : this.email, "password" : this.password}, {headers: {'Authorization': 'Basic ' + btoa(this.email + ':' + this.password)}})
        .then(request => this.loginSuccessful(request))
        .catch(() => this.loginFailed())
    },
    loginSuccessful (req) {
      this.loading = false

      if (typeof(req.data.Api_user) === "undefined" || req.data.Api_user === null) {
        this.loginFailed()
        return
      }

      this.error = false

      localStorage.token = btoa(req.data.Api_user + ':' + req.data.Api_key)

      this.$router.replace(this.$route.query.redirect || '/records')
    },
    loginFailed () {
      this.loading = false
      this.error = 'Invalid email address or password'
      delete localStorage.token
    }
  }
}
</script>

<style lang="css">
body {
  background: #605B56;
}

.login-wrapper {
  background: #fff;
  width: 70%;
  margin: 12% auto;
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
  margin-bottom: 20px;
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