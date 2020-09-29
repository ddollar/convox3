<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">Resource</th>
        <th scope="col">Type</th>
        <th scope="col">Status</th>
        <th scope="col"></th>
      </tr>
    </thead>
    <tbody>
      <Resource v-for="resource in resources" :key="resource.name" :resource="resource" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    resources: {
      query: require("@/queries/Organization/Rack/Resources.graphql"),
      pollInterval: 5000,
      update: data => data.organization?.rack?.resources,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        };
      },
    },
  },
  components: {
    Resource: () => import("@/components/Organization/Rack/Resource.vue"),
  },
};
</script>
