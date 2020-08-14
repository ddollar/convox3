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
          </select>
        </div>
      </b-form-row>
      <b-form-row>
        <div class="form-group col-12">
          <label>Automatic Updates</label>
          <div class="input-group" id="automatic-update">
            <select class="custom-select mr-1" id="update-frequency" name="update-frequency">
              <option value="never">Never</option>
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
  },
  data() {
    return {
      rack: {},
    };
  },
  props: ["rid"],
};
</script>
