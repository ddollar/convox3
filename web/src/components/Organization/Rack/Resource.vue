<template>
  <tr :id="id">
    <td>
      <code>{{ resource.name }}</code>
    </td>
    <td>{{ resource.type }}</td>
    <td>{{ resource.status }}</td>
    <td>
      <button class="btn btn-danger btn-sm" @click="del(resource.name)">
        <i class="fa fa-times"></i>
      </button>
    </td>
  </tr>
</template>

<script>
export default {
  computed: {
    id() {
      return `resource-${this._uid}`;
    },
  },
  methods: {
    del(name, event) {
      event.target.disabled = true;
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Resource/Delete.graphql"),
          variables: {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid,
            name: name,
          },
        })
        .then(() => {
          this.$parent.$apollo.queries.resources.refresh();
        })
        .catch(() => {
          event.target.disabled = false;
        });
    },
  },
  props: ["resource"],
};
</script>
