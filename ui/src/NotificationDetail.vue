<template>
  <div class="detail-layout notification-detail">
    <div class="content">
      <cl-nav-tabs
          :active-key="activeKey"
          :items="tabs"
          @select="onTabSelect"
      />
      <cl-nav-actions
          class="nav-actions"
          :collapsed="false"
      >
        <cl-nav-action-group-detail-common
            @back="onBack"
            @save="onSave"
        />
      </cl-nav-actions>

      <div class="content-container">
        <template v-if="activeKey === 'overview'">
          <NotificationForm
              ref="formRef"
              v-model="form"
          />
        </template>
        <NotificationDetailTabTriggers
            v-else-if="activeKey === 'triggers'"
            :form="form"
            @change="onTriggersChange"
        />
        <NotificationDetailTabTemplate
            v-else-if="activeKey === 'template'"
            :form="form"
            @title-change="onTitleChange"
            @template-change="onTemplateChange"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, ref} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import {useRequest} from 'crawlab-ui';
import {ElMessage} from 'element-plus';
import NotificationForm from './NotificationForm.vue';
import NotificationDetailTabTemplate from './NotificationDetailTabTemplate.vue';
import NotificationDetailTabTriggers from './NotificationDetailTabTriggers.vue';

const endpoint = '/plugin-proxy/notification/settings';

const {
  get,
  post,
} = useRequest();

export default defineComponent({
  name: 'NotificationDetail',
  components: {
    NotificationDetailTabTriggers,
    NotificationForm,
    NotificationDetailTabTemplate,
  },
  setup() {
    const router = useRouter();

    const route = useRoute();

    const id = computed(() => route.params.id);

    const activeKey = ref('overview');

    const tabs = ref([
      {
        id: 'overview',
        title: 'Overview',
      },
      {
        id: 'triggers',
        title: 'Triggers',
      },
      {
        id: 'template',
        title: 'Template',
      },
    ]);

    const form = ref({});

    const formRef = ref();

    const onBack = () => {
      router.push(`/notifications`);
    };

    const onSave = async () => {
      if (formRef.value) await formRef.value.validate();
      await post(`${endpoint}/${id.value}`, form.value);
      ElMessage.success('Saved successfully');
    };

    onMounted(() => {
      (async () => {
        const res = await get(`${endpoint}/${id.value}`);
        const {data} = res;
        form.value = data;
      })();
    });

    const onTabSelect = (tabName) => {
      activeKey.value = tabName;
    };

    const allTriggers = ref([
      {
        event: 'model:tasks:add',
        name: 'Create Task',
      },
      {
        event: 'model:tasks:save',
        name: 'Update Task',
      },
    ]);

    const onTriggersChange = (value) => {
      const triggers = value.map(v => allTriggers.value.find(d => d.event === v));
      form.value.triggers = [].concat(triggers);
    };

    const onTitleChange = (value) => {
      form.value.title = value;
    };

    const onTemplateChange = (value) => {
      form.value.template = value;
    };

    return {
      activeKey,
      tabs,
      onBack,
      onSave,
      form,
      formRef,
      onTabSelect,
      onTriggersChange,
      onTitleChange,
      onTemplateChange,
    };
  },
});
</script>

<style scoped>
.detail-layout {
  display: flex;
  height: 100%;
}

.detail-layout .content {
  flex: 1;
  max-width: 100%;
  background-color: white;
  display: flex;
  flex-direction: column;
}

.detail-layout .content .content-container {
  margin: 20px;
}
</style>
