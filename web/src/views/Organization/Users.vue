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
          <li v-for="member in members" :key="member.user.id" class="list-group-item d-flex align-items-center">
            <div class="flex-grow-1">{{ member.user.email }}</div>
            <div class="flex-shrink-0">
              <button
                class="btn btn-secondary dropdown-toggle"
                type="button"
                id="dropdownMenuButton"
                data-toggle="dropdown"
                aria-haspopup="true"
                aria-expanded="false"
              >
                Administrator
              </button>
              <div class="dropdown-menu">
                <form
                  method="post"
                  action="/organizations/7e45b1e6-803e-4d74-8094-bc25a1ba2533/users/d53aebda-0b24-4376-a88f-acba3c43da4c"
                >
                  <input type="hidden" name="role" value="administrator" />
                  <button type="submit" class="dropdown-item">Administrator</button>
                </form>
                <form
                  method="post"
                  action="/organizations/7e45b1e6-803e-4d74-8094-bc25a1ba2533/users/d53aebda-0b24-4376-a88f-acba3c43da4c"
                >
                  <input type="hidden" name="role" value="operator" />
                  <button type="submit" class="dropdown-item">Operator</button>
                </form>
                <form
                  method="post"
                  action="/organizations/7e45b1e6-803e-4d74-8094-bc25a1ba2533/users/d53aebda-0b24-4376-a88f-acba3c43da4c"
                >
                  <input type="hidden" name="role" value="developer" />
                  <button type="submit" class="dropdown-item">Developer</button>
                </form>
              </div>
            </div>
            <div class="flex-shrink-0 ml-2">
              <button class="btn btn-danger">
                <i class="fa fa-times"></i>
              </button>
            </div>
          </li>
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
import Organization from "@/mixins/Organization";

export default {
  apollo: {
    members: {
      query: require("@/queries/Members.graphql"),
      skip() {
        return !this.organization.id;
      },
      update: (data) => data.organization?.members,
      variables() {
        return {
          oid: this.organization.id,
        };
      },
    },
  },
  mixins: [Organization],
};
</script>
