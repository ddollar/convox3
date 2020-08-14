<template>
  <b-modal id="rack-import" title="Import Rack" @hide="clear()">
    <div v-if="alert" class="alert alert-danger" role="alert">{{ alert }}</div>
    <div class="form-row">
      <div class="form-group col-12">
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
    </div>
    <div class="form-row">
      <div class="form-group col-12">
        <label for="hostname">Hostname</label>
        <input v-model="hostname" class="form-control" type="text" id="hostname" required pattern="[a-z0-9.-]+" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-group col-12">
        <label for="hostname">Password</label>
        <input v-model="password" class="form-control" type="password" id="password" required />
      </div>
    </div>
    <template v-slot:modal-footer>
      <button type="submit" class="btn btn-success" @click="submit()">
        <i class="fas fa-plus-circle mr-1" aria-hidden="true"></i>
        Import Rack
      </button>
    </template>
  </b-modal>
</template>

<script>
import Error from "@/mixins/Error";

export default {
  data() {
    return {
      alert: "",
      hostname: "",
      name: "",
      password: "",
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
};
</script>
