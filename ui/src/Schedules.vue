<template>
  <div>
    <div v-if="schedules.length == 0">No Schedules</div>
    <md-list class="md-scrollbar">
       <md-list-item v-for="schedule in schedules">
         <schedule-item :schedule="schedule" @deleted="getSchedules" />
       </md-list-item>
    </md-list> 
  </div>
</template>

<script>
import ScheduleItem from './ScheduleItem.vue'
import axios from 'axios'
import moment from 'moment'

export default {
  components: {
    ScheduleItem: ScheduleItem
  },
  name: 'Schedules',
  props: [ ],
  data: function() {
    return {
      schedules: []
    }
  },
  mounted: function() {
    this.getSchedules()
  },
  methods: {
    getSchedules() {
      let url = this.$buildurl("schedules")
      axios.get(url).then((response) => {
        this.schedules.splice(0)
        response.data.forEach((schedule, index) => {
          this.schedules.push(schedule)
        })
      })
    }
  }
}
</script>
<style>
</style>
