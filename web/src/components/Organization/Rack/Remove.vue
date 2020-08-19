<template>
  <b-modal :id="`rack-remove-${rid}`" :title="title" header-bg-variant="danger" header-text-variant="light">
    <div v-if="alert" class="alert alert-danger" role="alert">{{ alert }}</div>
    <div>
      Remove
      <strong>{{ rack.name }}</strong> from this organization?
    </div>
    <div v-if="rack.uninstallable" class="mt-3 text-danger font-weight-bold">
      The underlying infrastructure will be destroyed.
    </div>
    <template v-slot:modal-footer>
      <button class="btn btn-danger" @click="remove()">
        <i class="fa fa-times mr-1"></i>
        {{ title }}
      </button>
    </template>
  </b-modal>
</template>

<script>
import Error from "@/mixins/Error";

export default {
  apollo: {
    rack: {
      query: require("@/queries/Organization/Rack/Settings.graphql"),
      update: data => data.organization?.rack,
      variables() {
        return {
          oid: this.$route.params.oid,
          id: this.rid,
        };
      },
    },
  },
  computed: {
    title() {
      return this.rack.uninstallable ? "Uninstall Rack" : "Remove Rack";
    },
  },
  data() {
    return {
      alert: "",
      rack: {},
    };
  },
  methods: {
    remove() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Remove.graphql"),
          variables: {
            oid: this.$route.params.oid,
            id: this.rid,
          },
        })
        .then(() => {
          switch (this.$route.name) {
            case "organization/racks":
              this.$parent.$apollo.queries.racks?.refetch();
              this.$bvModal.hide(`rack-remove-${this.rid}`);
              break;
            default:
              this.$router.replace({
                name: "organization/racks",
                params: { oid: this.$route.params.oid },
              });
          }
        })
        .catch(err => {
          this.alert = this.graphQLErrors(err);
        });
    },
  },
  mixins: [Error],
  props: ["rid"],
};
</script>
