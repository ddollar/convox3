<template>
  <tr>
    <td>
      <code>{{ instance.id }}</code>
    </td>
    <td><Progress :current="instance.cpu" max="1" color="#0DA542" style="width:100px;" /></td>
    <td><Progress :current="instance.mem" max="1" style="width:100px;" /></td>
    <td>{{ instance.private }}</td>
    <td>{{ instance.public }}</td>
    <td>{{ timeago(instance.started * 1000) }}</td>
    <td>
      <button class="btn btn-danger btn-sm">
        <i class="fa fa-times"></i>
      </button>
    </td>
  </tr>
</template>

<script>
// const prettyBytes = require("pretty-bytes");

export default {
  // apollo: {
  //   processes: {
  //     query: require("@/queries/Organization/Rack/App/Processes.graphql"),
  //     update: data => data.organization?.rack?.app?.processes,
  //     variables() {
  //       return {
  //         oid: this.$route.params.oid,
  //         rid: this.$route.params.rid,
  //         app: this.app.name
  //       };
  //     }
  //   },
  //   services: {
  //     query: require("@/queries/Organization/Rack/App/Services.graphql"),
  //     update: data => data.organization?.rack?.app?.services,
  //     variables() {
  //       return {
  //         oid: this.$route.params.oid,
  //         rid: this.$route.params.rid,
  //         app: this.app.name
  //       };
  //     }
  //   }
  // },
  // computed: {
  //   cpu() {
  //     return this.services.reduce((ax, s) => ax + s.cpu * s.count, 0);
  //   },
  //   mem() {
  //     return this.pretty_memory(
  //       this.services.reduce((ax, s) => ax + s.mem * s.count, 0)
  //     );
  //   }
  // },
  // data() {
  //   return {
  //     processes: [],
  //     services: []
  //   };
  // },
  // methods: {
  //   pretty_memory(num) {
  //     return prettyBytes(num * 1000000);
  //   }
  // },
  components: {
    Progress: () => import("@/components/Progress.vue"),
  },
  methods: {
    percent(value) {
      return `${(value * 100).toFixed(1)}%`;
    },
    timeago(value) {
      return new Date(value);
    },
  },
  props: ["instance"],
};
</script>
