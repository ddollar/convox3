<template>
  <span>
    <i :id="id" :class="css" />
    <b-tooltip :target="id" placement="left">{{ status }}</b-tooltip>
  </span>
</template>

<script>
export default {
  computed: {
    css() {
      return `fa ${this.status_icon} ${this.status_color}`;
    },
    id() {
      return `rack-status-${this._uid}`;
    },
    status_color() {
      if (!this.color) return "";

      switch (this.status) {
        case "converging":
          return "text-secondary";
        case "deleting":
          return "text-danger";
        case "failed":
          return "text-danger";
        case "installing":
          return "text-success";
        case "rollback":
          return "text-danger";
        case "running":
          return "text-success";
        case "updating":
          return "text-secondary";
        case "uninstalling":
          return "text-danger";
        default:
          return "text-danger";
      }
    },
    status_icon() {
      switch (this.status) {
        case "converging":
          return "fas fa-sync-alt fa-spin";
        case "deleting":
          return "fas fa-sync-alt fa-spin";
        case "failed":
          return "fas fa-times-circle";
        case "installing":
          return "fas fa-sync-alt fa-spin";
        case "rollback":
          return "fas fa-sync-alt fa-spin";
        case "running":
          return "fas fa-check-square";
        case "updating":
          return "fas fa-sync-alt fa-spin";
        case "uninstalling":
          return "fas fa-sync-alt fa-spin";
        default:
          return "fas fa-question-circle";
      }
    }
  },
  props: {
    color: {
      type: Boolean,
      default: false
    },
    status: {
      type: String
    }
  }
};
</script>
