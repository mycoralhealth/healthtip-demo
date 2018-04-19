<template>
  <div class="login-overlay">
    <div class="login-wrapper border border-light">
      <form class="form-signin" @submit.prevent="login">
        <img class="logo" src="../assets/logo.png">
        <h2 class="form-signin-heading">Health Tips</h2>
        <div class="alert alert-danger" v-if="error">{{ error }}</div>
				<button class="btn btn-lg btn-primary btn-block" @click="login" :disabled="loading" type="button"><i class="fa fa-refresh fa-spin" v-if="loading"></i><div v-else="loading">Sign-in</div></button>
        <p class="text-muted copy"><small>Copyright &copy; 2018 <a href="https://mycoralhealth.com">Coral Health</a></small></p>
      </form>
    </div>
  </div>
</template>

<script>
import {mapGetters, mapActions} from 'vuex';

export default {
  name: 'Login',
  props: ['auth', 'error'],
  data() {
    return {
      email: '',
      password: '',
      loading: false,
    };
  },
  computed: {
    ...mapGetters(['isAuthenticated', 'currentUser']),
  },
  created() {
    this.checkCurrentLogin();
  },
  updated() {
    this.checkCurrentLogin();
  },
  methods: {
    checkCurrentLogin() {
      if (this.isAuthenticated) {
        this.$router.replace(this.$route.query.redirect || '/records');
      }
    },
    login() {
      this.auth.login();
    },
  },
};
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
