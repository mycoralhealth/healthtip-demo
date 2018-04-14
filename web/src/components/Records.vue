<template>
    <main role="main" class="records-overlay">
			<modals-container/>
      <!-- Main jumbotron for a primary marketing message or call to action -->
      <div class="jumbotron jumbotron-push">
        <div class="container">
          <h1 class="display-3">Welcome to Health Tips!</h1>
          <p>We, at Coral Health, are using the Blockchain to improve the state of healthcare. This is a demo app to illustrate how anyone can get
          personalized health tips easily through a full-fledged decentralized medical records system. Start by adding some test records.</p>
        </div>
      </div>

      <div class="container">
        <!-- Example row of columns -->
        <div class="row">
          <div class="col-md-4">
            <h2>Add Test Result</h2>
							<error-bar :error="error" :onDismissHandler="dismissError"/>
              <form class="form-new-result" @submit.prevent="addRecord">
                <div class="form-group">
                  <label for="inputAge" class="sr-only">Age</label>
                  <div class="input-group">
                    <input type="number" min="1" max="150" step="1" id="inputAge" class="form-control" placeholder="Age" v-model="record.age" required>
                  </div>
                  <label for="inputHeight" class="sr-only">Height</label>
                  <div class="input-group">
                    <input type="number" min="0" max="300" class="form-control" id="inputHeight" placeholder="Height" v-model="record.height" required>
                    <div class="input-group-append">
                      <span class="input-group-text">cm</span>
                    </div>
                  </div>
                  <label for="inputWeight" class="sr-only">Weight</label>
                  <div class="input-group">
                    <input type="number" min="0" max="1000" class="form-control" id="inputWeight" placeholder="Weight" v-model="record.weight" required>
                    <div class="input-group-append">
                      <span class="input-group-text">kg</span>
                    </div>
                  </div>
                  <label for="inputCholesterol" class="sr-only">Heart Rate</label>
                  <div class="input-group">
                    <input type="number" min="0" class="form-control" id="inputCholesterol" placeholder="Heart Rate: beats/minute" v-model="record.cholesterol" required>
                    <div class="input-group-append">
                      <span class="input-group-text">bpm</span>
                    </div>
                  </div>
                  <label for="inputBloodPressure" class="sr-only">Respiratory Rate</label>
                  <div class="input-group">
                    <input type="number" min="0" class="form-control" id="inputBloodPressure" placeholder="Breath Rate: breaths/minute" v-model="record.bloodPressure" required>
                    <div class="input-group-append">
                      <span class="input-group-text">bpm</span>
                    </div>
                  </div>
									<h2>For Insurance Approval</h2>
									<label for="inputNumberCysts" class="sr-only">Number of Painful Cysts</label>
									<div class="input-group">
										<input type="number" class="form-control" id="importNumberCysts" placeholder="Number of Painful Cysts" v-model="record.numberOfCysts" required>
									</div>
									<label for="inputBaldess" class="sr-only">Baldness</label>
									<div class="form-group">
										<select required class="form-control" id="selectBaldness" v-model="record.baldness">
											<option disabled selected value="">Baldness?</option>
											<option v-for="option in baldnessOptions" v-bind:value="option.value">
												{{option.text}}
											</option>
										</select>
									</div>
									<div class="form-group" v-if="record.baldness">
										<select required class="form-control" id="selectBaldFromDisease" v-model="record.baldnessFromDisease">
											<option disabled selected value="">Baldness due to disease?</option>
											<option v-for="option in baldnessOptions" v-bind:value="option.value">
												{{option.text}}
											</option>
										</select>
									</div>
                  <button class="btn btn-secondary" :disabled="loading" type="submit" href="#" role="button"><i class="fa fa-refresh fa-spin" v-if="loading"></i><div v-else="loading">Add</div></button>
                </div>
              </form>
          </div>
          <div class="col-md-8">
            <h2>Previous Results</h2>

            <div class="alert alert-warning text-center" role="alert" v-if="records.length == 0">
              It looks like you don't have any test results submitted yet. Go ahead add some.
            </div>

            <table class="table table-striped text-center">
              <thead>
                <tr>
                  <th>Age</th>
                  <th>Height (cm)</th>
                  <th>Weight (kg)</th>
                  <th>Heart Rate (bpm)</th>
                  <th>Breath Rate (bpm)</th>
									<th>Number of Cysts</th>
									<th>Baldness?</th>
									<th>Baldness from disease?</th>
                  <th></th>
                </tr>
              </thead>
              <tbody v-for="(r, index) in records">
                <tr class="">
                  <th scope="row">{{ r.age }}</th>
                  <td>{{ r.height }}</td>
                  <td>{{ r.weight }}</td>
                  <td>{{ r.cholesterol }}</td>
                  <td>{{ r.bloodPressure }}</td>
									<td>{{ r.numberOfCysts}}</td>
									<td>{{ r.baldness ? "Yes" : "No" }}</td>
									<td>{{ r.baldnessFromDisease ? "Yes" : "No" }}</td>
                  <td><a href="#" v-on:click="deleteRecord(index)"><i class="fa fa-trash-o"></i></a></td>
                </tr>
                <tr class="text-left">
                  <td colspan="2">
										<button type="button" class="btn btn-outline-secondary disabled" v-if="r.tipSent == 1">Tip Requested</button>
										<button type="button" class="btn btn-outline-success" v-else="r.tipSent == 1" v-on:click="requestTip(index)">Request Tip</button>
										<br><br>
									</td>
									<td colspan="3">
										<button type="button" class="btn btn-outline-success" @click="showApprovalPopup(index)">Request Insurance Approval</button>
									</td>
                </tr>
              </tbody>
            </table>
            <div class="list-end"></div>
          </div>
        </div>

        <hr>
        <p class="text-muted text-center copy"><small>Copyright &copy; 2018 <a href="https://mycoralhealth.com">Coral Health</a></small></p>
      </div> <!-- /container -->
      <simplert :useRadius="true"
                :useIcon="true"
                ref="simplert">
      </simplert>
    </main>
