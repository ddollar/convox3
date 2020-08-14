<template>
  <b-modal :id="`rack-remove-${rid}`" title="Remove Rack" header-bg-variant="danger" header-text-variant="light">
    Lorem ipsum sit dolor?
    <template v-slot:modal-footer>
      <button class="btn btn-danger" @click="remove()">
        <i class="fa fa-times mr-1"></i>
        Remove Rack
      </button>
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
  },
  data() {
    return {
      rack: {},
    };
  },
  methods: {
    remove() {
      this.$bvModal.hide(`rack-remove-${this.rid}`);
    },
  },
  props: ["rid"],
};
</script>
