<template>
	<div class="container text-center">
		<div class="row">
			<div class="col-md-12">
				<div class="spacer">
					<button type="button" class="close" aria-label="Close" @click="onCloseHandler()">
						<span>&times;</span>
					</button>
				</div>
				<form class="form-approval" v-if="formStage === 0" @submit.prevent="submitApprovalRequest">
					<div class="form-group">
						<h3> Select the procedure for which you want Insurance Approval </h3>
						<div class="input-group">
							<select required class="form-control" id="selectProcedure" v-model="approvalRequest.procedure">
								<option disabled selected value="">Select Procedure</option>
								<option v-for="option in options.procedures" v-bind:value="option">
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
								<option v-for="option in options.companies" v-bind:value="option">
									{{option.name}}
								</option>
							</select>
						</div>
						<div>
							<button :disabled="options.loading" role="button" class="btn btn-outline-success" type="submit">
								<i class="fa fa-refresh fa-spin" v-if="options.loading"></i>
								<div v-else="options.loading"> Submit Request </div>
							</button>
						</div>
					</div>
				</form>
				<div class="form-approval" v-if="formStage === 1">
					<div>
						<img :src="approvalImage">
					</div>
					<div>
						{{generateApprovalText()}}
					</div>
					<div>
						<button role="button" class="btn btn-outline-success" type="button" @click="downloadMedicalPolicy()">View {{approvalResponse.company.name}}'s medical policy</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { mapGetters } from 'vuex'
import checkmark from '../assets/icons8-checkmark.svg'
import cancel from '../assets/icons8-cancel.svg'

export default {
	name: 'ApprovalPopup',
	props: ['recordId','onCloseHandler'],
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
			options: {
				companies: [],
				procedures:[],
				loading:false
			},
			approvalImage: checkmark
		}
	},

	mounted() {
		this.fetchApprovalOptions()
	},

	methods: {
		submitApprovalRequest() {
			this.$http.post('/api/records/' + this.recordId + '/approval', this.approvalRequest, {headers: {'Authorization': this.currentUser.getAuth()}})
				.then(response => this.displayApprovalResults(response.data))
				.catch(function(err) {
					});
		},

		displayApprovalResults(response) {
			this.approvalResponse = response;
			if(response.approved === false) {
				this.approvalImage = cancel;
			} else {
				this.approvalImage = checkmark;
			}
			this.formStage = 1; 
		},

		generateApprovalText() {
			const approved = this.approvalResponse.approved
			var approvalText = "According to " + this.approvalResponse.company.name+"’s medical policy, "
			approvalText = approvalText + "you " + (!approved ? "do not " : " ")
			approvalText = approvalText + "qualify for " + this.approvalResponse.procedure.name + "."
			return approvalText
		},

		fetchApprovalOptions() {
			var that = this;
			that.options.loading = true;
			Promise.all([
				that.$http.get('/api/companies', {headers: {'Authorization': this.currentUser.getAuth()}}),
				that.$http.get('/api/procedures', {headers: {'Authorization': this.currentUser.getAuth()}})
			]).then(function ([companies, procedures]) {
				that.options.companies = companies.data;
				that.options.procedures = procedures.data;
				that.options.loading = false;
			}).catch(err => this.reportError(err));
		},

		downloadMedicalPolicy() {
			this.$http({
				url: '/api/companies/'+this.approvalResponse.company.id+'/procedures/'+this.approvalResponse.procedure.id+'/policy',
				method: 'GET',
				headers: {'Authorization': this.currentUser.getAuth()},
				responseType: 'blob', // important
			}).then(response => {
				const disposition = response.headers["content-disposition"];
				if(disposition && disposition.indexOf('=') !== -1) {
					var filename = disposition.split('=')[1];
				}
				const url = window.URL.createObjectURL(new Blob([response.data]));
				const link = document.createElement('a');
				link.href = url;
				link.setAttribute('download', filename);
				document.body.appendChild(link);
				link.click();
			});
		}
	}
}
</script>

<style lang="css" scoped>

.form-approval .input-group {
	margin-top:10px;
	margin-bottom:10px;
}

.close {
	font-size:24px
}

.text-field {
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
