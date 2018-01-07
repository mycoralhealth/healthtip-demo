<template>
    <main role="main" class="records-overlay">

      <!-- Main jumbotron for a primary marketing message or call to action -->
      <div class="jumbotron jumbotron-push">
        <div class="container">
          <h1 class="display-3">Hello, welcome to Health Tips!</h1>
          <p>We at Coral Health are working on the future of the healthcare ecosystem. This is a demo app to illustrate how anyone can get
          personalized health tips easily with a full fledged decentralized medical records system. Start by adding some test records.</p>
          <p><a class="btn btn-primary btn-lg" href="https://www.mycoralhealth.com" role="button">Learn more &raquo;</a></p>
        </div>
      </div>

      <div class="container">
        <!-- Example row of columns -->
        <div class="row">
          <div class="col-md-4">
            <h2>Add Test Result</h2>
              <div class="alert alert-danger" v-if="error">{{ error }}</div>
              <form class="form-new-result" @submit.prevent="addRecord">
                <div class="form-group">
                  <label for="inputAge" class="sr-only">Age</label>
                  <div class="input-group">
                    <input type="number" min="0" step="1" id="inputAge" class="form-control" placeholder="Age" v-model="record.age" required>
                  </div>
                  <label for="inputHeight" class="sr-only">Height</label>
                  <div class="input-group">
                    <input type="number" min="0" class="form-control" id="inputHeight" placeholder="Height" v-model="record.height" required>
                    <div class="input-group-append">
                      <span class="input-group-text">cm</span>
                    </div>
                  </div>
                  <label for="inputWeight" class="sr-only">Weight</label>
                  <div class="input-group">
                    <input type="number" min="0" class="form-control" id="inputWeight" placeholder="Weight" v-model="record.weight" required>
                    <div class="input-group-append">
                      <span class="input-group-text">kg</span>
                    </div>
                  </div>
                  <label for="inputCholesterol" class="sr-only">Cholesterol</label>
                  <div class="input-group">
                    <input type="number" min="0" class="form-control" id="inputCholesterol" placeholder="Cholesterol" v-model="record.cholesterol" required>
                    <div class="input-group-append">
                      <span class="input-group-text">mg/dL</span>
                    </div>
                  </div>
                  <label for="inputBloodPressure" class="sr-only">Blood pressure</label>
                  <div class="input-group">
                    <input type="number" min="0" class="form-control" id="inputBloodPressure" placeholder="Blood Pressure" v-model="record.bloodPressure" required>
                    <div class="input-group-append">
                      <span class="input-group-text">mmHg</span>
                    </div>
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
                  <th>Cholesterol (mg/dL)</th>
                  <th>Blood Pressure (mmHg)</th>
                  <th></th>
                </tr>
              </thead>
              <tbody v-for="(r, index) in records">
                <tr>
                  <th scope="row">{{ r.age }}</th>
                  <td>{{ r.height }}</td>
                  <td>{{ r.weight }}</td>
                  <td>{{ r.cholesterol }}</td>
                  <td>{{ r.bloodPressure }}</td>
                  <td><a href="#" v-on:click="deleteRecord(index)"><i class="fa fa-trash-o"></i></a></td>
                </tr>
                <tr class="text-left">
                  <td colspan="6"><button type="button" class="btn btn-outline-success">Request Tip</button><br><br></td>
                </tr>
              </tbody>
            </table>

          </div>
        </div>

        <hr>
        <p class="text-muted text-center copy"><small>Copyright &copy; 2018 <a href="https://mycoralhealth.com">Coral Health</a></small></p>
      </div> <!-- /container -->

    </main>

</template>


<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Records',
  computed: {
    ...mapGetters({ currentUser: 'currentUser' })
  },
  data () {
    return {
      record: { age: '', height: '', weight: '', cholesterol: '', bloodPressure: '' },
      records: [],
      loading: false,
      error: false
    }
  },
  created () {
    this.checkCurrentLogin()
    this.fetchRecords()
  },
  methods: {
    checkCurrentLogin () {
      if (!this.currentUser) {
        this.$router.push('/')
      }
    },

    fetchRecords () {
      this.$http.get('/api/records', {headers: {'Authorization': this.currentUser.getAuth()}})
        .then(request => this.recordsLoaded(request))
        .catch(() => this.loadAPIError())
    },

    addRecord () {
      this.loading = true
      console.log(this.currentUser.getAuth())

      if (this.record.age.trim()) {
        this.$http.post('/api/records', this.record, {headers: {'Authorization': this.currentUser.getAuth()}})
          .then(request => this.appendRecordResult(request))
          .catch(err => this.reportError(err));
      }
    },

    deleteRecord (index) {
      if (confirm('Are you sure you want to delete this record?')) {
        this.$http.delete('api/records/' + this.records[index].id, {headers: {'Authorization': this.currentUser.getAuth()}})
          .then(() => this.removeRecordFromResult(index))
          .catch(err => this.reportError(err));
      }
    },

    appendRecordResult(req) {
      this.loading = false
      this.records.push(req.data)
    },

    removeRecordFromResult(index) {
      this.records.splice(index, 1);
    },

    recordsLoaded (req) {
      this.records = req.data;
      console.log(this.records);
    },

    reportError(err) {
      this.loading = false
      console.log(err)
    },

    loadAPIError() {
      this.$store.dispatch('logout')
      this.$router.push('/')
    }
  }
}
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

.jumbotron-push {
  margin-top: 50px;
}

.form-new-result .input-group {
  margin-bottom: 10px;
}

.copy {
  margin-top: 10px;
  margin-bottom: 30px;
  width: 100%;
  text-align: center;
}

</style>
