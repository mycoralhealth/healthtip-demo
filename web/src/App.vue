<template>
  <div id="app">
    <template v-if="isAuthenticated">
      <Navbar></Navbar>
    </template>
    <div class="container">
      <router-view 
				:auth="auth">
			</router-view>
      <template v-if="isAuthenticated">
        <Footer></Footer>
      </template>
    </div>
  </div>
</template>

<script>
import Navbar from '@/components/Navbar';
import {mapGetters, mapActions} from 'vuex';
import AuthService from '@/backend/auth/AuthService';

export default {
  name: 'app',
  computed: {
    ...mapGetters(['isAuthenticated']),
    auth: () => new AuthService(),
  },
  created() {
    console.log(this.isAuthenticated);
  },
  created() {
    this.checkCurrentLogin();
  },
  updated() {
    this.checkCurrentLogin();
  },
  methods: {
    checkCurrentLogin() {
      var unprotectedRoutes = ['/', '/callback'];
      if (
        !this.isAuthenticated &&
        unprotectedRoutes.indexOf(this.$route.path) < 0
      ) {
        this.$router.push('/');
      }
    },
  },
  components: {
    Navbar,
  },
};
</script>

<style>
</style>
