<template>
  <div class="private">
    <Navbar>
      <Picker />
    </Navbar>
    <div class="d-flex" style="min-height:100%;">
      <Sidebar />

      <main role="main" class="p-4 pt-4 flex-grow-1">
        <router-view />
      </main>

      <div class="infobar d-none flex-shrink-0" style="width: 300px;"></div>
    </div>
  </div>
</template>

<style lang="scss">
main {
  height: 100vh;
  margin-top: 10px;
  overflow-y: auto;
  position: relative;
}
</style>

<script>
export default {
  components: {
    Navbar: () => import("@/components/Navbar.vue"),
    Picker: () => import("@/components/Organization/Picker.vue"),
    Sidebar: () => import("@/components/Sidebar.vue"),
  },
};
</script>
