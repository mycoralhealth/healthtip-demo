<template>
  <div id="app">
    <template v-if="currentUser">
      <Navbar></Navbar>
    </template>
    <div class="container">
      <router-view></router-view>
      <template v-if="currentUser">
        <Footer></Footer>
      </template>
    </div>
  </div>
</template>

<script>
import Navbar from '@/components/Navbar'
import { mapGetters } from 'vuex'

export default {
  name: 'app',
  computed: {
    ...mapGetters({ currentUser: 'currentUser' })
  },
  created () {
    console.log(this.currentUser)
  },
  created () {
    this.checkCurrentLogin()
  },
  updated () {
    this.checkCurrentLogin()
  },
  methods: {
    checkCurrentLogin () {
      var unprotectedRoutes = ['/', '/signup', '/forgot', '/changePass']
      if (!this.currentUser && (unprotectedRoutes.indexOf(this.$route.path) < 0)) {
        this.$router.push('/')
      }
    }
  },
  components: {
    Navbar
  }
}
</script>

<style>
</style>
