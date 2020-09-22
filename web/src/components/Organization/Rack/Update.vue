<template>
  <tr :id="id">
    <td>
      <code>{{ update.version }}</code>
    </td>
    <td>
      <code>{{ update.parameters }}</code>
    </td>
    <td>
      <RelativeTime :time="update.created" />
    </td>
  </tr>
</template>

<script>
export default {
  computed: {
    id() {
      return `update-${this._uid}`;
    },
  },
  props: ["update"],
};
</script>
