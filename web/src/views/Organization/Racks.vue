<template>
  <div class="row">
    <div class="col-12 d-flex mb-4">
      <div class="flex-grow-1">
        <span class="dropdown mr-2">
          <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown">
            <Provider v-if="runtime" :kind="runtime.provider" :text="runtime.title" />
            <span v-else>All Racks</span>
          </button>
          <div class="dropdown-menu">
            <a class="dropdown-item" @click="filter(null)" href="#">All Racks</a>
            <div v-if="integrations.length > 0">
              <div class="dropdown-divider"></div>
              <a v-for="integration in integrations" :key="integration.id" class="dropdown-item" @click="filter(integration)" href="#">
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
              <a v-for="integration in integrations" :key="integration.id" class="dropdown-item" href="#">
                <Provider :kind="integration.provider" :text="integration.title" />
              </a>
            </div>
          </div>
          <a class="btn btn-success ml-2" href="#" data-toggle="modal" data-target="#rack-add-manual">
            <i class="fas fa-plus-circle mr-1"></i>
            Import
          </a>
        </div>
      </div>
    </div>
    <div v-if="filteredRacks.length == 0" class="col-12">
      <div class="card">
        <div class="card-body">No Racks Found</div>
      </div>
    </div>
    <Rack v-else v-for="rack in filteredRacks" :key="rack.id" :rack="rack" />
  </div>
</template>

<style lang="scss" src="@/styles/rack.scss"></style>

<script>
import Organization from "@/mixins/Organization";

export default {
  apollo: {
    integrations: {
      query: require("@/queries/Integrations.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: (data) => data.organization?.integrations,
      variables() {
        return {
          kind: "runtime",
          oid: this.organization.id,
        };
      },
    },
    racks: {
      query: require("@/queries/Racks.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: (data) => data.organization?.racks,
      variables() {
        return {
          oid: this.organization.id,
        };
      },
    },
  },
  components: {
    Provider: () => import("@/components/Provider"),
    Rack: () => import("@/components/Rack"),
  },
  computed: {
    filteredRacks() {
      return (this.racks || []).filter((rack) => {
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
