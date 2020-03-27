<template>
  <div id="app">
    <div class="d-flex" style="min-height:100%;">
      <Sidebar></Sidebar>

      <main role="main" class="p-4 pt-4 flex-grow-1">
        <!-- {{ range flashes }}
				<div class="alert alert-{{.Kind}} alert-dismissible fade show" role="alert">
					<strong>{{ .Message }}</strong>
					<button type="button" class="close" data-dismiss="alert" aria-label="Close">
						<span aria-hidden="true">&times;</span>
					</button>
				</div>
			{{ end }} -->

        <router-view />
      </main>

      <div class="infobar d-none flex-shrink-0" style="width: 300px;"></div>
    </div>
  </div>

  <!-- <div class="modal fade" id="organization-create" tabindex="-1" role="dialog">
				<div class="modal-dialog" role="document">
					<div class="modal-content">
						<div class="modal-body">
							<div class="alert alert-danger fade show" role="alert" id="alert" style="display:none">
								<strong id="message"></strong>
							</div>
							<form onsubmit="javascript:organization_create();return false">
								<div class="form-group">
									<label for="name">Organization Name</label>
									<input type="text" id="organizationName" name="name" class="form-control" required minlength="2" pattern="[a-z0-9-]+" placeholder="company-name">
								</div>
								<button type="submit" class="btn btn-success">
									<i class="fa fa-plus-circle"></i>
									Create Organization
								</button>
							</form>
						</div>
					</div>
				</div>
			</div> -->

  <!-- <div id="nav">
			<router-link to="/">Home</router-link> |
			<router-link to="/about">About</router-link> |
			<router-link to="/test">Test</router-link>
		</div>
		<router-view /> -->
  <!-- </div> -->
</template>

<style lang="scss">
// #app {
//   font-family: Avenir, Helvetica, Arial, sans-serif;
//   -webkit-font-smoothing: antialiased;
//   -moz-osx-font-smoothing: grayscale;
//   text-align: center;
//   color: #2c3e50;
// }

@import "~bootstrap";
@import "~bootstrap-vue";

@import "@/styles/body.scss";
@import "@/styles/nav.scss";
@import "@/styles/sidebar.scss";
</style>

<script>
import Sidebar from "@/components/Sidebar.vue";

export default {
  components: {
    Sidebar
  }
};
</script>
