<template>
  <b-modal id="account" title="Account Settings" hide-header hide-footer body-bg-variant="dark">
    <div class="card mb-3">
      <div class="card-header">Account Information</div>
      <ul class="list-group list-group-flush">
        <div
          v-if="email_alert"
          class="alert alert-danger rounded-0 mb-0"
          role="alert"
        >{{ email_alert }}</div>
        <li class="list-group-item d-flex align-items-center">
          <div class="flex-shrink-0">
            <label class="mb-0">Email</label>
          </div>
          <div class="flex-grow-1 ml-3 mr-3">
            <input v-if="email_editing" type="email" v-model="email" class="form-control" />
            <code v-else class="text-secondary">{{ email }}</code>
          </div>
          <div class="flex-shrink-0">
            <div v-if="email_editing">
              <b-button variant="danger" @click="email_edit(false)" class="mr-2">
                <i class="fa fa-remove fa-fw" />
              </b-button>
              <b-button variant="success" @click="email_save()">
                <i class="fas fa-check fa-fw" />
              </b-button>
            </div>
            <b-button v-else variant="primary" @click="email_edit(true)">
              <i class="fas fa-edit fa-fw" />
            </b-button>
          </div>
        </li>
        <div
          v-if="password_alert"
          class="alert alert-danger rounded-0 mb-0"
          role="alert"
        >{{ password_alert }}</div>
        <li class="list-group-item">
          <div v-if="password_editing">
            <div class="row d-flex align-items-center mb-3">
              <div class="col-4">
                <label class="mb-0">Old Password</label>
              </div>
              <div class="col-8">
                <input ref="password_old" class="form-control" type="password" />
              </div>
            </div>
            <div class="row d-flex align-items-center mb-3">
              <div class="col-4">
                <label class="mb-0">New Password</label>
              </div>
              <div class="col-8">
                <input ref="password_new" class="form-control" type="password" />
              </div>
            </div>
            <div class="d-flex justify-content-end">
              <b-button variant="danger" @click="password_edit(false)">
                <i class="fa fa-ban mr-1" /> Cancel
              </b-button>
              <b-button variant="success" class="ml-2" @click="password_save()">
                <i class="fa fa-check mr-1" /> Save
              </b-button>
            </div>
          </div>
          <div v-else class="d-flex align-items-center">
            <div class="flex-shrink-0">
              <label class="mb-0">Password</label>
            </div>
            <div class="flex-grow-1 ml-3 mr-3">
              <code class="text-secondary">**********</code>
            </div>
            <div class="flex-shrink-0">
              <b-button variant="primary" @click="password_edit(true)">
                <i class="fas fa-edit fa-fw" />
              </b-button>
            </div>
          </div>
        </li>
        <li class="list-group-item d-flex align-items-center">
          <div class="flex-shrink-0">
            <label class="mb-0">CLI Token</label>
          </div>
          <div class="flex-grow-1 ml-3 mr-3">
            <code class="text-secondary">**********</code>
          </div>
          <div class="flex-shrink-0">
            <b-button variant="danger" @click="cli_token_reset()">
              <i class="fas fa-sync fa-fw" />
            </b-button>
          </div>
        </li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header d-flex align-items-center">
        <div class="flex-grow-1">Security Tokens</div>
        <div class="flex-shrink-0" style="font-size: 0.8em;">Last Used</div>
        <div class="flex-shrink-0 ml-3">
          <b-button variant="success" @click="token_register_request()" class="pl-2 pr-2">
            <i class="fas fa-plus-circle fa-fw" />
          </b-button>
        </div>
      </div>
      <ul class="list-group list-group-flush">
        <li
          v-for="token in tokens"
          :key="token.id"
          class="list-group-item d-flex align-items-center"
        >
          <div class="flex-grow-1">
            <code>{{ token.name }}</code>
          </div>
          <div class="flex-shrink-0 text-secondary">
            <Timeago v-if="token.used > 0" :datetime="token.used * 1000" />
            <div v-else>never</div>
          </div>
          <div class="flex-shrink-0 ml-3">
            <b-button variant="danger" class="pl-2 pr-2" @click="token_delete(token.id)">
              <i class="fa fa-times fa-fw" />
            </b-button>
          </div>
        </li>
      </ul>
    </div>
    <b-modal id="cli-token" hide-header hide-footer size="lg">
      <div class="mb-3">
        <h6 class="font-weight-bold mb-0">Run the following command to login with your new CLI Key</h6>
      </div>
      <div class="input-group">
        <input
          class="form-control bg-dark text-light text-monospace"
          v-model="login_command"
          ref="login_command"
        />
        <div class="input-group-append">
          <button
            class="btn btn-light clippy input-group-text btn-sm"
            title="Copied"
            ref="copy_login_command"
            @click="copy_login_command();"
          >
            <img class="clippy" src="/images/clippy.svg" width="14" />
          </button>
        </div>
      </div>
    </b-modal>
    <b-modal id="token-register" hide-header hide-footer>
      <div class="d-flex justify-content-center mt-3 mb-3">
        <h5 class="font-weight-bold mb-0">Please activate your security token now</h5>
      </div>
      <div
        v-if="token_alert"
        class="alert alert-danger d-flex justify-content-center"
        role="alert"
      >{{ token_alert }}</div>
    </b-modal>
  </b-modal>
