<template>
  <div class="logs">{{ logs }}</div>
</template>

<script>
export default {
  apollo: {
    $subscribe: {
      logs: {
        query: require("@/queries/Rack/Logs.graphql"),
        variables() {
          return {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid,
          };
        },
        result({ data }) {
          this.logs += `${data.rack_logs.line}\n`;
        },
      },
    },
  },
  data: function() {
    return {
      logs: "",
    };
  },
};
</script>
