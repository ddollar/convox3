<template>
  <b-modal id="account" title="Account Settings" hide-header hide-footer body-bg-variant="dark">
    <div class="card mb-3">
      <div class="card-header">Account Information</div>
      <ul class="list-group list-group-flush">
        <div v-if="alert" class="alert alert-danger rounded-0 mb-0" role="alert">{{ alert }}</div>
        <li class="list-group-item d-flex align-items-center">
          <div class="flex-shrink-0"><label class="mb-0">Email</label></div>
          <div class="flex-grow-1 ml-3 mr-3">
            <input v-if="editing" type="email" v-model="email" class="form-control" />
            <code v-else class="text-secondary">{{ email }}</code>
          </div>
          <div class="flex-shrink-0">
            <div v-if="editing">
              <b-button variant="danger" @click="edit(false)" class="mr-2"><i class="fa fa-remove fa-fw"/></b-button>
              <b-button variant="success" @click="save()"><i class="fas fa-check fa-fw"/></b-button>
            </div>
            <b-button v-else variant="primary" @click="edit(true)"><i class="fas fa-edit fa-fw"/></b-button>
          </div>
        </li>
        <li class="list-group-item d-flex align-items-center">
          <div class="flex-shrink-0"><label class="mb-0">Password</label></div>
          <div class="flex-grow-1 ml-3 mr-3"><code class="text-secondary">**********</code></div>
          <div class="flex-shrink-0">
            <b-button variant="primary">
              <i class="fas fa-edit fa-fw" />
            </b-button>
          </div>
        </li>
        <li class="list-group-item d-flex align-items-center">
          <div class="flex-shrink-0"><label class="mb-0">CLI Token</label></div>
          <div class="flex-grow-1 ml-3 mr-3"><code class="text-secondary">**********</code></div>
          <div class="flex-shrink-0">
            <b-button variant="danger">
              <i class="fas fa-sync fa-fw" />
            </b-button>
          </div>
        </li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header d-flex align-items-center">
        <div class="flex-grow-1">
          Security Tokens
        </div>
        <div class="flex-shrink-0" style="font-size: 0.8em;">Last Used</div>
        <div class="flex-shrink-0 ml-3">
          <b-button variant="success" @click="token_register_request()" class="pl-2 pr-2"><i class="fas fa-plus-circle fa-fw"/></b-button>
        </div>
      </div>
      <ul class="list-group list-group-flush">
        <li v-for="token in tokens" :key="token.id" class="list-group-item d-flex align-items-center">
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
    <b-modal id="token-register" hide-header hide-footer>
      <div class="d-flex justify-content-center mt-3 mb-3">
        <h5 class="font-weight-bold mb-0">Please activate your security token now</h5>
      </div>
      <div v-if="alert" class="alert alert-danger d-flex justify-content-center" role="alert">{{ alert }}</div>
    </b-modal>
  </b-modal>
</template>

<script>
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
      alert: "",
      editing: false,
      email: "",
    };
  },
  methods: {
    edit(editing) {
      this.alert = "";
      this.editing = editing;
    },
    save() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/User/Update.graphql"),
          variables: {
            email: this.email,
          },
        })
        .then(() => {
          this.$apollo.queries.user.refetch();
          this.alert = "";
          this.editing = false;
        })
        .catch(err => {
          this.alert = this.graphQLErrors(err);
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
      this.alert = "";
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
            that.alert = "invalid token";
            break;
          case 4:
            that.alert = "token already registered";
            break;
          case 5:
            that.alert = "timed out waiting for token";
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
                that.alert = "invalid token";
              });
        }
      };
    },
  },
  mixins: [Error],
  mounted() {
    this.alert = "";
  },
};
</script>
