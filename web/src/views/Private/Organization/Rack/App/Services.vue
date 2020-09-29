<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">Service</th>
        <th scope="col">Count</th>
        <th scope="col">CPU</th>
        <th scope="col">Memory</th>
      </tr>
    </thead>
    <tbody>
      <Service v-for="service in services" :key="service.id" :service="service" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    services: {
      query: require("@/queries/Organization/Rack/App/Services.graphql"),
      update: data => data.organization?.rack?.app?.services,
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
    Service: () => import("@/components/Organization/Rack/App/Service"),
  },
};
</script>