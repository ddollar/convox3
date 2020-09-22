<template>
  <div :class="`col-12 col-xl-6 col-xxl-4 rack clickable ${css}`" @click="goto()">
    <div class="card mb-4 border-bottom-0">
      <div class="card-header d-flex bg-secondary text-light align-items-center">
        <div class="flex-grow-1">
          <i class="fa fa-server mr-1" />
          {{ rack.name }}
        </div>
        <div class="flex-shrink-0">
          <i v-if="$apollo.queries.status.loading" class="spinner"></i>
          <Status v-else :status="status" color />
        </div>
      </div>
      <Installing v-if="installing" :rack="rack" />
      <Uninstalling v-else-if="uninstalling" :rack="rack" />
      <ul v-else class="list-group list-group-flush">
        <li class="list-group-item d-flex align-items-center p-0">
          <div class="flex-even p-3 border-right">
            <div class="font-weight-bold">Apps</div>
            <i v-if="$apollo.queries.apps.loading" class="spinner"></i>
            <div v-else-if="appsError">--</div>
            <div v-else>{{ apps.length }}</div>
          </div>
          <div class="flex-even p-3 border-right">
            <div class="font-weight-bold">CPU</div>
            <i v-if="$apollo.queries.capacity.loading" class="spinner"></i>
            <div v-else-if="capacityError">--</div>
            <div v-else>{{ capacity.cpu.used }} / {{ capacity.cpu.total }}</div>
          </div>
          <div class="flex-even p-3">
            <div class="font-weight-bold">Memory</div>
            <i v-if="$apollo.queries.capacity.loading" class="spinner"></i>
            <div v-else-if="capacityError">--</div>
            <div v-else>
              {{ capacity_bytes(capacity.mem.used) }} /
              {{ capacity_bytes(capacity.mem.total) }}
            </div>
          </div>
        </li>
        <li class="list-group-item p-0">
          <div class="row" style="padding-left: 15px; padding-right: 14px;">
            <div class="col-12 col-xxl-6 p-3 border-right border-bottom bg-light">
              <div
                style="width: 100%; background-color: #fff; height: 80px; border: 1px #eee solid;"
                class="d-flex align-items-center justify-content-center text-secondary"
              >CPU/Memory Graph</div>
            </div>
            <div class="col-12 col-xxl-6 p-3 border-right border-bottom bg-light">
              <div
                style="width: 100%; background-color: #fff; height: 80px; border: 1px #eee solid;"
                class="d-flex align-items-center justify-content-center text-secondary"
              >Network Graph</div>
            </div>
          </div>
        </li>
      </ul>
    </div>
    <Remove :rid="rack.id" />
    <Settings :rid="rack.id" />
  </div>
</template>

<script>
import Organization from "@/mixins/Organization";

const prettyBytes = require("pretty-bytes");

export default {
  apollo: {
    apps: {
      error(error) {
        this.appsError = error;
      },
      query: require("@/queries/Organization/Rack/Apps.graphql"),
      result() {
        this.appsError = null;
      },
      skip() {
        return !this.running;
      },
      update: data => data.organization?.rack?.apps,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.rack.id,
        };
      },
    },
    capacity: {
      error(error) {
        this.capacityError = error;
      },
      query: require("@/queries/Organization/Rack/Capacity.graphql"),
      result() {
        this.capacityError = null;
      },
      skip() {
        return !this.running;
      },
      update: data => data.organization?.rack?.capacity,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.rack.id,
        };
      },
    },
    status: {
      error(error) {
        this.statusError = error;
      },
      pollInterval: 5000,
      query: require("@/queries/Organization/Rack/Status.graphql"),
      result() {
        this.statusError = null;
      },
      update: data => data.organization?.rack?.status,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.rack.id,
        };
      },
    },
  },
  components: {
    Installing: () => import("@/components/Organization/Rack/Installing.vue"),
    Remove: () => import("@/components/Organization/Rack/Remove.vue"),
    Settings: () => import("@/components/Organization/Rack/Settings.vue"),
    Status: () => import("@/components/Organization/Rack/Status.vue"),
    Uninstalling: () => import("@/components/Organization/Rack/Uninstalling.vue"),
  },
  computed: {
    css() {
      return `status-${this.status}`;
    },
    installing() {
      switch (this.status) {
        case "incomplete":
        case "installing":
          return true;
        default:
          return false;
      }
    },
    running() {
      switch (this.status) {
        case "running":
          return true;
        default:
          return false;
      }
    },
    uninstalling() {
      switch (this.status) {
        case "failed":
        case "uninstalling":
          return true;
        default:
          return false;
      }
    },
  },
  data() {
    return {
      appsError: true,
      capacity: {
        cpu: { total: 0, used: 0 },
        mem: { total: 0, used: 0 },
      },
      capacityError: true,
      status: "unknown",
      statusError: null,
    };
  },
  methods: {
    capacity_bytes(num) {
      return prettyBytes(num * 1000000);
    },
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
            name: "organization/rack",
            params: { oid: this.organization.id, rid: this.rack.id },
          });
      }
    },
  },
  mixins: [Organization],
  props: ["rack"],
};
</script>
