<template>
  <table class="table table-light">
    <thead class="thead-dark" style="border-top-left-radius: 4px;">
      <tr>
        <th scope="col" class="expand">ID</th>
        <th scope="col">CPU</th>
        <th scope="col">Memory</th>
        <th scope="col">Private</th>
        <th scope="col">Public</th>
        <th scope="col">Launched</th>
        <th scope="col"></th>
      </tr>
    </thead>
    <tbody>
      <Instance v-for="instance in instances" :key="instance.id" :instance="instance" />
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
    instances: {
      query: require("@/queries/Organization/Rack/Instances.graphql"),
      update: (data) => data.organization?.rack?.instances,
      variables() {
        return {
          oid: this.$route.params.oid,
          rid: this.$route.params.rid,
        };
      },
    },
  },
  components: {
    Instance: () => import("@/components/Organization/Rack/Instance"),
  },
};
</script>
