<template>
  <b-modal :id="`rack-settings-${rid}`" title="Rack Settings" @hide="clear()">
    <div v-if="alert" class="alert alert-danger" role="alert">{{ alert }}</div>
    <b-container fluid>
      <b-form-row class="mb-3">
        <div class="form-group col-7">
          <label for="name">Name</label>
          <input v-model="name" class="form-control" id="name" name="name" maxlength="18" required />
        </div>
        <div class="form-group col-5">
          <label for="runtime">Runtime</label>
          <b-form-select v-model="runtime">
            <b-form-select-option value="">None</b-form-select-option>
            <b-form-select-option v-for="runtime in runtimes" :key="runtime.id" :value="runtime.id">
              {{ runtime.title }}
            </b-form-select-option>
          </b-form-select>
        </div>
      </b-form-row>
      <b-form-row>
        <div class="form-group col-12">
          <label>Automatic Updates</label>
          <div class="input-group" id="automatic-update">
            <b-form-select v-model="update_frequency" :options="options.frequency" />
            <b-form-select v-if="show_day" v-model="update_day" :options="options.day" class="ml-2" />
            <b-form-select v-if="show_hour" v-model="update_hour" :options="options.hour" class="ml-2" />
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
import Error from "@/mixins/Error";

export default {
  apollo: {
    rack: {
      query: require("@/queries/Organization/Rack/Settings.graphql"),
      result(res) {
        const rack = res.data.organization.rack;
        this.name = rack.name;
        this.runtime = rack.runtime || "";
        this.update_day = rack.update_day;
        this.update_frequency = rack.update_frequency;
        this.update_hour = rack.update_hour;
      },
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
        return this.rack?.provider == null;
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
  computed: {
    show_day() {
      switch (this.update_frequency) {
        case "weekly":
          return true;
        default:
          return false;
      }
    },
    show_hour() {
      switch (this.update_frequency) {
        case "weekly":
          return true;
        case "daily":
          return true;
        default:
          return false;
      }
    },
  },
  data() {
    return {
      alert: "",
      name: "",
      runtime: "",
      update_day: 0,
      update_frequency: "never",
      update_hour: 0,
      rack: {},
      options: {
        frequency: [
          { value: "never", text: "Never" },
          { value: "hourly", text: "Hourly" },
          { value: "daily", text: "Daily" },
          { value: "weekly", text: "Weekly" },
        ],
        day: [
          { value: "0", text: "Sunday" },
          { value: "1", text: "Monday" },
          { value: "2", text: "Tuesday" },
          { value: "3", text: "Wednesday" },
          { value: "4", text: "Thursday" },
          { value: "5", text: "Friday" },
          { value: "6", text: "Saturday" },
        ],
        hour: [
          { value: "0", text: "00:00 UTC" },
          { value: "1", text: "01:00 UTC" },
          { value: "2", text: "02:00 UTC" },
          { value: "3", text: "03:00 UTC" },
          { value: "4", text: "04:00 UTC" },
          { value: "5", text: "05:00 UTC" },
          { value: "6", text: "06:00 UTC" },
          { value: "7", text: "07:00 UTC" },
          { value: "8", text: "08:00 UTC" },
          { value: "9", text: "09:00 UTC" },
          { value: "10", text: "10:00 UTC" },
          { value: "11", text: "11:00 UTC" },
          { value: "12", text: "12:00 UTC" },
          { value: "13", text: "13:00 UTC" },
          { value: "14", text: "14:00 UTC" },
          { value: "15", text: "15:00 UTC" },
          { value: "16", text: "16:00 UTC" },
          { value: "17", text: "17:00 UTC" },
          { value: "18", text: "18:00 UTC" },
          { value: "19", text: "19:00 UTC" },
          { value: "20", text: "20:00 UTC" },
          { value: "21", text: "21:00 UTC" },
          { value: "22", text: "22:00 UTC" },
          { value: "23", text: "23:00 UTC" },
        ],
      },
    };
  },
  methods: {
    clear() {
      this.alert = "";
      this.$apollo.queries.rack.refetch();
    },
    remove() {
      this.$bvModal.hide(`rack-settings-${this.rid}`);
      this.$bvModal.show(`rack-remove-${this.rid}`);
    },
    save() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Save.graphql"),
          variables: {
            oid: this.$route.params.oid,
            id: this.rid,
            name: this.name,
            runtime: this.runtime,
            update_day: parseInt(this.update_day),
            update_frequency: this.update_frequency,
            update_hour: parseInt(this.update_hour),
          },
        })
        .then(() => {
          this.$parent.$apollo.queries.rack.refetch();
          this.$bvModal.hide(`rack-settings-${this.rid}`);
        })
        .catch((err) => {
          this.alert = this.graphQLErrors(err);
        });
    },
  },
  mixins: [Error],
  props: ["rid"],
};
</script>
