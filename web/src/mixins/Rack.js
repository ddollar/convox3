export default {
  apollo: {
    rack: {
      query: require("../queries/Rack.graphql"),
      skip() {
        return !(this.$route.params.oid && this.$route.params.rid);
      },
      update: (data) => data.organization?.rack,
      variables() {
        return {
          id: this.$route.params.rid,
          oid: this.$route.params.oid,
        };
      },
    },
  },
  data() {
    return {
      rack: {},
    };
  },
};
