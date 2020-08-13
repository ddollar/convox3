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
      return `app-status-${this._uid}`;
    },
    status_color() {
      if (!this.color) return "";

      switch (this.status) {
        case "creating":
          return "text-success";
        case "deleting":
          return "text-danger";
        case "failed":
          return "text-warning";
        case "rollback":
          return "text-danger";
        case "running":
          return "text-success";
        case "updating":
          return "text-secondary";
        default:
          return "text-danger";
      }
    },
    status_icon() {
      switch (this.status) {
        case "creating":
          return "fas fa-sync-alt fa-spin";
        case "deleting":
          return "fas fa-sync-alt fa-spin";
        case "failed":
          return "fas fa-check-square";
        case "rollback":
          return "fas fa-sync-alt fa-spin";
        case "running":
          return "fas fa-check-square";
        case "updating":
          return "fas fa-sync-alt fa-spin";
        default:
          return "fas fa-question-circle";
      }
    },
  },
  props: {
    color: {
      type: Boolean,
      default: false,
    },
    status: {
      type: String,
    },
  },
};
</script>
