<template>
  <div id="app">
    <template v-if="currentUser">
      <Navbar></Navbar>
    </template>
    <div class="container">
      <router-view 
				:auth="auth">
			</router-view>
      <template v-if="currentUser">
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
    ...mapGetters({currentUser: 'currentUser'}),
    auth: () => new AuthService(),
  },
  created() {
    console.log(this.currentUser);
  },
  created() {
    this.checkCurrentLogin();
  },
  updated() {
    this.checkCurrentLogin();
  },
  methods: {
    checkCurrentLogin() {
      var unprotectedRoutes = [
        '/',
        '/callback',
        '/signup',
        '/forgot',
        '/changePass',
      ];
      if (
        !this.currentUser &&
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
