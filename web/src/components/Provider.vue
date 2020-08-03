<template>
  <span>
    <i :class="icon_class(kind)"></i>
    <span v-if="text" class="ml-2">{{ text }}</span>
  </span>
</template>

<script>
export default {
  methods: {
    brand_class(kind) {
      switch (kind) {
        case "aws":
          return "fab fa-aws";
        case "azure":
          return "fab fa-windows";
        case "do":
          return "fab fa-digital-ocean";
        case "gcp":
          return "fab fa-google";
        case "github":
          return "fab fa-github";
        case "github-enterprise":
          return "fab fa-github-square";
        case "gitlab":
          return "fab fa-gitlab";
        case "slack":
          return "fab fa-slack";
        default:
          return "fas fa-question-circle";
      }
    },
    icon_class(kind) {
      return `fa ${this.brand_class(kind)} fa-fw`;
    },
  },
  props: ["kind", "text"],
};
</script>
