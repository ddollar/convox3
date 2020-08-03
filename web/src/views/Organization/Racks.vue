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
            aria-haspopup="true"
            aria-expanded="false"
          >
            <Provider v-if="runtime" :kind="runtime.provider" :text="runtime.title" />
            <span v-else>All Racks</span>
          </button>
          <div class="dropdown-menu">
            <a class="dropdown-item" @click="filter(null)" href="#">All Racks</a>

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
        </span>
      </div>
      <div class="flex-shrink-0">
        <div class="dropdown">
          <button class="btn btn-success dropdown-toggle" type="button" data-toggle="dropdown">
            <i class="fas fa-cloud-upload-alt mr-1" aria-hidden="true"></i>
            Install
          </button>
          <div class="dropdown-menu dropdown-menu-right">
            <a
              v-for="integration in integrations"
              :key="integration.id"
              class="dropdown-item"
              href="#"
            >
              <!-- data-toggle="modal.lazy"
              data-target="#runtime-1df0909d-f5bb-4eba-8df7-6e72655086b3"-->
              <Provider :kind="integration.provider" :text="integration.title" />
            </a>
          </div>

          <a
            class="btn btn-success ml-2"
            href="#"
            data-toggle="modal"
            data-target="#rack-add-manual"
          >
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
  <!-- <div class="row">
    <div class="col-12 col-xxl-8 offset-xxl-2">
      <div id="racks" class="card">
        <ul class="list-group list-group-flush">
          <div v-for="rack in racks" :key="rack.id">
            <li class="list-group-item rack" @click="toggle(rack.id)">
              <div class="flex-grow-1">
                <h5 class="mb-0 text-dark">
                  <i class="fas fa-check-square text-success mr-3"></i>
                  <span style="font-weight:650;letter-spacing:0.01em;">{{ rack.name }}</span>
                </h5>
              </div>
              <button class="btn btn-primary btn-sm ml-2" v-if="!expanded(rack.id)" @click.stop="settings(rack.id)">
                <i class="fa fa-cog"></i>
              </button>
            </li>
            <Rack v-if="expanded(rack.id)" :id="rack.id" :oid="organization.id" />
          </div>
        </ul>
      </div>
    </div>
  </div>-->
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
      update: data =>
        data.organization?.integrations.sort((a, b) => {
          return (
            a.provider.localeCompare(b.provider) ||
            a.title.localeCompare(b.title)
          );
        }),
      variables() {
        return {
          kind: "runtime",
          oid: this.organization.id
        };
      }
    },
    racks: {
      query: require("@/queries/Racks.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: data => data.organization?.racks,
      variables() {
        return {
          oid: this.organization.id
        };
      }
    }
  },
  components: {
    Provider: () => import("@/components/Provider"),
    Rack: () => import("@/components/Rack")
  },
  computed: {
    filteredRacks() {
      return this.racks?.filter(rack => {
        return this.runtime === null ? true : rack.runtime == this.runtime.id;
      });
    }
  },
  data: function() {
    return {
      runtime: null
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
        params: { oid: this.organization.id, rid: rid }
      });
    },
    toggle(rid) {
      if (this.expanded(rid)) {
        this.$router.push({
          name: "organization/racks",
          params: { oid: this.organization.id }
        });
      } else {
        this.$router.push({
          name: "organization/rack/apps",
          params: { oid: this.organization.id, rid: rid }
        });
      }
    }
  },
  mixins: [Organization]
};
</script>
