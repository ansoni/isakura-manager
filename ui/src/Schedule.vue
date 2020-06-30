<template>
  <div>
    <md-steppers md-alternative :md-active-step.sync="active" md-linear>
      <md-step id="first" md-label="Filter">
        <div class="md-layout md-gutter md-alignment-center-center">
          <div class="md-layout-item md-size-large">
          </div>
          <div class="md-layout-item md-size-large">
            <md-field>
              <label>Listing Filter</label>
              <md-input v-on:keyup.enter="matches" v-model="filter"></md-input>
            </md-field> 
          </div>
          <div class="md-layout-item md-size-large">
            <md-button @click="matches" class="md-icon-button">
              <md-icon>search</md-icon>
            </md-button>
            Matches {{ filter_matches.length }}
          </div>
        </div>
        <md-button :disabled="filter_matches.length == 0" class="md-raised md-primary" @click="setDone('first', 'second')">Next</md-button> 
        <div>
          <md-list class="matches md-scrollbar">
            <md-list-item v-for="match in filter_matches">
              <span>{{match.channel}} - {{match.start}} - {{match.title}}</span>
            </md-list-item>
          </md-list> 
        </div>
      </md-step>

      <md-step id="second" md-label="Set Name">
        <div class="md-layout md-gutter md-alignment-center-center">
          <div class="md-layout-item md-size-small">
            <md-checkbox v-model="prependDate" value="1">Prepend Date</md-checkbox>
            <md-checkbox v-model="prependTime" value="1">Prepend Time</md-checkbox>
            <md-checkbox v-model="appendDate" value="1">Append Date</md-checkbox>
            <md-checkbox v-model="appendTime" value="1">Append Time</md-checkbox>
          </div>
          <div class="md-layout-item md-size-small">
            <div class="md-layout md-gutter md-alignment-center-center">
              <div class="md-layout-item md-size-small">
                <md-field>
                  <label>Search</label>
                  <md-input v-model="search"></md-input>
                </md-field> 
              </div>
              <div class="md-layout-item md-size-small">
                <md-field>
                  <label>Replace</label>
                  <md-input v-model="replace"></md-input>
                </md-field>
              </div>
              <div class="md-layout-item md-size-small">
                <md-button :disabled="search == ''" class="md-raised md-primary" @click="saveSearchAndReplace()"><md-icon>add</md-icon></md-button> 
              </div>
              <div class="md-layout-item md-size-small">
                <div v-for="(search_and_replace, i) in searches">
                  search for {{search_and_replace.search}}, replace with {{search_and_replace.replace}}
                  <md-button @click="searches.splice(i,1)"><md-icon>delete</md-icon></md-button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="md-layout-item md-size-small">
          <md-button class="md-raised md-primary" @click="setDone('second', 'third')">Next</md-button> 
        </div>
        <div>
          <md-list class="matches md-scrollbar">
            <md-list-item v-for="match in filter_matches">
              <span>{{modifyTitle(match)}}</span>
            </md-list-item>
          </md-list> 
        </div>
      </md-step>

      <md-step id="third" md-label="Config">
        <div>
            <md-field>
              <label>Cool Name</label>
              <md-input v-model="name"></md-input>
            </md-field> 
        </div>
        <div>
          <md-autocomplete v-model="folder" :md-options="folders" @md-opened="getFolders">
            <label>Folder</label>
            <template slot="md-autocomplete-item" slot-scope="{ item }">{{ item }}</template>
          </md-autocomplete>
        </div>
        <div>
          <md-checkbox v-model="watch">Watch for future matches</md-checkbox>
        </div>
        <md-button class="md-raised md-primary" @click="setDone('third', 'forth')">Next</md-button> 
      </md-step>

      <md-step id="forth" md-label="Review">
        <div>
           <div>Cool Name: {{name}}</div>
           <div>Folder: {{folder}}</div>
           <div>Filter: {{filter}}</div>
           <div>Prepend Date: {{prependDate}}</div>
           <div>Prepend Time: {{prependTime}}</div>
           <div>Append Date: {{appendDate}}</div>
           <div>Append Time: {{appendTime}}</div>
        <div>
        <md-button class="md-raised md-primary" @click="save">Save</md-button> 
      </md-step>
    </md-steppers>
  </div>
</template>

<script>
import axios from 'axios'
import moment from 'moment'
  export default {
    name: 'Schedule',
    props: [ "listing" ],
    data: function() {
      return {
        active: "first",
        first: false,
        second: false,
        third: false,
        forth: false,
        filter_matches: [],
        filter: "",
        folders: [],
        folder: "",
        watch: true,
        search: "",
        replace: "",
        searches: [],
        appendDate: false,
        prependDate: false,
        appendTime: false,
        prependTime: false,
        name: ""
      }
    },
    methods: {
      save() {
        let url = this.$buildurl("schedules")
        let searches = []
        this.searches.forEach((search, index) => {
          searches.push({search:search.search, replace: search.replace})
        })
        let schedule={ filter: this.filter, folder: this.folder, appendDate: this.appendDate, preprendDate: this.prependDate, appendTime: this.appendTime, appendDate: this.appendDate, name: this.name, watch: this.watch, searches: searches} 
        console.log("New Schedule - ")
        console.log(schedule)
        axios.post(url, schedule).then((response) => {
          alert(response.data)
          this.reset()
          this.$emit("schedule") 
        }).catch(function (error) {
          alert(error)
        })
      },
      saveSearchAndReplace() {
        this.searches.push({search: this.search, replace: this.replace})
      },
      getFolders(value) {
        let url = this.$buildurl("content/folders")
        axios.get(url).then((response) => {
          this.folders.splice(0)
          response.data.forEach((folder, index) => {
            this.folders.push(folder)
          })
        })
      },
      modifyTitle(value) {
        let toreturn = ""
        if (this['prependDate']) {
          toreturn+=moment(value.airdate).format("YYYY-MM-DD") + " "
        }
        if (this.prependTime) {
          toreturn+=moment(value.airdate).format("HH:mm") + " "
        }
        let title = value.title
        this.searches.forEach((search, index) => {
          title = title.replace(search.search, search.replace)   
        })
        toreturn+=title
        if (this.appendDate) {
          toreturn+=" " + moment(value.airdate).format("YYYY-MM-DD") 
        }
        if (this.appendTime) {
          toreturn+=" " + moment(value.airdate).format("HH:mm") 
        }
        return toreturn
      },
      matches() {
        if (this.filter == "") {
          return
        }
        let url = this.$buildurl("channels/guide?search=" + this.filter)
        axios.get(url).then((response) => {
          this.filter_matches.splice(0) 
          response.data.forEach((channel, index) => {
            channel.guide.forEach((listing,index) => {
              event={ channel: channel.channelName, id: listing.id, airdate: listing.airdate, title: listing.name, start: moment(listing.airdate).format("YYYY-MM-DD HH:mm")};
              this.filter_matches.push(event)
            })
          })
        })
      },
      reset(value) {
        console.log("Reset")
        this.first=false
        this.second=false
        this.third=false
        this.forth=false
        this.appendDate=false
        this.appendTime=false
        this.prependDate=false
        this.prependTime=false
        this.searches=[]
        this.name=""
        this.folder=""
        this.filter_matches.splice(0)
        if (value) {
          this.filter=value.title
        }
        this.active="first"
      },
      
      setDone (id, index) {
        this[id] = true

        if (index) {
          this.active = index
        }
      }
    } 
  }
</script>
<style>
.matches {
  height: 400px;
  overflow: auto;
}
</style>
