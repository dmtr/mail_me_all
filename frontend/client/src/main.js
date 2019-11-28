import "regenerator-runtime/runtime";
import Vue from "vue";
import VueRouter from "vue-router";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";

import store from "./store";
import Settings from "./components/Settings";
import Home from "./components/Home";

Vue.config.productionTip = false;
Vue.use(VueRouter);

const routes = [
  { path: "/settings", component: Settings },
  { path: "/", component: Home }
];
const router = new VueRouter({ routes });

new Vue({
  store,
  vuetify,
  router,
  render: h => h(App)
}).$mount("#app");
