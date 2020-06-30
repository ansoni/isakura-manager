<template>
  <md-list v-if="show === true">
    <md-list-item md-expand :md-expanded.sync="expanded">
      <md-badge v-if="! loading" md-position="top" class="md-primary" :md-content="listings_length"></md-badge>
      <md-progress-spinner v-if="loading" :md-diameter="20" :md-stroke="3" md-mode="indeterminate"></md-progress-spinner>
      <span class="md-list-item-text">{{channel.id}}</span>
      <md-list slot="md-expand" class="listings md-scrollbar">
        <md-list-item v-if="expanded" v-for="(listing,i) in listings">
          <div>
            <span>
              <md-button class="md-icon-button md-dense md-raised" @click="schedule(listing)">
                <md-icon>schedule</md-icon>
              </md-button>
            </span>
            <preview v-if="listing.airdate < dateNow" :channel_name="channel.id" :listing_id="listing.id"></preview>
            <md-button v-else disabled class="md-icon-button md-dense md-raised">
              <md-icon>play_circle_outline</md-icon>
            </md-button>
            <span>{{listing.start}} - {{listing.title}}</span>
            <md-divider></md-divider>
          </div>
        </md-list-item>
      </md-list>
    </md-list-item>
  </md-list>
</template>

<script>
import Preview from './Preview.vue'
import axios from 'axios'
import moment from 'moment'

module.exports = {
  components: {
    Preview: Preview
  },
  props: [ 'channel', 'search', 'dateNow' ],
  mounted: function() {
    this.channelRef=this.channel
    this.events()
  }, 
  watch: {
    search: function(search) {
      this.events()
    }
  },
  methods: {
    schedule: function(value) {
       this.$emit('schedule', value)
    },
    events() {
      let url = this.$buildurl("channels/" + encodeURI(this.channel.id) + "/guide?search=" + this.search)
      let channel = this.channel.id
      this.show=true
      this.$nextTick(function() {
        this.loading=true
      })
      axios.get(url).then((response) => {
        if (response.status == 200) {
          this.listings.splice(0) 
          response.data.guide.forEach((value, index) => {
            event={ id: value.id, airdate: parseInt(value.id), title: value.name, start: moment(value.airdate).format("YYYY-MM-DD HH:mm")};
            this.listings.push(event)
          })
          this.listings_length=this.listings.length
          this.$emit('results_count', this.listings_length)
        } else {
          console.log("Failure! - " + response);
        }
        if (this.listings_length > 0) {
          this.show=true
          this.$nextTick(function () {
            this.loading=false
          })
        } else {
          this.show=false
        }
      })
    },
  },
  data: function () {
    return {
      show: false,
      loading: true,
      channelRef: {},
      expanded: false,
      listings: [],
      listings_length: 0,
    }
  }
}
</script>
<style>
.listings {
  max-height: 400px;
  overflow: auto;
}
</style>
