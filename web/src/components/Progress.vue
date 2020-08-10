<template>
  <div>
    <div :id="id" class="progress bg-white" :title="title">
      <div :style="style"></div>
    </div>
    <b-tooltip :target="id" placement="bottom">{{ title }}</b-tooltip>
  </div>
</template>

<style lang="scss" scoped>
.progress {
  border: 1px #bbb solid;
}
</style>

<script>
export default {
  computed: {
    id() {
      return `progress-${this._uid}`;
    },
    percent() {
      return (parseFloat(this.current) / parseFloat(this.max)) * 100;
    },
    style() {
      return `background-color:${this.color}; width:${this.percent.toFixed(0)}%`;
    },
    title() {
      return `${this.percent.toFixed(1)}%`;
    },
  },
  props: {
    current: { type: Number, default: 0.0 },
    color: { type: String, default: "#0184bf" },
    max: { type: Number, default: 1.0 },
  },
};
</script>
