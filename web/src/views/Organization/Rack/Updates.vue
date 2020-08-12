<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col">Version</th>
        <th scope="col" class="expand">Params</th>
        <th scope="col">Started</th>
      </tr>
    </thead>
    <tbody>
      <Update v-for="update in updates" :key="update.id" :update="update" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    updates: {
      query: require("@/queries/Organization/Rack/Updates.graphql"),
      pollInterval: 5000,
      update: (data) => data.organization?.rack?.updates,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        };
      },
    },
  },
  components: {
    Update: () => import("@/components/Organization/Rack/Update.vue"),
  },
};
</script>
