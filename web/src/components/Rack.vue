<template>
  <li class="list-group-item bg-light">
    <nav>
      <div class="nav nav-tabs">
        <router-link class="nav-item nav-link" :to="tab('apps')">Apps</router-link>
        <router-link class="nav-item nav-link" :to="tab('instances')">Instances</router-link>
        <router-link class="nav-item nav-link" :to="tab('resources')">Resources</router-link>
        <router-link class="nav-item nav-link" :to="tab('updates')">Updates</router-link>
      </div>
      <div class="tab-content border border-top-0">
        <Apps v-if="active('apps')" :oid="oid" :rid="id" />
      </div>
    </nav>
  </li>
</template>

<script>
export default {
  apollo: {
    rack: {
      query: require("@/queries/Rack.graphql"),
      skip() {
        return !this.id;
      },
      update: data => data.organization?.racks,
      variables() {
        return {
          oid: this.oid,
          id: this.id
        };
      }
    }
  },
  components: {
    Apps: () => import("@/components/Apps")
  },
  data() {
    return {
      rack: {}
    };
  },
  methods: {
    active(name) {
      return this.$route.name == `organization/rack/${name}`;
    },
    tab(name) {
      return {
        name: `organization/rack/${name}`,
        params: { oid: this.oid, rid: this.id }
      };
    }
  },
  props: ["id", "oid"]
};
</script>