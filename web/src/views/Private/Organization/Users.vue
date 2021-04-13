<template>
  <div class="row">
    <div class="col-12 col-xxl-7">
      <div class="card mb-4">
        <div class="card-header d-flex align-items-center">
          <div class="flex-grow-1">Users</div>
          <div class="flex-shrink-0">
            <button class="btn btn-success"><i class="fa fa-plus-circle mr-1"></i>Invite</button>
          </div>
        </div>
        <ul class="list-group list-group-flush">
          <Member v-for="member in members" :key="member.id" :member="member" />
        </ul>
      </div>
    </div>
    <div class="col-12 col-xxl-5">
      <div class="card">
        <div class="card-header">Pending Invitations</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  apollo: {
    members: {
      query: require("@/queries/Organization/Members.graphql"),
      update: data => data.organization?.members,
      variables() {
        return {
          oid: this.$route.params.oid,
        };
      },
    },
  },
  components: {
    Member: () => import("@/components/Organization/Member"),
  },
};
</script>
