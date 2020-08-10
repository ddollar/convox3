<template>
  <div class="card-columns">
    <Integrations :integrations="notifications" title="Notification" />
    <Integrations :integrations="runtimes" title="Runtime" />
    <Integrations :integrations="sources" title="Source" />
  </div>
</template>

<style lang="scss" src="@/styles/rack.scss"></style>

<script>
import Organization from "@/mixins/Organization";

export default {
  apollo: {
    notifications: {
      query: require("@/queries/Organization/Integrations.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: data => data.organization?.integrations,
      variables() {
        return {
          kind: "notification",
          oid: this.organization.id
        };
      }
    },
    runtimes: {
      query: require("@/queries/Organization/Integrations.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: data => data.organization?.integrations,
      variables() {
        return {
          kind: "runtime",
          oid: this.organization.id
        };
      }
    },
    sources: {
      query: require("@/queries/Organization/Integrations.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: data => data.organization?.integrations,
      variables() {
        return {
          kind: "source",
          oid: this.organization.id
        };
      }
    }
  },
  components: {
    Integrations: () => import("@/components/Organization/Integrations")
  },
  mixins: [Organization]
};
</script>