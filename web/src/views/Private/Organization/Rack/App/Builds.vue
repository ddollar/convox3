<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">ID</th>
        <th scope="col">Description</th>
        <th scope="col">Started</th>
        <th scope="col">Elapsed</th>
        <th scope="col">Status</th>
        <th scope="col"></th>
      </tr>
    </thead>
    <tbody>
      <Build v-for="build in builds" :key="build.id" :build="build" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    builds: {
      query: require("@/queries/Organization/Rack/App/Builds.graphql"),
      update: data => data.organization?.rack?.app?.builds,
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
    Build: () => import("@/components/Organization/Rack/App/Build"),
  },
};
</script>