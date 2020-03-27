<template>
  <div class="organization">
    <Navbar>
      <div class="dropdown" id="organization-chooser">
        <button
          type="button"
          class="btn btn-dark dropdown-toggle"
          data-toggle="dropdown"
        >
          <span class="pl-2 pr-3">{{ organization.name }}</span>
        </button>
        <div class="dropdown-menu dropdown-menu-right bg-dark">
          <router-link
            class="dropdown-item text-light bg-dark"
            v-for="organization in organizations"
            :to="{
              name: 'organization/racks',
              params: { oid: organization.id }
            }"
            :key="organization.id"
          >
            {{ organization.name }}
          </router-link>
          <div class="dropdown-divider"></div>
          <a
            class="dropdown-item text-light bg-dark"
            data-toggle="modal"
            data-target="#organization-create"
          >
            <i class="fa fa-plus-circle pr-1" aria-hidden="true"></i>
            Create Organization
          </a>
        </div>
      </div>
    </Navbar>
    <router-view />
  </div>
</template>

<script>
import Navbar from "@/components/Navbar.vue";

export default {
  apollo: {
    organizations: require("../queries/Organizations.graphql"),
    organization: {
      query: require("../queries/Organization.graphql"),
      variables() {
        return {
          id: this.$route.params.oid
        };
      }
    }
  },
  components: {
    Navbar
  },
  data() {
    return {
      organization: {}
    };
  }
};
</script>
