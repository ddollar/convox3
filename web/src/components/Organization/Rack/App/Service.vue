<template>
  <tr :id="id">
    <td>
      <code>{{ service.name }}</code>
    </td>
    <td class="text-right">{{ service.count }}</td>
    <td class="text-right">{{ service.cpu }}</td>
    <td class="text-right">{{ service.mem }}</td>
  </tr>
</template>

<script>
export default {
  computed: {
    id() {
      return `service-${this._uid}`;
    },
    percent(value) {
      return `${(value * 100).toFixed(1)}%`;
    },
  },
  methods: {},
  props: ["service"],
};
</script>
