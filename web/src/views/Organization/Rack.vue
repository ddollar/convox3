<template>
  <div class="rack pb-4">
    <div class="d-flex align-items-center flex-wrap">
      <div class="mb-4 order-0">
        <router-link :to="back()" class="btn btn-dark">
          <i class="fa fa-chevron-circle-left mr-1"></i>
          All Racks
        </router-link>
      </div>
      <div class="flex-fill ml-4 mr-4 mb-4 order-1">
        <h4 class="font-weight-bold mb-0">{{ rack.name }}</h4>
      </div>
      <div class="mr-4 mb-4 order-3 order-lg-2">
        <nav class="nav nav-pills flex-nowrap">
          <router-link :to="route('apps')" class="nav-item nav-link">Apps</router-link>
          <router-link :to="route('instances')" class="nav-item nav-link">Instances</router-link>
          <router-link :to="route('logs')" class="nav-item nav-link">Logs</router-link>
          <router-link :to="route('processes')" class="nav-item nav-link">Processes</router-link>
          <router-link :to="route('resources')" class="nav-item nav-link">Resources</router-link>
          <router-link :to="route('updates')" class="nav-item nav-link">Updates</router-link>
        </nav>
      </div>
      <div class="mb-4 order-2 order-lg-3">
        <button class="btn btn-secondary">
          <i class="fa fa-cog"></i>
        </button>
      </div>
    </div>

    <router-view />
  </div>
</template>

<script>
import Organization from "@/mixins/Organization";
import Rack from "@/mixins/Rack";

export default {
  methods: {
    back() {
      return {
        name: "organization/racks",
        params: { oid: this.$route.params.oid },
      };
    },
    route(page) {
      return {
        name: `organization/rack/${page}`,
        params: { oid: this.$route.params.oid, rid: this.$route.params.rid },
      };
    },
  },
  mixins: [Organization, Rack],
  mounted() {
    if (this.$route.name == "organization/rack") {
      this.$router.push({
        name: "organization/rack/apps",
        params: { oid: this.$route.params.oid, rid: this.$route.params.rid },
      });
    }
  },
};
</script>
