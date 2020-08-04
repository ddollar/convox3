<template>
  <div class="rack">
    <div class="d-flex align-items-center flex-wrap">
      <div class="mb-4 order-0">
        <router-link :to="back()" class="btn btn-dark">
          <i class="fa fa-chevron-circle-left mr-1"></i>
          All Racks
        </router-link>
      </div>
      <div class="flex-fill ml-4 mr-4 mb-4 order-1">
        <h4 class="font-weight-bold mb-0">production</h4>
      </div>
      <div class="mr-4 mb-4 order-3 order-lg-2">
        <nav class="nav nav-pills flex-nowrap">
          <a class="nav-item nav-link active" href="#">Dashboard</a>
          <a class="nav-item nav-link" href="#">Apps</a>
          <a class="nav-item nav-link" href="#">Instances</a>
          <a class="nav-item nav-link" href="#">Logs</a>
        </nav>
      </div>
      <div class="mb-4 order-2 order-lg-3">
        <button class="btn btn-secondary">
          <i class="fa fa-cog"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import Organization from "@/mixins/Organization";

export default {
  methods: {
    back() {
      return {
        name: "organization/racks",
        params: { oid: this.organization.id },
      };
    },
  },
  mixins: [Organization],
};
</script>
