<template>
  <div class="row">
    <div class="col-12 col-xl-6 col-xxl-4 mb-4">
      <div class="card">
        <div class="card-header">User Information</div>
        <div class="card-body">
          <div class="d-flex align-items-center">
            <div class="flex-shrink-0">Email</div>
            <div class="flex-grow-1 ml-3 mr-3"><input type="email" v-model="email" class="form-control" /></div>
            <div class="flex-shrink-0">
              <b-button variant="primary">
                <i class="fa fa-save mr-1" />
                Save
              </b-button>
            </div>
          </div>
        </div>
        <div class="card-footer d-flex">
          <div class="flex-grow-1">
            <b-button variant="danger">
              <i class="fa fa-sync mr-1" />
              Reset CLI Key
            </b-button>
          </div>
          <div class="flex-shrink-0">
            <b-button variant="primary">
              <i class="fa fa-edit mr-1" />
              Change Password
            </b-button>
          </div>
        </div>
      </div>
    </div>
    <div class="col-12 col-xl-6 col-xxl-4">
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
            <div class="flex-shrink-0"><Timeago :datetime="token.used * 1000" /></div>
            <div class="flex-shrink-0 ml-3">
              <b-button variant="danger" class="pl-2 pr-2">
                <i class="fa fa-times fa-fw" />
              </b-button>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <b-modal id="token-register" hide-header hide-footer>
      <div class="d-flex justify-content-center mt-3 mb-3">
        <h5 class="font-weight-bold mb-0">Please activate your security token now</h5>
      </div>
    </b-modal>
  </div>
</template>

<script>
import u2f from "@/scripts/u2f";

export default {
  apollo: {
    tokens: {
      query: require("@/queries/Tokens.graphql"),
      update: data => data.tokens,
    },
  },
  data() {
    return {
      email: "",
    };
  },
  methods: {
    token_register_request() {
      this.$apollo
        .mutate({
          mutation: require("@/queries/Token/RegisterRequest.graphql"),
        })
        .then(result => {
          const req = result.data.token_register_request;
          const data = JSON.parse(req.data);
          this.$bvModal.show("token-register");
          u2f.register(data.appId, data.registerRequests, data.registeredKeys || [], this.token_register_response(req.id), 30);
        })
        .catch(err => {
          // TODO handle this error
          alert(err);
        });
    },
    token_register_response(id) {
      const apollo = this.$apollo;
      const modal = this.$bvModal;
      return function(token) {
        if (token.errorCode > 0) {
          // handle error
          return;
        }
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
          })
          .catch(err => {
            // TODO handle this error
            alert(err);
          });
      };
    },
  },
};
</script>
