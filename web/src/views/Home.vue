<template>
  <div class="home"></div>
</template>

<script>
export default {
  async mounted() {
    var orgs = (
      await this.$apollo.query({
        query: require("../queries/Organizations.graphql")
      })
    ).data.organizations;
    if (orgs.length > 0) {
      this.$router.push({
        name: "organization/dashboard",
        params: { oid: orgs[0].id }
      });
    }
  }
};
</script>
