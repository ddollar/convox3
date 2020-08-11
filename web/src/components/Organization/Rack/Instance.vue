<template>
  <tr :id="id">
    <td>
      <code>{{ instance.id }}</code>
    </td>
    <td>{{ instance.status }}</td>
    <td><Progress :current="instance.cpu" color="#0DA542" style="width:100px;" /></td>
    <td><Progress :current="instance.mem" style="width:100px;" /></td>
    <td>{{ instance.private }}</td>
    <td>{{ instance.public }}</td>
    <td><Timeago v-if="instance.started > 0" :datetime="datetime" /></td>
    <td>
      <button class="btn btn-danger btn-sm" @click="terminate(instance.id, $event)">
        <i class="fa fa-times"></i>
      </button>
    </td>
  </tr>
</template>

<script>
export default {
  components: {
    Progress: () => import("@/components/Progress.vue"),
  },
  computed: {
    datetime() {
      return new Date(this.instance.started * 1000);
    },
    id() {
      return `instance-${this._uid}`;
    },
    percent(value) {
      return `${(value * 100).toFixed(1)}%`;
    },
  },
  methods: {
    terminate(id, event) {
      event.target.disabled = true;
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Instance/Terminate.graphql"),
          variables: {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid,
            iid: id,
          },
        })
        .then(() => {
          this.$parent.$apollo.queries.instances.refresh();
        })
        .catch(() => {
          event.target.disabled = false;
        });
    },
  },
  props: ["instance"],
};
</script>
