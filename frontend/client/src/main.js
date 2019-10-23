import "regenerator-runtime/runtime";
import Vue from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import axios from "axios";
import VueAxios from "vue-axios";
import VueAuthenticate from "vue-authenticate";

import store from "./store";

Vue.config.productionTip = false;

Vue.use(VueAxios, axios);

Vue.use(VueAuthenticate, {
  baseUrl: "https://localhost",

  providers: {
    github: {
      redirectUri: "https://localhost:8080/tw/callback"
    }
  }
});

new Vue({
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
