<template>
  <div class="public">
    <Navbar></Navbar>
    <div class="d-flex" style="min-height:100%;">
      <main role="main" class="p-4 pt-4 flex-grow-1">
        <!-- {{ range flashes }}
				<div class="alert alert-{{.Kind}} alert-dismissible fade show" role="alert">
					<strong>{{ .Message }}</strong>
					<button type="button" class="close" data-dismiss="alert" aria-label="Close">
						<span aria-hidden="true">&times;</span>
					</button>
				</div>
        {{ end }}-->

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
import Navbar from "@/components/Navbar.vue";

export default {
  components: {
    Navbar,
  },
};
</script>
