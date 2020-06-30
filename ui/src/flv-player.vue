<template>
  <video class="flv-player"
    controls></video>
</template>

<script type="text/ecmascript-6">
/* eslint-disable */
import flvjs from 'flv.js'
export default {
  name: 'flv-player',
  props: {
    mediaDataSource: {
      type: Object,
      required: true
    },
    config: {
      type: Object,
      required: false
    },
    autoplay: {
      type: Boolean,
      required: false,
      default: true
    },
    poster: {
      type: String,
      required: false
    },
  },
  data() {
    return {
      isSupported: true,
      flvPlayer: {},
      state: {
        load: false
      }
    }
  },
  // computed: {},
  // watch: {},
  // beforeCreate: function () {},
  created: function () {
    this.isSupported = flvjs.isSupported()
    //console.log(this.mediaDataSource)
    //console.log(this.config)
    this.flvPlayer = flvjs.createPlayer(this.mediaDataSource, this.config)
  },
  // beforeMount: function () {},
  mounted: function () {
    this.flvPlayer.attachMediaElement(this.$el)
    if (this.autoplay) {
      this.play()
    }
  },
  // beforeUpdate: function () {},
  // updated: function () {},
  // beforeDestroy: function () {},
  // destroyed: function () {},
  methods: {
    constructor: function (mediaDataSource, config) {
      this.flvPlayer.constructor(mediaDataSource, config)
    },
    destroy: function () {
      this.flvPlayer.destroy()
    },
    on: function (event, listener) {
      this.flvPlayer.on(event, listener)
    },
    off: function (event, listener) {
      this.flvPlayer.off(event, listener)
    },
    load: function () {
      this.flvPlayer.load()
      this.state.load = true
      if (this.autoplay) {
        this.play()
      }
    },
    unload: function () {
      this.flvPlayer.unload()
      this.state.load = false
    },
    play: function () {
      this.state.load || this.load()
      this.flvPlayer.play()
    },
    pause: function () {
      this.flvPlayer.pause()
    },
    refresh: function () {
      this.pause()
      this.unload()
      this.play()
    },
  }
}
</script>

<style scoped>
</style>
