<template>
  <nav id="sidebar" class="d-none d-md-block flex-shrink-0">
    <ul class="nav flex-column">
      <Nav icon="fa-server" to="racks">Racks</Nav>
      <Nav icon="fa-check-square-o" to="workflows">Workflows</Nav>
      <Nav icon="fa-wrench" to="jobs">Jobs</Nav>
      <Nav icon="fa-user" to="users">Users</Nav>
      <Nav icon="fa-search" to="audits">Audit Log</Nav>
      <Nav icon="fa-plug" to="integrations">Integrations</Nav>
      <Nav icon="fa-credit-card" to="billing" v-if="billing">Billing</Nav>
      <Nav icon="fa-cog" to="settings">Settings</Nav>
      <Nav icon="fa-question-circle" to="support">Support</Nav>
    </ul>
  </nav>
</template>

<style lang="scss" src="@/styles/sidebar.scss"></style>

<script>
export default {
  components: {
    Nav: () => import("@/components/Sidebar/Nav"),
  },
  computed: {
    billing() {
      return true;
    },
  },
};
</script>
