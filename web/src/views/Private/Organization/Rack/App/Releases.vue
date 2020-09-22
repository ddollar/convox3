<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">ID</th>
        <th scope="col">Description</th>
        <th scope="col">Build</th>
        <th scope="col">Created</th>
        <th scope="col"></th>
      </tr>
    </thead>
    <tbody>
      <Release v-for="release in releases" :key="release.id" :release="release" />
    </tbody>
  </table>
</template>

<script>
export default {
  apollo: {
    releases: {
      query: require("@/queries/Organization/Rack/App/Releases.graphql"),
      update: data => data.organization?.rack?.app?.releases,
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
    Release: () => import("@/components/Organization/Rack/App/Release"),
  },
};
</script>