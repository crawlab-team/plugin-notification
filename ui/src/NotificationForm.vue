<template>
  <cl-form :model="internalForm" ref="formRef">
    <cl-form-item :span="2" label="Name" prop="name" required>
      <el-input v-model="internalForm.name" placeholder="Name" @change="onChange"/>
    </cl-form-item>
    <cl-form-item :span="2" label="Description" prop="description">
      <el-input v-model="internalForm.description" type="textarea" placeholder="Description" @change="onChange"/>
    </cl-form-item>
    <cl-form-item :span="2" label="Type" prop="type">
      <el-select v-model="internalForm.type" @change="onChange">
        <el-option value="mail" label="Mail"/>
        <el-option value="mobile" label="Mobile"/>
      </el-select>
    </cl-form-item>
    <cl-form-item :span="1" label="Enabled" prop="enabled">
      <cl-switch v-model="internalForm.enabled" @change="onChange"/>
    </cl-form-item>
    <cl-form-item :span="1" label="Global" prop="global">
      <cl-switch v-model="internalForm.global" @change="onChange"/>
    </cl-form-item>

    <template v-if="internalForm.type === 'mail'">
      <cl-form-item :span="2" label="SMTP Server" prop="mail.server" required>
        <el-input v-model="internalForm.mail.server" placeholder="SMTP Server" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="SMTP Port" prop="mail.port" required>
        <el-input v-model="internalForm.mail.port" placeholder="SMTP Port" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="SMTP User" prop="mail.user">
        <el-input v-model="internalForm.mail.user" placeholder="SMTP User" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="SMTP Password" prop="mail.password">
        <el-input v-model="internalForm.mail.password" placeholder="SMTP Password" type="password" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="Sender Email" prop="mail.sender_email">
        <el-input v-model="internalForm.mail.sender_email" placeholder="Sender Email" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="Sender Identity" prop="mail.sender_identity">
        <el-input v-model="internalForm.mail.sender_identity" placeholder="Sender Identity" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="To" prop="mail.to">
        <el-input v-model="internalForm.mail.to" placeholder="To" @change="onChange"/>
      </cl-form-item>
      <cl-form-item :span="2" label="Cc" prop="mail.cc">
        <el-input v-model="internalForm.mail.cc" placeholder="Cc" @change="onChange"/>
      </cl-form-item>
    </template>

    <template v-else-if="internalForm.type === 'mobile'">
      <cl-form-item :span="4" label="Webhook" prop="mobile.webhook">
        <el-input
            v-model="internalForm.mobile.webhook"
            placeholder="Webhook"
            @change="onChange"
        />
      </cl-form-item>
    </template>

  </cl-form>
</template>

<script lang="ts">
import {defineComponent, ref, watch} from 'vue';

export default defineComponent({
  name: 'NotificationForm',
  props: {
    modelValue: {
      type: Object,
      default: () => {
        return {};
      }
    },
  },
  emits: [
    'update:modelValue',
  ],
  setup(props, {emit}) {
    const formRef = ref();

    const internalForm = ref({
      name: '',
      description: '',
      type: 'mail',
      enabled: true,
      global: true,
      mail: {
        server: '',
        port: '465',
        user: '',
        password: '',
        sender_email: '',
        sender_identity: '',
        title: '',
        template: '',
        to: '',
        cc: '',
      },
      mobile: {
        webhook: '',
        title: '',
        template: '',
      },
    });

    watch(() => props.modelValue, () => {
      internalForm.value = props.modelValue;
    });

    const onChange = () => {
      emit('update:modeValue', internalForm.value);
    };

    const validate = async () => {
      await formRef.value.validate();
    };

    return {
      formRef,
      internalForm,
      onChange,
      validate,
    };
  },
});
</script>

<style scoped>

</style>
