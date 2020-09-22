<template>
  <tr :id="id">
    <td>
      <code>{{ build.id }}</code>
    </td>
    <td>{{ build.description }}</td>
    <td>
      <RelativeTime :time="build.started" />
    </td>
    <td>{{ build.ended }}</td>
    <td>{{ build.status }}</td>
    <td>
      <button class="btn btn-sm btn-secondary">MM</button>
    </td>
  </tr>
</template>

<script>
export default {
  components: {
    RelativeTime: () => import("@/components/RelativeTime.vue"),
  },
  computed: {
    datetime() {
      return new Date(this.build.started * 1000);
    },
    id() {
      return `build-${this._uid}`;
    },
  },
  methods: {},
  props: ["build"],
};
</script>
