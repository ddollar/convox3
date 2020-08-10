<template>
  <div class="card mb-4">
    <div class="card-header">{{ title }}</div>
    <ul class="list-group list-group-flush">
      <li v-for="integration in integrations" :key="integration.id" class="list-group-item">
        <Provider :kind="integration.provider" :text="integration.title" />
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  components: {
    Provider: () => import("@/components/Provider")
  },
  props: ["integrations", "title"]
};
</script>