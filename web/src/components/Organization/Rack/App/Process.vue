<template>
  <tr :id="id">
    <td>
      <code>{{ process.id }}</code>
    </td>
    <td>{{ process.service }}</td>
    <td>{{ process.status }}</td>
    <td>
      <code>{{ process.release }}</code>
    </td>
    <td>
      <Progress :current="process.cpu" color="#0DA542" style="width:100px;" />
    </td>
    <td>
      <Progress :current="process.mem" style="width:100px;" />
    </td>
    <td>
      <RelativeTime :time="process.started" />
    </td>
    <td>
      <button class="btn btn-danger btn-sm" @click="stop(process.app, process.id, $event)">
        <i class="fa fa-times"></i>
      </button>
    </td>
  </tr>
</template>

<script>
export default {
  components: {
    Progress: () => import("@/components/Progress.vue"),
    RelativeTime: () => import("@/components/RelativeTime.vue"),
  },
  computed: {
    id() {
      return `process-${this._uid}`;
    },
    percent(value) {
      return `${(value * 100).toFixed(1)}%`;
    },
  },
  methods: {
    stop(app, id, event) {
      event.target.disabled = true;
      this.$apollo
        .mutate({
          mutation: require("@/queries/Organization/Rack/Process/Stop.graphql"),
          variables: {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid,
            app: app,
            pid: id,
          },
        })
        .then(() => {
          this.$parent.$apollo.queries.processes.refresh();
        })
        .catch(() => {
          event.target.disabled = false;
        });
    },
  },
  props: ["process"],
};
</script>
