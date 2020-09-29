<template>
  <div>
    <b-button variant="success" v-b-modal.app-create size="sm">
      <i class="fas fa-plus-circle" />
    </b-button>
    <b-modal id="app-create" title="Create App" @shown="focus()" @hide="clear()">
      <div v-if="alert" class="alert alert-danger" role="alert">{{ alert }}</div>
      <b-container fluid>
        <b-form-row>
          <label>App Name</label>
          <input v-model="name" ref="name" class="form-control" type="text" required pattern="[a-z0-9-]" />
        </b-form-row>
      </b-container>
      <template v-slot:modal-footer>
        <b-button block variant="success" @click="create()">
          <i class="fas fa-plus-circle mr-1" />
          Create App
        </b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import Error from "@/mixins/Error";

export default {
  data() {
    return {
      alert: "",
      name: "",
    };
  },
  methods: {
    clear() {
      this.name = "";
    },
    create() {
      this.alert = "";
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/App/Create.graphql"),
          variables: {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid,
            name: this.name,
          },
        })
        .then(() => {
          this.$bvModal.hide("app-create");
        })
        .catch(err => {
          this.alert = this.graphQLErrors(err);
        });
    },
    focus() {
      this.$refs.name.focus();
    },
  },
  mixins: [Error],
};
</script>
