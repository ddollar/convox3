<template>
  <b-modal :id="id" size="lg" title="Install Rack" @hide="clear()">
    <div v-if="alert" class="alert alert-danger" role="alert">{{ alert }}</div>
    <div class="form-row">
      <div class="form-group col-8 d-flex">
        <div class="flex-grow-1">
          <label for="name">Name</label>
          <input
            v-model="name"
            class="form-control"
            type="text"
            id="name"
            placeholder="production"
            required
            minlength="3"
            maxlength="18"
            pattern="[a-z0-9-]+"
          />
        </div>
        <div class="flex-shrink-0 ml-2" v-if="runtime.engines.length > 0">
          <label>Engine</label>
          <b-form-select v-model="engine" :options="runtime.engines" value-field="name" text-field="description"></b-form-select>
        </div>
      </div>
      <div class="form-group col-4">
        <label>Region</label>
        <b-form-select v-model="region">
          <b-form-select-option v-for="region in runtime.regions" :key="region" :value="region">{{ region }}</b-form-select-option>
        </b-form-select>
      </div>
    </div>
    <template v-slot:modal-footer>
      <button type="submit" class="btn btn-success" @click="submit()">
        <i class="fas fa-plus-circle mr-1" aria-hidden="true"></i>
        Install Rack
      </button>
    </template>
  </b-modal>
</template>

<script>
import Error from "@/mixins/Error";

export default {
  apollo: {
    runtime: {
      query: require("@/queries/Organization/Runtime.graphql"),
      update: (data) => data.organization?.runtime,
      variables() {
        return {
          oid: this.$route.params.oid,
          id: this.iid,
        };
      },
    },
  },
  computed: {
    id() {
      return `rack-install-${this.iid}`;
    },
  },
  data() {
    return {
      alert: "",
      engine: "v3",
      name: "",
      region: "",
      runtime: { engines: [] },
    };
  },
  methods: {
    clear() {
      this.alert = "";
    },
    submit() {
      this.alert = "";
      const { name, hostname, password } = this.$data;
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Import.graphql"),
          variables: {
            oid: this.$route.params.oid,
            name,
            hostname,
            password,
          },
        })
        .then(() => {
          this.$bvModal.hide("rack-import");
          this.$parent.$apollo.queries.racks.refetch();
        })
        .catch((err) => {
          this.alert = this.graphQLErrors(err);
        });
    },
  },
  mixins: [Error],
  props: ["iid"],
};
</script>
