import Vue from "vue";
import App from "./App.vue";
import VueMaterial from 'vue-material'
//import 'vue-material/dist/vue-material.min.css'
//import 'vue-material/dist/theme/default.css'
import axios from 'axios'
import moment from 'moment'

Vue.prototype.$http = axios
Vue.prototype.$baseurl = "http://127.0.0.1:8080/"
Vue.prototype.$buildurl = function(url) {
    return this.$baseurl + url;
}

Vue.config.productionTip = false;
Vue.use(VueMaterial)
Vue.prototype.moment = moment

new Vue({
  render: h => h(App)
}).$mount("#app-placeholder");
