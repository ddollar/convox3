<template>
  <b-modal :id="`rack-settings-${rid}`" title="Rack Settings">
    <b-container fluid>
      <b-form-row class="mb-3">
        <div class="form-group col-7">
          <label for="name">Name</label>
          <input class="form-control" id="name" name="name" maxlength="18" required :value="rack.name" />
        </div>
        <div class="form-group col-5">
          <label for="runtime">Runtime</label>
          <select class="form-control" name="runtime" id="runtime">
            <option value="">None</option>
            <option v-for="runtime in runtimes" :key="runtime.id" :value="runtime.id">{{ runtime.title }}</option>
          </select>
        </div>
      </b-form-row>
      <b-form-row>
        <div class="form-group col-12">
          <label>Automatic Updates</label>
          <div class="input-group" id="automatic-update">
            <select class="custom-select mr-1" id="update-frequency" name="update-frequency">
              <option value="never">Never</option>
              <option value="hourly">Hourly</option>
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
            </select>
            <select class="custom-select ml-1 mr-1" id="update-day" name="update-day">
              <option value="0">Sunday</option>
            </select>
            <select class="custom-select ml-1" id="update-hour" name="update-hour">
              <option value="0">00:00 UTC</option>
            </select>
          </div>
        </div>
      </b-form-row>
    </b-container>
    <template v-slot:modal-footer class="d-flex align-items-center">
      <div class="flex-grow-1">
        <button class="btn btn-danger" @click="remove()">
          <i class="fa fa-times mr-1"></i>
          Remove Rack
        </button>
      </div>
      <div class="flex-shrink-0">
        <button type="submit" class="btn btn-primary" @click="save()">
          <i class="fas fa-check mr-1" aria-hidden="true"></i>
          Save Changes
        </button>
      </div>
    </template>
  </b-modal>
</template>

<script>
export default {
  apollo: {
    rack: {
      query: require("@/queries/Organization/Rack/Settings.graphql"),
      update: (data) => data.organization?.rack,
      variables() {
        return {
          oid: this.$route.params.oid,
          id: this.rid,
        };
      },
    },
    runtimes: {
      query: require("@/queries/Organization/Rack/Runtimes.graphql"),
      skip() {
        return this.rack.provider == null;
      },
      update: (data) => data.organization?.integrations,
      variables() {
        return {
          oid: this.$route.params.oid,
          provider: this.rack.provider,
        };
      },
    },
  },
  data() {
    return {
      rack: {},
    };
  },
  methods: {
    remove() {
      this.$bvModal.hide(`rack-settings-${this.rid}`);
      this.$bvModal.show(`rack-remove-${this.rid}`);
    },
    save() {
      this.$bvModal.hide(`rack-settings-${this.rid}`);
    },
  },
  props: ["rid"],
};
</script>
