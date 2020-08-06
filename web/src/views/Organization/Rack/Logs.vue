<template>
  <div class="logs" id="rlt">{{ logs }}</div>
</template>

<script>
export default {
  apollo: {
    $subscribe: {
      logs: {
        query: require("@/queries/Rack/Logs.graphql"),
        variables() {
          return {
            oid: this.$route.params.oid,
            rid: this.$route.params.rid,
          };
        },
        result({ data }) {
          this.logs += `${data.rack_logs.line}\n`;
        },
      },
    },
  },
  data: function() {
    return {
      logs:
        "1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n1\n",
    };
  },
  methods: {
    scrollMonitor: function() {
      const el = this.$el;
      el.dataset.bottom = el.scrollTop === el.scrollHeight - el.offsetHeight;
    },
  },
  mounted() {
    this.$el.addEventListener("scroll", this.scrollMonitor);
    this.scrollMonitor();
  },
  updated() {
    const el = this.$el;

    if (el.dataset.bottom === "true") {
      el.scrollTop = el.scrollHeight;
    }
  },
};
</script>
