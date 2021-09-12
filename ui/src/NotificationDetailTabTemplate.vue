<template>
  Template
</template>

<script lang="ts">
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'NotificationDetailTabTemplate',
  props: {
    form: {
      type: Object,
      default: () => {
        return {};
      },
    },
  },
  setup() {
    return {};
  },
});
</script>

<style scoped>

</style>
