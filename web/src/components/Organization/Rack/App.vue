<template>
  <div class="col-12 col-xl-6 col-xxl-4 app clickable">
    <div class="card mb-4 border-bottom-0">
      <div class="card-header d-flex bg-light">
        <div class="flex-grow-1">{{ app.name }}</div>
        <div class="flex-shrink-0">
          <i class="fa fa-check-square"></i>
        </div>
      </div>
      <ul class="list-group list-group-flush">
        <li class="list-group-item d-flex align-items-center p-0">
          <div class="flex-even p-3 border-right">
            <div class="font-weight-bold">Processes</div>
            <i v-if="$apollo.queries.processes.loading" class="spinner"></i>
            <div v-else>{{ processes.length }}</div>
          </div>
          <div class="flex-even p-3 border-right">
            <div class="font-weight-bold">CPU</div>
            <i v-if="$apollo.queries.services.loading" class="spinner"></i>
            <div v-else>{{ cpu }}</div>
          </div>
          <div class="flex-even p-3">
            <div class="font-weight-bold">Memory</div>
            <i v-if="$apollo.queries.services.loading" class="fas fa-circle-notch fa-spin text-muted"></i>
            <div v-else>{{ mem }}</div>
          </div>
        </li>
        <li class="list-group-item p-0">
          <div class="row" style="padding-left: 15px; padding-right: 14px;">
            <div class="col-12 col-xxl-6 p-3 border-right border-bottom bg-light">
              <div
                style="width: 100%; background-color: #fff; height: 80px; border: 1px #eee solid;"
                class="d-flex align-items-center justify-content-center text-secondary"
              >
                CPU/Memory Graph
              </div>
            </div>
            <div class="col-12 col-xxl-6 p-3 border-right border-bottom bg-light">
              <div
                style="width: 100%; background-color: #fff; height: 80px; border: 1px #eee solid;"
                class="d-flex align-items-center justify-content-center text-secondary"
              >
                Network Graph
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
const prettyBytes = require("pretty-bytes");

export default {
  apollo: {
    processes: {
      query: require("@/queries/Organization/Rack/App/Processes.graphql"),
      update: (data) => data.organization?.rack?.app?.processes,
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
      update: (data) => data.organization?.rack?.app?.services,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
          app: this.app.name,
        };
      },
    },
  },
  computed: {
    cpu() {
      return this.services.reduce((ax, s) => ax + s.cpu * s.count, 0);
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
    pretty_memory(num) {
      return prettyBytes(num * 1000000);
    },
  },
  props: ["app"],
};
</script>
