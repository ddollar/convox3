<template>
  <div class="row">
    <div class="col-12 d-flex mb-4">
      <div class="flex-grow-1">
        <span class="dropdown mr-2">
          <button
            class="btn btn-secondary dropdown-toggle"
            type="button"
            id="dropdownMenuButton"
            data-toggle="dropdown"
          >
            <Provider v-if="runtime" :kind="runtime.provider" :text="runtime.title" />
            <span v-else>All Racks</span>
          </button>
          <div class="dropdown-menu">
            <a class="dropdown-item" @click="filter(null)" href="#">All Racks</a>
            <div v-if="integrations.length > 0">
              <div class="dropdown-divider"></div>
              <a
                v-for="integration in integrations"
                :key="integration.id"
                class="dropdown-item"
                @click="filter(integration)"
                href="#"
              >
                <Provider :kind="integration.provider" :text="integration.title" />
              </a>
            </div>
          </div>
        </span>
      </div>
      <div class="flex-shrink-0">
        <div class="dropdown">
          <button class="btn btn-success dropdown-toggle" type="button" data-toggle="dropdown">
            <i class="fas fa-cloud-upload-alt mr-1" aria-hidden="true"></i>
            Install
          </button>
          <div class="dropdown-menu dropdown-menu-right">
            <div v-if="integrations.length == 0" class="dropdown-item">No Runtime Integrations</div>
            <div v-else>
              <b-dropdown-item
                v-for="integration in integrations"
                :key="integration.id"
                @click="install(integration.id)"
              >
                <Provider :kind="integration.provider" :text="integration.title" />
                <Install :iid="integration.id" />
              </b-dropdown-item>
            </div>
          </div>
          <b-button v-b-modal.rack-import variant="success" class="ml-2">
            <i class="fas fa-plus-circle mr-1"></i>
            Import
          </b-button>
          <Import />
        </div>
      </div>
    </div>
    <div v-if="filteredRacks.length == 0" class="col-12">
      <div class="card">
        <div class="card-body">
          <span>This organization does not yet have any Racks. Use the</span>
          <span class="text-success font-weight-bold ml-1">
            <i class="fa fa-cloud-upload-alt"></i>
            Install
          </span>
          <span>or</span>
          <span class="text-success font-weight-bold ml-1">
            <i class="fas fa-plus-circle"></i>
            Import
          </span>
          <span>buttons above to add one.</span>
        </div>
      </div>
    </div>
    <Rack v-else v-for="rack in filteredRacks" :key="rack.id" :rack="rack" />
  </div>
</template>

<script>
import Organization from "@/mixins/Organization";

export default {
  apollo: {
    integrations: {
      query: require("@/queries/Organization/Integrations.graphql"),
      update: data => data.organization?.integrations,
      variables() {
        return {
          kind: "runtime",
          oid: this.$route.params.oid,
        };
      },
    },
    racks: {
      query: require("@/queries/Organization/Racks.graphql"),
      pollInterval: 5000,
      update: data => data.organization?.racks,
      variables() {
        return {
          oid: this.$route.params.oid,
        };
      },
    },
  },
  components: {
    Import: () => import("@/components/Organization/Rack/Import.vue"),
    Install: () => import("@/components/Organization/Rack/Install.vue"),
    Provider: () => import("@/components/Provider.vue"),
    Rack: () => import("@/components/Organization/Rack.vue"),
  },
  computed: {
    filteredRacks() {
      return (this.racks || []).filter(rack => {
        return this.runtime === null ? true : rack.runtime == this.runtime.id;
      });
    },
  },
  data: function() {
    return {
      integrations: [],
      runtime: null,
    };
  },
  methods: {
    expanded(rid) {
      return this.$route.params.rid == rid && this.$route.meta.expand;
    },
    filter(iid) {
      this.runtime = iid;
    },
    install(id) {
      this.$bvModal.show(`rack-install-${id}`);
    },
    settings(rid) {
      this.$router.push({
        name: "organization/rack/settings",
        params: { oid: this.organization.id, rid: rid },
      });
    },
    toggle(rid) {
      if (this.expanded(rid)) {
        this.$router.push({
          name: "organization/racks",
          params: { oid: this.organization.id },
        });
      } else {
        this.$router.push({
          name: "organization/rack/apps",
          params: { oid: this.organization.id, rid: rid },
        });
      }
    },
  },
  mixins: [Organization],
};
</script>
