<template>
  <div>
     <div v-if="contents.length == 0">No Content</div>
     <md-card v-for="content in contents">
       <div class="md-title">{{content.local_path}}</div>
       <div class="md-subhead">{{content.download_date}}</div>
       <md-card-actions>
         <md-button @click="deleteContent(content.name)">Delete</md-button>
        </md-card-actions>
     </md-card>
  </div>
</template>

<script>
import axios from 'axios'
import moment from 'moment'

export default {
  components: {
  },
  name: 'Library',
  props: [ ],
  data: function() {
    return {
      contents: []
    }
  },
  mounted: function() {
    this.getContents()
  },
  methods: {
    deleteContent(value) {
      let url = this.$buildurl("content/" + value)
      axios.delete(url).then((response) => {
        this.getContents()
      })
    },
    getContents() {
      let url = this.$buildurl("content")
      axios.get(url).then((response) => {
        this.contents.splice(0)
        response.data.forEach((content, index) => {
          this.contents.push(content)
        })
      })
      setTimeout(function () { this.getContents() }, 60000)
    }
  }
}
</script>
<style>
</style>
