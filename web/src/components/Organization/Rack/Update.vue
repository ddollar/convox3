<template>
  <tr :id="id">
    <td>
      <code>{{ update.version }}</code>
    </td>
    <td>
      <code>{{ update.parameters }}</code>
    </td>
    <td><Timeago v-if="update.created > 0" :datetime="datetime" /></td>
  </tr>
</template>

<script>
export default {
  computed: {
    datetime() {
      return new Date(this.process.started * 1000);
    },
    id() {
      return `update-${this._uid}`;
    },
  },
  props: ["update"],
};
</script>
