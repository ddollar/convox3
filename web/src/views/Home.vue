<template>
  <div class="home"></div>
</template>

<script>
export default {
  async mounted() {
    try {
      var orgs = (
        await this.$apollo.query({
          query: require("../queries/Organizations.graphql"),
        })
      ).data.organizations;
      if (orgs.length > 0) {
        this.$router.push({
          name: "organization/racks",
          params: { oid: orgs[0].id },
        });
      }
    } catch (err) {
      this.$router.push({
        name: "login",
      });
    }
  },
};
</script>
