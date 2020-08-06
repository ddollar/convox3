<template>
  <textarea class="bg-dark text-light" style="width: 100%; height: 100%;" v-model="logs"></textarea>
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
          console.log("data", data);
          console.log("dd", this);
          this.data.logs += `${data}\n"`;
        },
      },
    },
  },
  data: function() {
    return {
      logs: "foo\n",
    };
  },
};
</script>
