<template>
  <nav class="navbar navbar-dark bg-dark fixed-top">
    
    <span class="navbar-left logo"><img style="width: 85px; float: left" src="../assets/logo_small.png"></img><h4 class="navbar-text title text-muted">CORAL HEALTH</h4></span>
    <span class="navbar-text">
      <div class="name">
        {{ currentUser.firstName }}&nbsp;{{ currentUser.lastName }}
      </div>
      <router-link to="/logout" class="btn btn-logout btn-sm">Logout</router-link>
    </span>

  </nav>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Navbar',
  computed: {
    ...mapGetters({ currentUser: 'currentUser' })
  }
}
</script>

<style lang="css" scoped>

@media screen and (max-width: 550px) {
  .name { display: none; }
}

@media screen and (min-width: 550px) {
  .name { display: inline; }
}

.logo {
  margin-top: 5px; 
  margin-bottom: -5px;
}

.title {
  margin-left: 10px;
  margin-top: 0px;
  height: 100%;
}

.bg-records {
  background: #3798cf !important;
}

.btn-logout {
  margin-left: 10px;
  background-color: #837A75;
  border-color: #837A75;
  color: #fff;
}

.btn-logout:hover {
  background-color: darken(#837A75, 20%);
  border-color: darken(#837A75, 20%);
  color: #fff;
}
</style>
