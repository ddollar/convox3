<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">App</th>
        <th scope="col">Status</th>
        <th scope="col" class="text-right">Processes</th>
        <th scope="col" class="text-right">CPU</th>
        <th scope="col" class="text-right">Memory</th>
        <th scope="col">
          <Create />
        </th>
      </tr>
    </thead>
    <tbody>
      <App v-for="app in apps" :key="app.id" :app="app" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    apps: {
      query: require("@/queries/Organization/Rack/Apps.graphql"),
      pollInterval: 5000,
      update: data => data.organization?.rack?.apps,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        };
      },
    },
  },
  components: {
    App: () => import("@/components/Organization/Rack/App"),
    Create: () => import("@/components/Organization/Rack/App/Create"),
  },
  computed: {
    empty() {
      if (this.$apollo.queries.apps.loading) return false;
      if (this.apps.length > 0) return false;
      return true;
    },
  },
  data() {
    return {
      apps: [],
    };
  },
};
</script>
