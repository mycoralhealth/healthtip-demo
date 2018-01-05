<template>
  <div>

    <div class="container" id="records">
      <div class="col-sm-7">
        <div class="panel panel-default">
          <div class="panel-heading">
            <h3>Add new medical record:</h3>
          </div>
          <div class="panel-body">
            <div>
              <input type="number" min="0" step="1" class="form-control" placeholder="Age" v-model="record.Age">
              <input type="number" min="0" class="form-control" placeholder="Height" v-model="record.Height">
              <input type="number" min="0" class="form-control" placeholder="Weight" v-model="record.Weight">
              <input type="number" min="0" class="form-control" placeholder="Cholesterol" v-model="record.Cholesterol">
              <input type="number" min="0" class="form-control" placeholder="Blood Pressure" v-model="record.Blood_pressure">
              <button class="btn btn-primary" v-on:click="addRecord">Add</button>
            </div>
          </div>
        </div>
      </div>
      <div class="col-sm-5">
        <div class="list-group">
          <a href="#" class="list-group-item" v-for="(r, index) in records">
            <h4 class="list-group-item-heading"><i class="glyphicon glyphicon-bullhorn"></i> {{ r.age }}</h4>
            <p class="list-group-item-text" v-if="r.height">Height: {{ r.height }}</p>
            <p class="list-group-item-text" v-if="r.weight">Weight: {{ r.weight }}</p>
            <p class="list-group-item-text" v-if="r.cholesterol">Cholesterol: {{ r.cholesterol }}</p>
            <p class="list-group-item-text" v-if="r.blood_pressure">Blood Pressure: {{ r.blood_pressure }}</p>
            <button class="btn btn-xs btn-danger" v-on:click="deleteRecord(index)">Delete</button>
          </a>
        </div>
      </div>
    </div>

  </div>
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
      record: { Age: '', Height: '', Weight: '', Cholesterol: '', Blood_pressure: '' },
      records: []
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
      console.log(this.currentUser.getAuth())

      if (this.record.Age.trim()) {
        this.$http.post('/api/records', this.record, {headers: {'Authorization': this.currentUser.getAuth()}})
          .then(request => this.appendRecordResult(request))
          .catch(err => this.reportError(err));
      }
    },

    deleteRecord (index) {
      if (confirm('Are you sure you want to delete this record?')) {
        // this.events.splice(index, 1);
        this.$http.delete('api/records/' + this.records[index].id, {headers: {'Authorization': this.currentUser.getAuth()}})
          .then(() => this.removeRecordFromResult(index))
          .catch(err => this.reportError(err));
      }
    },

    appendRecordResult(req) {
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
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
