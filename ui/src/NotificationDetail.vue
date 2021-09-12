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
        <NotificationDetailTabTemplate
            v-else-if="activeKey === 'template'"
            :form="form"
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

const endpoint = '/plugin-proxy/notification/settings';

const {
  get,
  post,
} = useRequest();

export default defineComponent({
  name: 'NotificationDetail',
  components: {
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
        id: 'template',
        title: 'Template',
      }
    ]);

    const form = ref({});

    const formRef = ref();

    const onBack = () => {
      router.push(`/notifications`);
    };

    const onSave = async () => {
      await formRef.value.validate();
      await post(`${endpoint}/${id.value}`, form.value);
      ElMessage.success('Saved successfully');
    };

    onMounted(() => {
      (async () => {
        const res = await get(`${endpoint}/${id.value}`);
        form.value = res.data;
      })();
    });

    const onTabSelect = (tabName) => {
      activeKey.value = tabName;
    };

    return {
      activeKey,
      tabs,
      onBack,
      onSave,
      form,
      formRef,
      onTabSelect,
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
