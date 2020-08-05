<template>
  <div class="row">
    <App v-for="app in apps" :key="app.id" :app="app" />
  </div>
</template>

<script>
export default {
  apollo: {
    apps: {
      query: require("@/queries/Apps.graphql"),
      update: (data) => data.organization?.rack?.apps,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        };
      },
    },
  },
  components: {
    App: () => import("@/components/App"),
  },
};
</script>
