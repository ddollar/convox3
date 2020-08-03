import { onLogin, onLogout } from "@/vue-apollo";
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters(["authenticated", "token"]),
  },
  methods: {
    login(token) {
      this.$store.commit("setToken", token);
      onLogin(this.$apollo.provider.defaultClient, token);
    },
    logout() {
      this.$store.commit("removeToken");
      onLogout(this.$apollo.provider.defaultClient);
    },
  },
};
