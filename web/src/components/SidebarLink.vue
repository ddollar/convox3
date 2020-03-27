<template>
  <li class="nav-item" v-if="accessible(to)">
    <router-link class="nav-link" :to="route(to)">
      <i :class="`fa ${icon}`" aria-hidden="true"></i>
      <slot></slot>
    </router-link>
  </li>
</template>

<script>
import Organization from "@/mixins/Organization";
import { accessible } from "@/scripts/access";

export default {
  methods: {
    accessible(name) {
      if (!this.organization.id) {
        return false;
      }

      const { role } = this.$router.resolve(this.route(name)).route.meta;

      return accessible(role, this.organization.access);
    },
    route(name) {
      return {
        name: `organization/${name}`,
        params: { oid: this.organization.id }
      };
    }
  },
  mixins: [Organization],
  props: ["icon", "to"]
};
</script>
