<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">ID</th>
        <th scope="col">Service</th>
        <th scope="col">Status</th>
        <th scope="col">Release</th>
        <th scope="col">CPU</th>
        <th scope="col">Memory</th>
        <th scope="col">Started</th>
        <th scope="col"></th>
      </tr>
    </thead>
    <tbody>
      <Process v-for="process in processes" :key="process.id" :process="process" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    processes: {
      query: require("@/queries/Organization/Rack/App/Processes.graphql"),
      pollInterval: 5000,
      update: data => data.organization?.rack?.app?.processes,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
          app: this.$route.params.aid,
        };
      },
    },
  },
  components: {
    Process: () => import("@/components/Organization/Rack/App/Process.vue"),
  },
};
</script>
