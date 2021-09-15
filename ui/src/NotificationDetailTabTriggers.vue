<template>
  <cl-transfer
      :titles="titles"
      :data="data"
      :value="enabled"
      @change="onChange"
  />
</template>

<script lang="ts">
import {computed, defineComponent, ref} from 'vue';

export default defineComponent({
  name: 'NotificationDetailTabTriggers',
  props: {
    form: {
      type: Object,
      default: () => {
        return {};
      },
    },
  },
  emits: [
    'change',
  ],
  setup(props, {emit}) {
    const titles = ref([
      'Available',
      'Enabled',
    ]);

    const data = ref([
      {
        key: 'model:tasks:add',
        label: 'Create Task',
      },
      {
        key: 'model:tasks:save',
        label: 'Update Task',
      },
    ]);

    const enabled = computed(() => {
      const {triggers} = props.form;
      if (!triggers) return [];
      return triggers.map(t => t.event);
    });

    const onChange = (value) => {
      emit('change', value);
    };

    return {
      titles,
      data,
      enabled,
      onChange,
    };
  },
});
</script>

<style scoped>

</style>
