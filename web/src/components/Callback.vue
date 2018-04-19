<template>
	<div>
	</div>
</template>

<script>
export default {
  name: 'Callback',
  props: ['auth'],
  mounted() {
    this.auth.handleAuthentication();
  },
  data() {
    return {};
  },
};
</script>
