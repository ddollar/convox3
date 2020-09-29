<template>
  <div class="async">
    <i v-if="loading" class="spinner"></i>
    <div v-else>{{ value }}</div>
  </div>
</template>

<script>
export default {
  props: ["loading", "value"],
};
</script>