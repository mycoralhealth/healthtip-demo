<template>
  <div>

    <div class="container" id="events">
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
          <a href="#" class="list-group-item" v-for="record in records">
            <h4 class="list-group-item-heading"><i class="glyphicon glyphicon-bullhorn"></i> {{ record.Age }}</h4>
            <p class="list-group-item-text" v-if="record.Height">Height: {{ record.Height }}</p>
            <p class="list-group-item-text" v-if="record.Weight">Weight: {{ record.Weight }}</p>
            <p class="list-group-item-text" v-if="record.Cholesterol">Cholesterol: {{ record.Cholesterol }}</p>
            <p class="list-group-item-text" v-if="record.Blood_pressure">Blood Pressure: {{ record.Blood_pressure }}</p>
            <button class="btn btn-xs btn-danger" v-on:click="deleteRecord($index)">Delete</button>
          </a>
        </div>
      </div>
    </div>

  </div>
</template>


<script>
export default {
  name: 'Records',
  data () {
    return {
      record: { Age: '', Height: '', Weight: '', Cholesterol: '', Blood_pressure: '' },
      records: []
    }
  },
  ready: function() {
    this.fetchRecords();
  },
  methods: {

    fetchRecords: function () {
      var events = [];

      this.$http.get('/api/records')
        .success(function (records) {
          this.$set('records', records);
          console.log(records);
        })
        .error(function (err) {
          console.log(err);
        });
    },

    addRecord: function () {
      if (this.record.Age.trim()) {
        this.$http.post('/api/records', this.record)
          .success(function (res) {
            this.records.push(this.record);
            console.log('Record added!');
          })
          .error(function (err) {
            console.log(err);
          });
      }
    },

    deleteRecord: function (index) {
      if (confirm('Are you sure you want to delete this record?')) {
        // this.events.splice(index, 1);
        this.$http.delete('api/records/' + record.id)
          .success(function (res) {
            console.log(res);
            this.events.splice(index, 1);
          })
          .error(function (err) {
            console.log(err);
          });
      }
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