</template>

<script>
import $ from "jquery";
import Error from "@/mixins/Error";
import u2f from "@/scripts/u2f";

export default {
  apollo: {
    tokens: {
      query: require("@/queries/Tokens.graphql"),
      update: data => data.tokens,
    },
    user: {
      query: require("@/queries/User.graphql"),
      result(res) {
        const user = res.data.user;
        this.email = user.email;
      },
    },
  },
  data() {
    return {
      email: "",
      email_alert: "",
      email_editing: false,
      login_command: "",
      password_alert: "",
      password_editing: false,
      token_alert: "",
    };
  },
  methods: {
    cli_token_reset() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/User/CliTokenReset.graphql"),
        })
        .then(res => {
          const key = res.data.user_cli_token_reset;
          this.login_command = `convox login ${window.location.hostname} -t ${key}`;
          this.$bvModal.show("cli-token");
        });
    },
    copy_login_command(btn) {
      this.$refs.login_command.select();
      document.execCommand("copy");
      $(this.$refs.copy_login_command).tooltip("show");
    },
    email_edit(editing) {
      this.email_alert = "";
      this.email_editing = editing;
    },
    email_save() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/User/Update.graphql"),
          variables: {
            email: this.email,
          },
        })
        .then(() => {
          this.$apollo.queries.user.refetch();
          this.email_alert = "";
          this.email_editing = false;
        })
        .catch(err => {
          this.email_alert = this.graphQLErrors(err);
        });
    },
    password_edit(editing) {
      this.password_alert = "";
      this.password_editing = editing;
    },
    password_save() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/User/PasswordUpdate.graphql"),
          variables: {
            old: this.$refs.password_old.value,
            new: this.$refs.password_new.value,
          },
        })
        .then(() => {
          this.$apollo.queries.user.refetch();
          this.password_alert = "";
          this.password_editing = false;
        })
        .catch(err => {
          this.password_alert = this.graphQLErrors(err);
        });
    },
    token_delete(id) {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Token/Delete.graphql"),
          variables: {
            id: id,
          },
        })
        .then(() => {
          this.$apollo.queries.tokens.refetch();
        });
    },
    token_register_request() {
      this.token_alert = "";
      this.$apollo
        .mutate({
          mutation: require("@/queries/Token/RegisterRequest.graphql"),
        })
        .then(result => {
          const req = result.data.token_register_request;
          const data = JSON.parse(req.data);
          this.$bvModal.show("token-register");
          u2f.register(data.appId, data.registerRequests, data.registeredKeys || [], this.token_register_response(req.id), 30);
        });
    },
    token_register_response(id) {
      const apollo = this.$apollo;
      const modal = this.$bvModal;
      const that = this;
      return function(token) {
        switch (token.errorCode) {
          case 1:
          case 2:
          case 3:
            that.token_alert = "invalid token";
            break;
          case 4:
            that.token_alert = "token already registered";
            break;
          case 5:
            that.token_alert = "timed out waiting for token";
            break;
          default:
            apollo
              .mutate({
                mutation: require("@/queries/Token/RegisterResponse.graphql"),
                variables: {
                  id: id,
                  data: JSON.stringify(token),
                },
              })
              .then(() => {
                modal.hide("token-register");
                that.$apollo.queries.tokens.refetch();
              })
              .catch(() => {
                that.token_alert = "invalid token";
              });
        }
      };
    },
  },
  mixins: [Error],
};
</script>
