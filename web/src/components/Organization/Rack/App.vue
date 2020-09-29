<template>
  <tr :id="id" class="clickable" @click="goto()">
    <td>
      <code>{{ app.name }}</code>
    </td>
    <td>{{ app.status }}</td>
    <td class="text-right"><Async :loading="$apollo.queries.processes.loading" :value="processes.length" /></td>
    <td class="text-right"><Async :loading="$apollo.queries.services.loading" :value="cpu" /></td>
    <td class="text-right"><Async :loading="$apollo.queries.services.loading" :value="mem" /></td>
    <td>
      <b-button variant="danger" size="sm"><i class="fa fa-times" style="font-size:1.1em; padding-left:2px; padding-right:2px;"/></b-button>
    </td>
  </tr>
</template>

<script>
const prettyBytes = require("pretty-bytes");

export default {
  apollo: {
    processes: {
      query: require("@/queries/Organization/Rack/App/Processes.graphql"),
      update: data => data.organization?.rack?.app?.processes,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
          app: this.app.name,
        };
      },
    },
    services: {
      query: require("@/queries/Organization/Rack/App/Services.graphql"),
      update: data => data.organization?.rack?.app?.services,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
          app: this.app.name,
        };
      },
    },
  },
  components: {
    Async: () => import("@/components/Async.vue"),
  },
  computed: {
    cpu() {
      return this.services.reduce((ax, s) => ax + s.cpu * s.count, 0);
    },
    id() {
      return `app-${this._uid}`;
    },
    mem() {
      return this.pretty_memory(this.services.reduce((ax, s) => ax + s.mem * s.count, 0));
    },
  },
  data() {
    return {
      processes: [],
      services: [],
    };
  },
  methods: {
    goto() {
      switch (this.status) {
        case "installing":
        case "uninstalling":
          break;
        case "incomplete":
        case "failed":
        case "unknown":
          this.$bvModal.show(`rack-remove-${this.rack.id}`);
          break;
        default:
          this.$router.push({
            name: "organization/rack/app",
            params: { oid: this.$route.params.oid, rid: this.$route.params.rid, aid: this.app.name },
          });
      }
    },
    pretty_memory(num) {
      return prettyBytes(num * 1000000);
    },
  },
  props: ["app"],
};
</script>
