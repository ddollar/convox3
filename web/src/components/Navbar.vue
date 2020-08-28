<template>
  <header>
    <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
      <a class="navbar-brand" href="/">
        <img src="/images/logo.png" style="height:32px; margin-top:-6px; margin-bottom:-4px; margin-left: -2px;" />
      </a>
      <button
        class="navbar-toggler"
        type="button"
        data-toggle="collapse"
        data-target="#navbar-links"
        aria-controls="navbar-links"
        aria-expanded="false"
        aria-label="Toggle Navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbar-links">
        <ul class="navbar-nav ml-auto">
          <!-- {{ if current_organization }}
					{{ $cl := checklist }}
					{{ if checklist_count $cl }}
						<li class="nav-item dropdown pr-4">
							<button type="button" class="btn btn-danger font-weight-bold" data-toggle="dropdown" style="border-radius:16px; padding:2px 11px; padding-left:10px; margin-top:5px; font-weight:900;">
								{{ checklist_count $cl }}
							</button>
							<div class="dropdown-menu dropdown-menu-right bg-dark">
								{{ template "checklist" $cl }}
							</div>
						</li>
					{{ end }} -->
          <li class="nav-item pr-4">
            <slot></slot>
          </li>
          <li class="nav-item pr-4">
            <a class="nav-link" href="https://docs.convox.com" target="_blank">
              <i class="fas fa-book-open mr-1"></i>
              Documentation
            </a>
          </li>
          <li class="nav-item pr-4" v-if="authenticated">
            <router-link to="/account" class="nav-link">
              <i class="fa fa-user mr-1" aria-hidden="true"></i>
              Account
            </router-link>
          </li>
          <li class="nav-item pr-4" v-if="authenticated">
            <a class="nav-link" @click="logout()" href="/login">
              <i class="fa fa-sign-out mr-1" aria-hidden="true"></i>
              Logout
            </a>
          </li>
          <!-- {{ if authenticated }}
					{{ if forums }}
						<li class="nav-item pr-4">
							<a class="nav-link" href="https://community.convox.com">
								<i class="fas fa-comments mr-1" aria-hidden="true"></i>
								Forums
							</a>
						</li>
					{{ end }}
				{{ end }} -->
        </ul>
      </div>
    </nav>
  </header>
</template>

<style lang="scss">
nav.bg-dark {
  background-color: #101921 !important;
}
</style>

<script>
import Authentication from "@/mixins/Authentication";

export default {
  mixins: [Authentication],
};
</script>
