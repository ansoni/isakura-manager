<template>
  <div class="app">
    <md-tabs class="md-primary" :md-active-tab="currenttab" @md-changed="tabChanged">
      <md-tab id="tab-guide" md-label="Guide" to="/guide">
        <div class="md-layout md-gutter">
          <div class="md-layout-item md-size-large"></div>
          <div class="md-layout-item md-size-xsmall">
            <md-field>
              <label for="movie">Channel</label>
              <md-select v-model="current_channel" id="channel" name="channel">
                <md-option v-for="channel in channels" id="channel.id" value="channel.id">{{channel.id}}</md-option>
              </md-select>
            </md-field>
          </div>
          <div class="md-layout-item md-size-large">
            <md-field md-clearable @md-clear="doClear">
              <label>Search the Guide</label>
              <md-input v-on:keyup.enter="doSearch" v-model="search_text"></md-input>
            </md-field>
          </div>
          <div class="md-layout-item md-size-xsmall">
            <md-button @click="doSearch" class="md-icon-button">
              <md-icon>search</md-icon>
            </md-button>
          </div>
          <div class="md-layout-item md-size-large"></div>
        </div>
        <!--channel v-if="current_channel != null" v-bind:channel="current_channel" v-bind:search="search" v-bind:dateNow="dateNow" @schedule="onScheduleRequest" /-->
        <!--div class="md-layout md-gutter md-alignment-center-center">
          <div class="md-layout-item md-size-large">
          </div>
          <div class="md-layout-item md-size-large">
            <span v-if="total_results > 0">Results: {{total_results}}</span>
          </div>
          <div class="md-layout-item md-size-large">
          </div>
        </div>
        <div class="md-layout md-gutter">
          <div class="md-layout-item">
            <channel v-for="channel in channels" v-bind:channel="channel" v-bind:search="search" v-bind:dateNow="dateNow" @results_count="onResultsCount" @schedule="onScheduleRequest" />
          </div>
       </div-->
      </md-tab>
      <md-tab id="tab-schedule" md-label="Schedule" to="/schedule">
        <schedule ref="schedule" @schedule="newSchedule"/>
      </md-tab>
      <md-tab id="tab-schedules" md-label="Schedules" to="/schedules">
        <schedules ref="schedules"/>
      </md-tab>
      </md-tab>
      <md-tab id="tab-library" md-label="Library" to="/library">
        <library />
      </md-tab>
    </md-tabs>
  </div>
</template>

<script>

import Library from './Library.vue'
import Channel from './Channel.vue'
import Schedule from './Schedule.vue'
import Schedules from './Schedules.vue'
import axios from 'axios'
import moment from 'moment'

export default {
  components: {
    Channel: Channel,
    Schedule: Schedule,
    Schedules: Schedules,
    Library: Library
  },
  data: function() {
    return {
      listing: {},
      currenttab: "tab-guide",
      dateNow: new Date().valueOf()/1000,
      total_results: 0,
      search_text: "",
      search: "",
      channelGuides: {},
      current_channel: "",
      channels: [],
      channelUpdated: 0,
      selectedChannels: []
    };
  },
  created: function() {
    setInterval(() => this.dateNow = new Date().valueOf()/1000, 1000 * 60)
  },
  mounted: function() {
    this.getChannels();
  },
  methods: {
    tabChanged(value) {
      this.currenttab=value
    },
    onResultsCount(value) {
      this.total_results+=value
    },
    newSchedule() {
      this.currenttab="tab-schedules"
      this.$refs.schedules.getSchedules()
    },
    onScheduleRequest(value) {
      console.log(value)
      this.listing=value
      this.currenttab="tab-schedule"
      this.$refs.schedule.reset(value)
      console.log("Schedule Request" + value)
    },
    containsKey(obj, key) {
      return Object.keys(obj).includes(key);
    },
    channelGuide(channel) {
      if (channel in this.channelGuides) {
        return this.channelGuides[channel]
      }
      return []
    },
    events() {
      this.channels.forEach((value, index) => {
        let url = this.$buildurl("channels/" + value.id + "/guide?search=" + this.search)
        let channel = value.id
        axios.get(url).then((response) => {
          if (response.status == 200) {
            let channelGuides=[]
            response.data.guide.forEach((value, index) => {
              event={ title: value.name, start: moment(value.airdate).format("YYYY-MM-DD HH:mm")};
              channelGuides.push(event)
            })
            this.$set(this.channelGuides, channel, channelGuides)
            this.channelUpdated+=1;
          } else {
            console.log("Failure! - " + response);
          }
        })
      })
    },
    doClear() {
      this.search_text=""
      this.doSearch()
    },
    doSearch() {
      this.total_results=0
      console.log("Search for " + this.search_text)
      this.search = this.search_text; // promote to children
    },
    getChannels() {
      let url = this.$buildurl("channels");
      axios.get(url).then((response) => {
        this.channels.slice(0)
        response.data.forEach((value, index) => {
          this.$set(this.channelGuides, value.channelName, [])
          this.channels.push({ id: value.channelName })
        })
        this.events();
        this.loading = false;
      }).catch((error) => { console.log(error); });  
    },
    toggleWeekends() {
      this.calendarWeekends = !this.calendarWeekends; // update a property
    },
    gotoPast() {
      let calendarApi = this.$refs.fullCalendar.getApi(); // from the ref="..."
      calendarApi.gotoDate("2000-01-01"); // call a method on the Calendar object
    },
    handleDateClick(arg) {
      /*if (confirm("Would you like to add an event to " + arg.dateStr + " ?")) {
        this.calendarEvents.push({
          // add new event data
          title: "New Event",
          start: arg.date,
          allDay: arg.allDay
        });
      }*/
      alert("Click");
    },
    hasGuide(channel) {
       return this.containsKey(this.channelGuides, channel)
    }
  }
};
</script>

<style>
.app {
  font-family: Arial, Helvetica Neue, Helvetica, sans-serif;
  font-size: 14px;
}

.md-app {
    border: 1px solid rgba(#000, .12);
  }

.app-top {
  margin: 0 0 3em;
}

.app-calendar {
 margin: 0 auto;
 width: auto;
 z-index: 1000;
}
</style>
