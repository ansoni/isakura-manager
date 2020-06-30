<template>
  <div>
    <md-card class="schedule">
      <div class="md-title">
          Name: {{schedule.name}}
      </div>
      <div class="md-subhead">
         Filter: {{schedule.filter}}
      </div>
      <div class="md-subhead">
         PrependDate: {{schedule.prependDate}}
      </div>
      <div class="md-subhead">
         PrependTime: {{schedule.prependTime}} </div>
      <div class="md-subhead">
         AppendDate: {{schedule.apependDate}}
      </div>
      <div class="md-subhead">
         AppendTime: {{schedule.apependTime}}
      </div>
      <div class="md-subhead">
         Watch: {{schedule.watch}}
      </div>
      <div class="md-subhead">
         Search/Replace
         <div v-for="search in schedule.searches">
           Replace {{search.search}} with {{search.replace}}
         </div>
      </div>
      <md-card-actions>
        <md-button @click="deleteScheduleItem">Delete</md-button>
      </md-card-actions>
    </md-card
  </div>
</template>

<script>
import axios from 'axios'
import moment from 'moment'

export default {
  name: 'ScheduleItem',
  props: [ "schedule" ],
  mounted: function() {
    console.log(this.schedule)
  },
  methods: {
    deleteScheduleItem() {
      let url = this.$buildurl("schedules/" + this.schedule.name)
      axios.delete(url).then((response) => {
        this.$emit("deleted")
      })
    }
  }
}
</script>
<style>
.schedule {
  width: 600px;
}
</style>
