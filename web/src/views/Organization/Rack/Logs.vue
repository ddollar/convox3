<template>
  <div class="logs" id="rlt">{{ logs }}</div>
</template>

<script>
export default {
  apollo: {
    $subscribe: {
      logs: {
        query: require("@/queries/Organization/Rack/Logs.graphql"),
        variables() {
          return {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid
          };
        },
        result({ data }) {
          this.logs += `${data.rack_logs.line}\n`;
        }
      }
    }
  },
  data: function() {
    return {
      logs: ""
    };
  },
  methods: {
    scrollMonitor: function() {
      const el = this.$el;
      el.dataset.bottom = el.scrollTop === el.scrollHeight - el.offsetHeight;
    }
  },
  mounted() {
    this.$el.addEventListener("scroll", this.scrollMonitor);
    this.$el.dataset.bottom = true;
  },
  updated() {
    const el = this.$el;

    if (el.dataset.bottom === "true") {
      el.scrollTop = el.scrollHeight;
    }
  }
};
</script>
