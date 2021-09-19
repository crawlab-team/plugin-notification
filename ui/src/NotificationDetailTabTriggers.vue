<template>
  <cl-transfer
      :titles="titles"
      :data="triggerList"
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
    triggerList: {
      type: Array,
      default: () => {
        return [];
      },
    }
  },
  emits: [
    'change',
  ],
  setup(props, {emit}) {
    const titles = ref([
      'Available',
      'Enabled',
    ]);

    const enabled = computed(() => {
      const {triggers} = props.form;
      return triggers || [];
    });

    const onChange = (value) => {
      emit('change', value);
    };

    return {
      titles,
      enabled,
      onChange,
    };
  },
});
</script>

<style scoped>

</style>
