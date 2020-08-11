<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col" class="expand">ID</th>
        <th scope="col">App</th>
        <th scope="col">Service</th>
        <th scope="col">Status</th>
        <th scope="col">Release</th>
        <th scope="col">CPU</th>
        <th scope="col">Memory</th>
        <th scope="col">Started</th>
        <th scope="col"></th>
      </tr>
    </thead>
    <tbody>
      <Process v-for="process in processes" :key="process.id" :process="process" />
    </tbody>
  </table>
  <!-- <div class="row">
    <div class="col-4">
      <div class="card">
        <div class="card-header">Foo</div>
        <ul class="list-group list-group-flush">
          <Instance v-for="instance in instances" :key="instance.id" :instance="instance" />
        </ul>
      </div>
    </div>
  </div>-->
</template>

<script>
export default {
  apollo: {
    processes: {
      query: require("@/queries/Organization/Rack/Processes.graphql"),
      pollInterval: 5000,
      update: (data) => data.organization?.rack?.processes,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        };
      },
    },
  },
  components: {
    Process: () => import("@/components/Organization/Rack/Process.vue"),
  },
};
</script>
