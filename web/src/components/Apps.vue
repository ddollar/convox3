<template>
  <div class="tab-pane active" id="rack-apps">
    <ul class="list-group list-group-flush">
      <li class="list-group-item d-flex align-items-center" style="font-weight:bold;">
        <div class="flex-grow-1">NAME</div>
        <div class="flex-shrink-0 ml-4" style="width:100px;">RELEASE</div>
        <div class="flex-shrink-0 ml-4">
          <button
            class="btn btn-success btn-sm"
            data-toggle="modal"
            data-target="#app-create-9bed02be-d1df-425e-b02d-a755f9188e45"
          >
            <i class="fa fa-plus-circle"></i>
          </button>
        </div>
      </li>
      <li class="list-group-item d-flex align-items-center" v-for="app in apps" :key="app.name">
        <div class="flex-grow-1 flex-shrink-0">
          <i class="fas fa-check-square text-success mr-2"></i>
          {{app.name}}
        </div>
        <div class="flex-shrink-0 ml-4" style="width:100px;">{{app.release}}</div>
        <div class="flex-shrink-0 ml-4">
          <button class="btn btn-danger btn-sm app-delete">
            <i
              class="fa fa-remove"
              style="font-size:1.3em; padding-left:1px; padding-right:1px; padding-top:2px; padding-bottom:2px;"
            ></i>
          </button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  apollo: {
    apps: {
      query: require("@/queries/Apps.graphql"),
      update: data => data.organization?.rack?.apps,
      variables() {
        return {
          oid: this.oid,
          rid: this.rid
        };
      }
    }
  },
  props: ["oid", "rid"]
};
</script>