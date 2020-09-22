<template>
  <div>
    <b-button variant="secondary" size="sm" v-b-modal="`manifest-${this._uid}`">
      <i class="far fa-file-alt" />
    </b-button>
    <b-modal :id="`manifest-${this._uid}`" hide-footer :title="title">
      <pre class="text-light bg-dark p-3 mb-0">{{ manifest }}</pre>
    </b-modal>
  </div>
</template>

<script>
export default {
  props: ["manifest", "title"],
};
</script>