<template>
  <div>
    <BusinessSelect :onSelectChange="onSelectChange" :businessId="businessId" />
  </div>
</template>

<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useUserInfo } from '/@/stores/userInfo';
import { useBusinessApi } from '/@/api/business/index';
import { useRouter } from 'vue-router';
import BusinessSelect from '/@/components/businessSelect/index.vue'

export default defineComponent({
  name: 'Business',
  components: {
    BusinessSelect
  },
  setup() {
    const router = useRouter();
    const businessApi = useBusinessApi();
    const businessList = ref([]);
    const userInfo = useUserInfo().userInfos;
    const businessId = ref(userInfo.businessId);

    const onSelectChange = async (businessId: number) => {
      if (!businessId) return
      await businessApi.change({ businessId });
      router.go(0);
    }

    const getBusinessList = async () => {
      const res: any = await businessApi.getList();
      if (res && res.code === 0) {
        businessList.value = res.data.list;
      }
    };
    getBusinessList();

    return {
      businessList,
      businessId,
      onSelectChange
    }
  },
});
</script>
<style scoped>
.example-showcase .el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
</style>
