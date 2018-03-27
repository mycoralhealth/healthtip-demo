<template>
	<div class="alert alert-danger alert-dismissible" v-if="error">
		<button type="button" class="close" data-dismiss="alert" aria-label="Close" v-on:click="onDismissHandler()">
			<span aria-hidden="true">&times;</span>
		</button>
		{{ error }}
	</div>
</template>

<script>
export default {
	name: "ErrorBar",
	props: ['error', 'onDismissHandler'],
}
</script>
