<template>
  <div class="row">
    <div class="col-sm-10 offset-sm-1 col-md-8 offset-md-2 col-lg-6 offset-lg-3 col-xl-4 offset-xl-4">
      <div class="card">
        <div class="card-header">Login</div>
        <div class="position-relative d-flex">
          <div class="card-body">
            <div class="form-group">
              <label for="login-email">Email</label>
              <input id="login-email" class="form-control" type="email" v-model="email" />
            </div>
            <div class="form-group">
              <label for="login-password">Password</label>
              <input id="login-password" class="form-control" type="password" v-model="password" />
            </div>
          </div>
          <div
            id="login-alert"
            class="alert alert-danger rounded-0 h-100 w-100 position-absolute d-none align-items-center justify-content-center"
          ></div>
        </div>
        <div class="card-footer d-flex align-items-center">
          <div class="flex-fill">
            <a href>Forgot Password</a>
          </div>
          <div class="flex-fill text-right">
            <button id="login" class="btn btn-primary" @click="submit($event)">Login</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import $ from "jquery";

import Authentication from "@/mixins/Authentication";

export default {
  data() {
    return {
      email: "",
      password: "",
    };
  },
  methods: {
    submit(event) {
      event.target.disabled = true;
      this.$apollo
        .mutate({
          mutation: require("@/queries/Login.graphql"),
          variables: {
            email: this.$data.email,
            password: this.$data.password,
          },
        })
        .then((result) => {
          this.login(result.data.login.token);
        })
        .catch((error) => {
          var message = error.graphQLErrors.map((err) => err.message).join(", ");
          if (message != "") {
            var alert = $("#login-alert");
            alert.html(message);
            alert.removeClass("d-none");
            alert.addClass("d-flex");
            alert.css("opacity", "0%");
            alert.animate({ opacity: "100%" }, () => {
              window.setTimeout(() => {
                alert.animate({ opacity: "0%" }, () => {
                  alert.removeClass("d-flex");
                  alert.addClass("d-none");
                });
              }, 800);
            });
          }
        })
        .finally(() => {
          this.$router.push({
            name: "home",
          });
        });
    },
  },
  mixins: [Authentication],
};
</script>
