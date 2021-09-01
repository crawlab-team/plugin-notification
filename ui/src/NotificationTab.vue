<template>
  <el-form
      label-width="150px"
      inline
  >
    <el-form-item label="Name">
      <el-input v-model="name" size="small" placeholder="name"/>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue';

export default defineComponent({
  name: 'NotificationTab',
  setup(props, {emit}) {
    const name = ref('test');

    return {
      name,
    };
  },
});
</script>

<style scoped>

</style>
