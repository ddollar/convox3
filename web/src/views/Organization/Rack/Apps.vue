<template>
  <div class="row">
    <div v-if="empty" class="col-12">
      <div class="card">
        <div class="card-body">
          <span>This Rack does not yet have any Apps. Use the</span>
          <span class="text-success font-weight-bold ml-2 mr-1">
            <i class="fas fa-plus-circle"></i>
            Create App
          </span>
          <span>button above to create one.</span>
        </div>
      </div>
    </div>
    <App v-for="app in apps" :key="app.id" :app="app" />
    <!-- <div class="col-12 col-xl-6 col-xxl-4 app create clickable">
      <div class="card mb-4 border-success">
        <div class="card-body d-flex align-items-center justify-content-center">
          <h4 class="text-success mb-0 font-weight-bold">
            <i class="fas fa-plus-circle mr-1"></i>
            Create App
          </h4>
        </div>
      </div>
    </div>-->
  </div>
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
