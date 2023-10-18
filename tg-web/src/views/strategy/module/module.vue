<template>
    <div class="module-container">
        <div class="module-header">
            <div class="filter-form">
                <div class="selects">
                    <label>系统名称</label>
                    <n-select size="small" :options="selectAppNameData" v-model:value="appName"
                              style="width: 200px;" max-tag-count="responsive"/>
                    <label style="margin-left: 30px">代码分支</label>
                    <n-select size="small" :options="branchData" v-model:value="branchName"
                              style="width: 200px;" max-tag-count="responsive"/>
                </div>
                <n-button size="small" type="info" @click="queryModules">查询</n-button>
            </div>
        </div>
        <ModuleBox :appName="curAppName" :branch-name="curBranchName"/>
    </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";
import {getGitBranches, getScenesConfigAppname,} from "@/api/strategy";
import {NButton, useMessage} from "naive-ui";
import {valueAndLabelI} from "@/types/globalInterface";
import {useRouter} from "vue-router";
import ModuleBox from "@/views/strategy/module/ModuleBox.vue";


// ====================== data ======================

// appname
const selectAppNameData = ref([] as valueAndLabelI[]);
const appName = ref("" as string)
const curAppName = ref("" as string)

const branchName = ref("master" as string)
const curBranchName = ref("master" as string)
const branchData = ref([] as valueAndLabelI[])

// message box
const message = useMessage();
const router = useRouter();

// ====================== methods ======================

const queryModules = () => {
    curAppName.value = appName.value;
    curBranchName.value = branchName.value;
}

const getSelectAppNames = async () => {
    let res = await getScenesConfigAppname({
        operator: window.userInfo.username,
    });
    if (typeof res === "object") {
        selectAppNameData.value = res;
    }
};

const getBranchesList = async () => {
    let res = await getGitBranches({
        operator: window.userInfo.username,
        app_name: appName.value,
    })
    if (res) {
        branchData.value = res.map((item: string) => {
            return {
                label: item,
                value: item,
            }
        })
    }
}

const refresh = async () => {
    await getSelectAppNames();
};

// ====================== watch ======================

watch(appName, async (newVal, oldVal) => {
    if (newVal !== oldVal) {
        await getBranchesList();
    }
});

// ====================== hooks =======================

onMounted(() => {
    refresh();
});
</script>

<style lang="less" scoped>
.module-container {
  width: 100%;

  .module-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 5px auto;

    .filter-form {
      flex: 1;

      .selects {
        display: inline-flex;
        justify-content: center;
        align-items: center;
        margin: 5px;

        label {
          width: 60px;
          margin-right: 5px;
        }
      }
    }
  }
}

:deep(.n-button) {
  margin: auto 5px;
  text-align: left;
}
</style>
