export default {
  apollo: {
    organization: {
      query: require("../queries/Organization.graphql"),
      skip() {
        return !this.$route.params.oid;
      },
      variables() {
        return {
          id: this.$route.params.oid,
        };
      },
    },
  },
  data() {
    return {
      organization: {},
    };
  },
};
