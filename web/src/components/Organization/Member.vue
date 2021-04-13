<template>
  <li class="list-group-item d-flex align-items-center">
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
        <form method="post" action="/organizations/7e45b1e6-803e-4d74-8094-bc25a1ba2533/users/d53aebda-0b24-4376-a88f-acba3c43da4c">
          <input type="hidden" name="role" value="administrator" />
          <button type="submit" class="dropdown-item">
            Administrator
          </button>
        </form>
        <form method="post" action="/organizations/7e45b1e6-803e-4d74-8094-bc25a1ba2533/users/d53aebda-0b24-4376-a88f-acba3c43da4c">
          <input type="hidden" name="role" value="operator" />
          <button type="submit" class="dropdown-item">Operator</button>
        </form>
        <form method="post" action="/organizations/7e45b1e6-803e-4d74-8094-bc25a1ba2533/users/d53aebda-0b24-4376-a88f-acba3c43da4c">
          <input type="hidden" name="role" value="developer" />
          <button type="submit" class="dropdown-item">Developer</button>
        </form>
      </div>
    </div>
    <div class="flex-shrink-0 ml-2">
      <button class="btn btn-danger" @click="remove()">
        <i class="fa fa-times"></i>
      </button>
    </div>
  </li>
</template>

<script>
export default {
  methods: {
    remove() {
      this.$apollo.mutate({
        mutation: require("@/queries/Organization/Member/Delete.graphql"),
        variables: {
          oid: this.$route.params.oid,
          uid: this.member.user.id,
        },
      });
    },
  },
  props: ["member"],
};
</script>
