<template>
  <div class="rack pb-4">
    <div class="d-flex align-items-center flex-wrap">
      <div class="mb-4 order-0">
        <router-link :to="back()" class="btn btn-secondary">
          <i class="fa fa-chevron-circle-left mr-2"></i>
          <i class="fa fa-server mr-1" />
          <span v-if="rack">{{ rack.name }}</span>
        </router-link>
      </div>
      <div class="flex-fill ml-4 mr-4 mb-4 order-1">
        <h4 class="font-weight-bold mb-0">
          <i class="fa far fa-window-maximize mr-1" />
          {{ $route.params.aid }}
        </h4>
      </div>
      <div class="mr-3 mb-4 order-lg-2" style="font-size: 1.4em; font-weight: bold;">
        <i v-if="$apollo.loading" class="spinner" />
      </div>
      <div class="mb-4 order-4 order-lg-4">
        <nav class="nav nav-pills flex-nowrap">
          <router-link :to="route('services')" class="nav-item nav-link">Services</router-link>
          <router-link :to="route('processes')" class="nav-item nav-link">Processes</router-link>
          <router-link :to="route('builds')" class="nav-item nav-link">Builds</router-link>
          <router-link :to="route('releases')" class="nav-item nav-link">Releases</router-link>
        </nav>
      </div>
    </div>
    <router-view />
  </div>
</template>

<script>
export default {
  apollo: {
    rack: {
      query: require("@/queries/Organization/Rack.graphql"),
      update: data => data.organization?.rack,
      variables() {
        return {
          oid: this.$route.params.oid,
          id: this.$route.params.rid,
        };
      },
    },
  },
  methods: {
    back() {
      return {
        name: "organization/rack/apps",
        params: {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        },
      };
    },
    route(page) {
      return {
        name: `organization/rack/app/${page}`,
        params: {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
          aid: this.$route.params.aid,
        },
      };
    },
  },
};
</script>