import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    token: localStorage.getItem("token"),
  },
  getters: {
    authenticated: state => state.token !== null,
    token: state => state.token,
  },
  mutations: {
    removeToken(state) {
      state.token = undefined;
      localStorage.removeItem("token");
    },
    setToken(state, value) {
      localStorage.setItem("token", value);
      state.token = value;
    },
  },
  actions: {},
  modules: {},
});
