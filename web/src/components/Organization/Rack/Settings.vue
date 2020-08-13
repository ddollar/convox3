<template>
  <b-modal :id="`rack-settings-${rid}`" title="test">{{ rid }}</b-modal>
</template>

<script>
export default {
  props: ["rid"],
};
</script>
