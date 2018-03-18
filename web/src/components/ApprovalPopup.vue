<template>
	<div class="container text-center">
		<div class="row">
			<div class="col-md-12">
				<div class="spacer">
					<button type="button" class="close" aria-label="Close" @click="onCloseHandler()">
						<span>&times;</span>
					</button>
				</div>
				<form class="form-approval" v-if="formStage === 0">
					<div class="form-group">
						<h3> Select the procedure for which you want Insurance Approval </h3>
						<div class="input-group">
							<select required class="form-control" id="selectProcedure" v-model="approvalRequest.procedure">
								<option disabled selected value="">Select Procedure</option>
								<option v-for="option in options.procedures" v-bind:value="option.id">
									{{option.name}}
								</option>
							</select>
						</div>	
					</div>
					<div class="form-group">
						<h3> Select Insurance Company </h3>
						<div class="input-group">
							<select required class="form-control" id="selectCompany" v-model="approvalRequest.company">
								<option disabled selected value="">Select Company</option>
								<option v-for="option in options.companies" v-bind:value="option.id">
									{{option.name}}
								</option>
							</select>
						</div>
						<div>
							<button role="button" class="btn btn-outline-success" type="submit" @click="submitApprovalRequest()">Submit Request</button>
						</div>
					</div>
				</form>
				<div v-if="formStage === 1">
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
	name: 'ApprovalPopup',
	props: ['recordId','onCloseHandler', 'options'],
  computed: {
    ...mapGetters({ currentUser: 'currentUser' })
  },
	data() {
		return {
			formStage: 0,
			approvalRequest: {
				procedure: '',
				company: ''
			},
			approvalResponse: {},
		}
	},
	methods: {
		submitApprovalRequest() {
			this.$http.post('/api/records/' + this.recordId + '/approval', this.approvalRequest, {headers: {'Authorization': this.currentUser.getAuth()}})
				.then(response => this.displayApprovalResults(reponse))
		},

		displayApprovalResults(response) {
			this.approvalResponse = response;
			this.formStage = 1; 
		},
	}
}
</script>

<style lang="css" scoped>

.form-approval .input-group {
	margin-top:10px;
	margin-bottom:10px;
}

.spacer {
	margin-top:15px;
}

select:invalid {
	color: #868e95;
}

</style>
