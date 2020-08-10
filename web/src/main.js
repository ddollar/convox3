import Vue from "vue";
import router from "./router";
import store from "./store";

import App from "./App.vue";

import { createProvider } from "./vue-apollo";

import { BootstrapVue, IconsPlugin } from "bootstrap-vue";
import "bootstrap";

Vue.config.productionTip = false;

Vue.use(BootstrapVue);
Vue.use(IconsPlugin);

new Vue({
  render: (h) => h(App),
  router,
  store,
  apolloProvider: createProvider(),
}).$mount("#app");

import VueTimeago from "vue-timeago";

Vue.use(VueTimeago, {
  name: "Timeago",
  locale: "en",
});
