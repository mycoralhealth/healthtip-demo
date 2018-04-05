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
    this.auth.handleAuthentication();
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
