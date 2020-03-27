<template>
  <li class="nav-item" v-if="accessible(to)">
    <router-link class="nav-link" :to="route(to)">
      <i :class="`fa ${icon}`" aria-hidden="true"></i>
      <slot></slot>
    </router-link>
  </li>
</template>

<script>
import { accessible } from "@/scripts/access";

export default {
  apollo: {
    organization: {
      query: require("@/queries/Organization.graphql"),
      variables() {
        return {
          id: this.$route.params.oid
        };
      },
      skip() {
        return !this.$route.params.oid;
      }
    }
  },
  data() {
    return {
      organization: {}
    };
  },
  methods: {
    accessible(name) {
      if (!this.$route.params.oid) {
        return false;
      }

      const { role } = this.$router.resolve(this.route(name)).route.meta;

      return accessible(role, this.organization.access);
    },
    route(name) {
      return {
        name: `organization/${name}`,
        params: { oid: this.$route.params.oid }
      };
    }
  },
  props: ["icon", "to"]
};
</script>
