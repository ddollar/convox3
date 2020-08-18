<template>
  <b-modal :id="id" size="lg" title="Install Rack" @hide="clear()" hide-footer>
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
    <div class="d-flex align-items-center">
      <div class="flex-grow-1">
        <a class="text-success" href="#" @click="parameter_add()">
          <i class="fa fas fa-plus-circle" />
          Add Parameter
        </a>
      </div>
      <div class="flex-shrink-0">
        <button type="submit" class="btn btn-success" @click="submit()">
          <i class="fas fa-plus-circle mr-1" aria-hidden="true"></i>
          Install Rack
        </button>
      </div>
    </div>
    <div v-for="(parameter, index) in parameters" :key="index" class="form-row mt-3">
      <div class="col-6">
        <b-form-select v-model="parameters[index].key">
          <b-form-select-option v-for="parameter in parameters_unused(index)" :key="parameter" :value="parameter">
            {{ parameter }}
          </b-form-select-option>
        </b-form-select>
      </div>
      <div class="col-6 d-flex">
        <div class="flex-grow-1">
          <b-form-input v-model="parameters[index].value" />
        </div>
        <div class="flex-shrink-0 ml-2">
          <b-button variant="danger" @click="parameter_remove(index)">
            <i class="fa fa-times" />
          </b-button>
        </div>
      </div>
    </div>
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
      parameters: [],
      region: "",
      runtime: { engines: [] },
    };
  },
  methods: {
    clear() {
      this.alert = "";
      this.name = "";
      this.engine = "v3";
      this.parameters = [];
      this.region = "";
    },
    parameter_add() {
      if (this.parameters.length < this.runtime.parameters.length) {
        this.parameters.push({ key: "", value: "" });
      }
    },
    parameter_remove(index) {
      this.parameters.splice(index, 1);
    },
    parameters_unused(index) {
      const used = {};
      for (var i in this.parameters) {
        if (i != index) {
          used[this.parameters[i].key] = true;
        }
      }
      return this.runtime.parameters.filter((name) => !used[name]);
    },
    submit() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Install.graphql"),
          variables: {
            oid: this.$route.params.oid,
            iid: this.iid,
            engine: this.engine,
            name: this.name,
            parameters: this.parameters,
            region: this.region,
          },
        })
        .then(() => {
          // this.$bvModal.hide(this.id);
          // this.$parent.$apollo.queries.racks.refetch();
        })
        .catch((err) => {
          console.log("err", err);
          this.alert = this.graphQLErrors(err);
        });
    },
  },
  mixins: [Error],
  props: ["iid"],
};
</script>