</template>


<script>
import {mapGetters} from 'vuex';
import Simplert from 'vue2-simplert';
import ApprovalPopup from './ApprovalPopup.vue';
import ErrorBar from './ErrorBar.vue';

export default {
  name: 'Records',
  components: {Simplert, ApprovalPopup, ErrorBar},
  computed: {
    ...mapGetters(['authString', 'currentUser', 'isAuthenticated']),
  },
  data() {
    return {
      record: {
        age: '',
        height: '',
        weight: '',
        cholesterol: '',
        bloodPressure: '',
        numberOfCysts: '',
        baldness: '',
        baldnessFromDisease: '',
      },
      records: [],
      loading: false,
      error: false,
      baldnessSelected: '',
      baldnessOptions: [{text: 'No', value: false}, {text: 'Yes', value: true}],
      error: false,
    };
  },
  created() {
    this.checkCurrentLogin();
    this.fetchRecords();
  },

  methods: {
    checkCurrentLogin() {
      if (!this.isAuthenticated) {
        this.$router.push('/');
      }
    },

    fetchRecords() {
      this.$http
        .get('/api/records', {
          headers: {Authorization: this.authString},
        })
        .then(request => this.recordsLoaded(request))
        .catch(() => this.loadAPIError());
    },

    addRecord() {
      this.loading = true;
      if (this.record.baldnessFromDisease === '') {
        this.record.baldnessFromDisease = false;
      }
      if (this.record.age.trim()) {
        this.$http
          .post('/api/records', this.record, {
            headers: {Authorization: this.authString},
          })
          .then(request => this.appendRecordResult(request))
          .catch(err => this.reportError(err.response.data));
      }
    },

    deleteRecord(index) {
      var that = this;

      let confirmFn = function() {
        that.$http
          .delete('api/records/' + that.records[index].id, {
            headers: {Authorization: that.authString},
          })
          .then(() => that.removeRecordFromResult(index))
          .catch(err => that.reportError(err.response.data));
      };

      let obj = {
        title: 'Delete Test Result',
        message: 'Are you sure you want to delete this test result?',
        type: 'warning',
        customConfirmBtnText: 'Delete',
        customConfirmBtnClass:
          'simplert__confirm simplert__confirm--radius bg-danger',
        useConfirmBtn: true,
        onConfirm: confirmFn,
      };
      this.$refs.simplert.openSimplert(obj);
    },

    clearExistingRecord() {
      this.record.age = '';
      this.record.height = '';
      this.record.weight = '';
      this.record.cholesterol = '';
      this.record.bloodPressure = '';
      this.record.numberOfCysts = '';
      this.record.baldness = '';
      this.record.baldnessFromDisease = '';
    },

    appendRecordResult(req) {
      this.clearExistingRecord();
      this.loading = false;
      this.records.push(req.data);
      var container = this.$el.querySelector('.list-end');
      container.scrollIntoView();
    },

    removeRecordFromResult(index) {
      this.records.splice(index, 1);
    },

    recordsLoaded(req) {
      this.records = req.data;
    },

    reportError(err) {
      this.loading = false;
      this.error = err;
      var container = this.$el.querySelector('.container');
      container.scrollIntoView();
    },

    loadAPIError() {
      this.$store.dispatch('logout');
      this.$router.push({
        name: 'Login',
        params: {error: 'Remote server unavailable. Please try again later.'},
      });
    },

    requestTip(index) {
      this.$http
        .post(
          '/api/records/' + this.records[index].id + '/tip',
          {name: this.currentUser.name, email: this.currentUser.email},
          {
            headers: {Authorization: this.authString},
          },
        )
        .then(req => this.requestTipSuccess(req.data, index))
        .catch(err => this.reportError(err.response.data));
    },

    requestTipSuccess(record, index) {
      this.records.splice(index, 1, record);
      let obj = {
        title: 'Request Sent',
        message:
          'Your request for a Health Tip was sent. You should receive a reponse in the next 48 hours.',
        type: 'success',
      };
      this.$refs.simplert.openSimplert(obj);
    },

    dismissError() {
      this.error = false;
    },

    showApprovalPopup(index) {
      this.$modal.show(
        ApprovalPopup,
        {
          recordId: this.records[index].id,
          onCloseHandler: this.hideApprovalPopup,
        },
        {
          name: 'insurance-approval',
          adaptive: true,
          width: 600,
          height: 300,
        },
      );
    },

    hideApprovalPopup() {
      this.$modal.hide('insurance-approval');
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="css" scoped>

.records-overlay {
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  position: absolute;
  width: 100%;
  height: 100%;
  top: 50;
  left: 0;
}

.jumbotron-push { margin-top: 50px; }

.form-new-result .input-group {
  margin-bottom: 10px;
}

.copy {
  margin-top: 10px;
  margin-bottom: 30px;
  width: 100%;
  text-align: center;
}

.btn-spacer {
  margin-right: 15px;
}

select:invalid {
	color: #868e95;
}

</style>
