<template>
  <div class="row">
    <div class="col-sm-10 offset-sm-1 col-md-8 offset-md-2 col-lg-6 offset-lg-3 col-xl-4 offset-xl-4">
      <div class="card">
        <div class="card-header">Login</div>
        <div class="position-relative d-flex">
          <div class="card-body">
            <div v-if="alert" class="alert alert-danger" role="alert">{{ alert }}</div>
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
            <b-button id="login" variant="primary" @click="submit()">
              Login
            </b-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Authentication from "@/mixins/Authentication";
import Error from "@/mixins/Error";
import u2f from "@/scripts/u2f";

export default {
  data() {
    return {
      alert: "",
      email: "",
      password: "",
    };
  },
  methods: {
    submit() {
      this.alert = "";
      this.$apollo
        .mutate({
          mutation: require("@/queries/Login.graphql"),
          variables: {
            email: this.$data.email,
            password: this.$data.password,
          },
        })
        .then(result => {
          console.log("result", result);
          if (result.data.login.session === null) {
            this.token_challenge();
          } else {
            this.login(result.data.login.key);
            this.$router.push({
              name: "home",
            });
          }
        })
        .catch(err => {
          this.alert = this.graphQLErrors(err);
        });
    },
    token_challenge() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Token/AuthenticationRequest.graphql"),
        })
        .then(async result => {
          console.log("result", result);
        })
        .catch(err => {
          console.log("err", err);
        });
    },
  },
  mixins: [Authentication, Error],
};
</script>
