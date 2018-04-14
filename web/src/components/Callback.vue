<template>
	<div>
	</div>
</template>

<script>
import {mapActions} from 'vuex';

export default {
  name: 'Callback',
  props: ['auth'],
  mounted() {
    let error = this.auth.handleAuthentication();
    if (error) {
      this.$router.push({name: 'Login', params: {error: error}});
    }
    this.login();
    this.$router.push('/');
  },
  data() {
    return {};
  },

  methods: {
    ...mapActions(['login']),
  },
};
</script>
