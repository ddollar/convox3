<template>
  <div class="uninstall logs" id="rlt">{{ logs }}</div>
</template>

<script>
export default {
  apollo: {
    $subscribe: {
      logs: {
        query: require("@/queries/Organization/Rack/Uninstall/Logs.graphql"),
        skip() {
          return this.rack.install == null;
        },
        variables() {
          return {
            oid: this.$route.params.oid,
            uid: this.rack.uninstall.id
          };
        },
        result({ data }) {
          this.logs += `${data.uninstall_logs.line}\n`;
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
      el.dataset.bottom = el.scrollTop >= el.scrollHeight - el.offsetHeight;
    }
  },
  mounted() {
    this.$el.addEventListener("scroll", this.scrollMonitor);
    this.$el.dataset.bottom = true;
  },
  props: ["rack"],
  updated() {
    const el = this.$el;

    if (el.dataset.bottom === "true") {
      el.scrollTop = el.scrollHeight;
    }
  }
};
</script>
