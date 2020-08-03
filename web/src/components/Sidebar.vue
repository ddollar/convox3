<template>
  <nav id="sidebar" class="d-none d-md-block flex-shrink-0">
    <ul class="nav flex-column">
      <SidebarLink icon="fa-server" to="racks">Racks</SidebarLink>
      <SidebarLink icon="fa-check-square-o" to="workflows">Workflows</SidebarLink>
      <SidebarLink icon="fa-wrench" to="jobs">Jobs</SidebarLink>
      <SidebarLink icon="fa-user" to="users">Users</SidebarLink>
      <SidebarLink icon="fa-search" to="audits">Audit Log</SidebarLink>
      <SidebarLink icon="fa-plug" to="integrations">Integrations</SidebarLink>
      <SidebarLink icon="fa-credit-card" to="billing" v-if="billing">Billing</SidebarLink>
      <SidebarLink icon="fa-cog" to="settings">Settings</SidebarLink>
      <SidebarLink icon="fa-question-circle" to="support">Support</SidebarLink>
    </ul>
  </nav>
</template>

<style lang="scss" src="@/styles/sidebar.scss"></style>

<script>
import SidebarLink from "@/components/SidebarLink.vue";

export default {
  components: {
    SidebarLink
  },
  computed: {
    billing() {
      return true;
    }
  }
};
</script>
