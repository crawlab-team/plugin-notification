<template>
  <div class="notification-detail-tab-template">
    <textarea :value="form.template" class="simple-mde" ref="simpleMDERef"/>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onBeforeUnmount, ref, computed} from 'vue';

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
  emits: [
    'template-change',
    'title-change',
  ],
  setup(props, {emit}) {
    const simpleMDERef = ref();

    const simpleMDE = ref();

    onMounted(() => {
      simpleMDE.value = new window.SimpleMDE({
        element: simpleMDERef.value,
        spellChecker: false,
      });
      simpleMDE.value.codemirror.on('change', () => {
        emit('template-change', simpleMDE.value.value());
      });
    });

    onBeforeUnmount(() => {
      simpleMDE.value.toTextArea();
      simpleMDE.value = null;
    });

    return {
      simpleMDERef,
    };
  },
});
</script>

<style scoped>
</style>
