<template>
  <div class="row">
    <div class="col-sm-10 offset-sm-1 col-md-8 offset-md-2 col-lg-6 offset-lg-3 col-xl-4 offset-xl-4">
      <div class="card">
        <form @submit="submit">
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
              <b-button type="submit" id="login" variant="primary">
                Login
              </b-button>
            </div>
          </div>
        </form>
      </div>
    </div>
    <b-modal id="token-authenticate" hide-header hide-footer>
      <div class="d-flex justify-content-center mt-3 mb-3">
        <h5 class="font-weight-bold mb-0">Please activate your security token now</h5>
      </div>
    </b-modal>
  </div>
</template>

<script>
import Authentication from "@/mixins/Authentication";
import Error from "@/mixins/Error";
import u2f from "@/scripts/u2f";

export default {
  created() {
    const tag = document.createElement("script");
    tag.setAttribute("src", "/scripts/u2f.js");
  },
  data() {
    return {
      alert: "",
      email: "",
      password: "",
    };
  },
  methods: {
    submit(e) {
      e.preventDefault();
      this.alert = "";
      this.$apollo
        .mutate({
          mutation: require("@/queries/Login.graphql"),
          variables: {
            email: this.$data.email,
            password: this.$data.password,
          },
        })
        .then(async result => {
          await this.login(result.data.login.key);
          if (result.data.login.session === null) {
            this.token_authentication_request();
          } else {
            this.$router.push({ name: "home" });
          }
        })
        .catch(err => {
          this.alert = this.graphQLErrors(err);
        });
    },
    token_authentication_request() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Token/AuthenticationRequest.graphql"),
        })
        .then(result => {
          const req = result.data.token_authentication_request;
          const data = JSON.parse(req.data);
          this.$bvModal.show("token-authenticate");
          u2f.sign(data.appId, data.challenge, data.registeredKeys, this.token_authentication_response(req.id), 30);
        })
        .catch(() => {
          this.alert = "invalid token";
        });
    },
    token_authentication_response(id) {
      const that = this;
      const apollo = this.$apollo;
      const login = this.login;
      const modal = this.$bvModal;
      const router = this.$router;
      return function(token) {
        if (token.errorCode > 0) {
          that.alert = "invalid token";
          return;
        }
        apollo
          .mutate({
            mutation: require("@/queries/Token/AuthenticationResponse.graphql"),
            variables: {
              id: id,
              data: JSON.stringify(token),
            },
          })
          .then(async result => {
            modal.hide("token-authenticate");
            await login(result.data.token_authentication_response.key);
            router.push({ name: "home" });
          })
          .catch(() => {
            that.alert = "invalid token";
          });
      };
    },
  },
  mixins: [Authentication, Error],
};
</script>
