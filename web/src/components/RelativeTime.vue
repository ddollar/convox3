<template>
  <div>
    <b-tooltip :target="id" placement="left">{{ datetime }}</b-tooltip>
    <Timeago v-if="time > 0" :id="id" :datetime="datetime" />
    <div v-else>{{ zero }}</div>
  </div>
</template>

<script>
export default {
  computed: {
    datetime() {
      return new Date(this.time * 1000);
    },
    id() {
      return `timeago-${this._uid}`;
    },
  },
  props: ["time", "zero"],
};
</script>