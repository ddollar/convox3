<template>
  <div id="app">
    <router-view></router-view>
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
  </div>-->

  <!-- <div id="nav">
			<router-link to="/">Home</router-link> |
			<router-link to="/about">About</router-link> |
			<router-link to="/test">Test</router-link>
		</div>
  <router-view />-->
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
@import "@/styles/bootstrap.scss";

@import "~bootstrap";
@import "~bootstrap-vue";

@import "@/styles/body.scss";
@import "@/styles/card.scss";
@import "@/styles/clickable.scss";
@import "@/styles/logs.scss";
@import "@/styles/nav.scss";

$fa-font-path: "~@fortawesome/fontawesome-free/webfonts";
@import "~@fortawesome/fontawesome-free/scss/fontawesome";
@import "~@fortawesome/fontawesome-free/scss/solid";
@import "~@fortawesome/fontawesome-free/scss/regular";
@import "~@fortawesome/fontawesome-free/scss/brands";
@import "~@fortawesome/fontawesome-free/scss/v4-shims";
</style>

<script>
// var scrollMonitorTop = 0;
// var scrollMonitorEnd = true;

// $.initialize(".scroll-monitor", function() {
//   var monitor = $(this);
//   if (scrollMonitorEnd) {
//     monitor.prop("scrollTop", monitor.prop("scrollHeight"));
//   } else {
//     monitor.prop("scrollTop", scrollMonitorTop);
//   }
//   monitor.on("scroll", function() {
//     scrollMonitorTop = $(this).prop("scrollTop");
//     scrollMonitorEnd = $(this).scrollTop() + $(this).innerHeight() >= $(this)[0].scrollHeight;
//   });
// });

export default {};
</script>
