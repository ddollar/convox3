<template>
  <div class="row">
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
  </div>
</template>

<style lang="scss" src="@/styles/rack.scss"></style>

<script>
import Organization from "@/mixins/Organization";

export default {
  apollo: {
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
    Rack: () => import("@/components/Rack")
  },
  methods: {
    expanded(rid) {
      return this.$route.params.rid == rid && this.$route.meta.expand;
    },
    settings(rid) {
      this.$router.push({ name: "organization/rack/settings", params: { oid: this.organization.id, rid: rid } });
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
