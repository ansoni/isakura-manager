<template>
  <span>
    <md-dialog :md-active.sync="show">
      <md-dialog-title>Preview</md-dialog-title>
      <md-dialog-actions>
        <md-button class="md-primary" @click="toggleShow">Close</md-button>
      </md-dialog-actions>
      <flv-player ref="fp" :mediaDataSource='{ isLive: true, type: "flv", url: this.previewUrl()}' :config='{ lazyLoadMaxDuration: 5 }'/>
    </md-dialog>

    <md-button class="md-icon-button md-dense md-raised" @click="play">
      <md-icon>play_circle_outline</md-icon>
    </md-button>
  </span>
</template>

<script>
  import flvPlayer from './flv-player';

  export default {
    name: 'preview',
    components: {
      flvPlayer: flvPlayer
    },
    mounted: function() {
      //this.$refs.fp.refresh()
    },
    props: [
      'channel_name',
      'listing_id'
    ],
    methods: {
      refresh(){
        this.$refs.fp.refresh()
      },
      pause(){
        this.$refs.fp.pause()
      },
      play() {
        //this.playerOptions.sources.push({ type: "video/x-flv", src: this.preview_url })
        //console.log(this.preview_url)
        //this.$refs.fp.play()
        this.show = true
      },
      previewUrl() {
        let url="/channels/" + encodeURI(this.channel_name) + "/content/" + this.listing_id + "/preview"
        //console.log(url)
        return url
      },
      toggleShow() {
        this.show=!this.show
      },
      // listen event
      onPlayerPlay(player) {
        // console.log('player play!', player)
      },
      onPlayerPause(player) {
        // console.log('player pause!', player)
      },
      // ...player event
 
      // or listen state event
      playerStateChanged(playerCurrentState) {
        // console.log('player current update state', playerCurrentState)
      },
 
      // player is ready
      playerReadied(player) {
        console.log('the player is readied', player)
        // you can use it to do something...
        // player.[methods]
      }
    },
    data: function() {
      return {
        preview_url: "",
        show: false,
        autoplay: true,
        techOrder: ['flash'],
        playerOptions: {
          sources: [],
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .md-dialog {
    max-width: 768px;
  }
</style>
