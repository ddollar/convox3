<template>
  <tr :id="id">
    <td>
      <code>{{ release.id }}</code>
    </td>
    <td>{{ release.description }}</td>
    <td>
      <code>{{ release.build }}</code>
    </td>
    <td>
      <RelativeTime :time="release.created" />
    </td>
    <td>
      <Manifest :manifest="release.manifest" :title="release.id" />
    </td>
  </tr>
</template>

<script>
export default {
  components: {
    Manifest: () => import("@/components/Manifest.vue"),
    RelativeTime: () => import("@/components/RelativeTime.vue"),
  },
  computed: {
    datetime() {
      return new Date(this.release.started * 1000);
    },
    id() {
      return `release-${this._uid}`;
    },
  },
  methods: {},
  props: ["release"],
};
</script>
