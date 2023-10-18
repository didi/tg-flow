<template>
  <n-spin :show="show">
    <template #icon>
      <n-icon>
        <Reload />
      </n-icon>
    </template>
    <div class="system-container">
      <div class="system-header">
        <div class="filter-form">
         <div class="selects">
          <label>场景名称</label>
          <n-select size="small" :options="selectNameData" v-model:value="names" style="width: 500px;" multiple max-tag-count="responsive"/>
         </div>
         <div class="selects">
          <label>系统名称</label>
          <n-select size="small" :options="selectAppNameData" v-model:value="appNames" style="width: 500px;" multiple max-tag-count="responsive" />
         </div>
         <n-button size="small" type="info" @click="filterDatas">查询</n-button>
        </div>
        <div class="select-right">
          <n-button size="small"  @click="onAddScenesConfig">
          <template #icon>
            <n-icon><Add /></n-icon>
          </template>
          新增
          </n-button>
          <n-button size="small"  @click="refresh">
            <template #icon>
              <n-icon><Refresh /></n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </div>
      <n-data-table
        :columns="columns"
        :data="data"
        ref="dataTableInst"
        size="small"
        :pagination="pagination"
        :bordered="true"
        :scroll-x="1600"
      />
    </div>
  </n-spin>
</template>

<script lang="ts" setup>
import { h, onMounted, ref } from "vue";
import {
  ScenesConfigContentI,
  getScenesConfigAppname,
  getScenesConfigName,
  deleteScenesConfig,
  getScenesConfig,
} from "@/api/strategy";
import { DialogReactive, NButton, NSpace, useDialog, useMessage } from "naive-ui";
import { Add, Refresh, Reload } from "@vicons/ionicons5";
import { renderTooltip} from '@/utils/renderComponents';
import FormData from "@/views/strategy/scenes/formData.vue";
import { computed } from "@vue/reactivity";
import { valueAndLabelI } from "@/types/globalInterface";
import {useRoute} from "vue-router";

//state

const appNames = ref([] as any[]) 
const names = ref([] as any[]) 

const dataTableInst = ref(null as any) ;

let selectAppNameData = ref([] as valueAndLabelI[]);

let selectNameData = ref([] as valueAndLabelI[]);

const pagination = { pageSize: 15 }; // 表格一页显示的数据条数

const dialog = useDialog();
const route = useRoute();
const message = useMessage();

const show = ref(false);

let dialogRevise!: DialogReactive;
let dialogDelete!: DialogReactive;
let dialogAdd!: DialogReactive;

const data = ref([] as Array<ScenesConfigContentI>);

//methods
const filterDatas = () => {
  dataTableInst.value.filter({
    name:names.value,
    appname:appNames.value
  })
}

const getSelectAppNames = async () => {
  let res = await getScenesConfigAppname({
    operator: window.userInfo.username,
  });
  if (typeof res === "object") {
    selectAppNameData.value = res;
  }
};

const getSelectNames = async () => {
  let res = await getScenesConfigName({ operator: window.userInfo.username });
  if (typeof res === "object") {
    res.shift()
    selectNameData.value = res;
  }
};

/**
 * 在如果是从其他位置跳转过来的，根据参数初始化过滤条件
 */
const initFilterIfRouteQuery = () => {
  if (route.query.appName) {
    appNames.value = [route.query.appName as string]
  }
  if (route.query.sceneName) {
    names.value = [route.query.sceneName as string]
  }
  filterDatas()
}

const refresh = async () => {
  show.value = true;
  await getSelectNames();
  await getSelectAppNames();
  await initGetSystemConfig();

  initFilterIfRouteQuery()

  if (data.value) {
    message.success("刷新成功");
    show.value = false;
  } else {
    message.error("刷新失败");
    show.value = false;
  }
};

const onAddScenesConfig = () => {
  dialogAdd = dialog.info({
    title: "添加",
    style: "width: 70vw",
    showIcon: false,
    content: () =>
      h(FormData as Object, {
        scenesConfigData: {},
        onAddSuccess: addSuccess,
      }),
  });
};

const addSuccess = (row: boolean | null) => {
  if (row) {
    refresh();
    console.log("添加成功");
  } else {
    console.log("添加失败");
  }
  dialogAdd.destroy();
}; // 添加成功之后的回调函数

const onReviseConfigData = (row: ScenesConfigContentI) => {
  dialogRevise = dialog.info({
    title: "修改",
    showIcon: false,
    style: "width: 70vw",
    content: () =>
      h(FormData as Object, {
        scenesConfigData: {
          id: row.id,
          name: row.name,
          appname: row.appname,
          expname: row.expname,
          buckettype: row.buckettype,
          flow_type: row.flow_type,
        },
        onUpdateSuccess: updateSuccess,
      }),
  });
}; // 修改systemconfig

const updateSuccess = (row: boolean | null) => {
  if (row) {
    refresh();
    console.log("修改成功");
  } else {
    console.log("修改失败");
  }
  dialogRevise.destroy();
}; // 修改systemconfig 之后 FormData的回调函数

const onDeleteConfigData = (row: ScenesConfigContentI) => {
  dialogDelete = dialog.warning({
    title: "警告",
    content: "你确定要删除吗？",
    positiveText: "确定",
    negativeText: "取消",
    onPositiveClick: () => {
      onRightDelete(row.id);
    },
  });
}; // 删除config的弹窗函数

const onRightDelete = async (row: number) => {
  let res = await deleteScenesConfig({
    id: row,
    operator: window.userInfo.username,
  });
  if (res) {
    message.success("删除成功");
    await initGetSystemConfig();
  } else {
    message.success("删除失败");
  }
}; // 确认config删除的函数

const initGetSystemConfig = async () => {
  let res = await getScenesConfig({ operator: window.userInfo.username });
  if (res) {
    data.value = res;
  }
}; // 初始化config的函数

// watch && computed

let columns = computed(() => {
  return [
    {
      title: "场景编号",
      key: "id",
      sorter: (a: ScenesConfigContentI, b: ScenesConfigContentI) => a.id - b.id,
    },
    {
      title: "场景名称",
      key: "name",
      defaultFilterOptionValues: [],
      filterOptions: selectNameData.value,
      filter(value: string, row: ScenesConfigContentI) {
        return !!~row.name.indexOf(String(value));
      },
    },
    {
      title: "系统名称",
      key: "appname",
      defaultFilterOptionValues: [],
      filterOptions: selectAppNameData.value,
      filter(value: string, row: ScenesConfigContentI) {
        return !!~row.appname.indexOf(String(value));
      },
    },
    {
      title: "实验名称",
      key: "expname",
    },
    {
      title: "分流类型",
      key: "flow_type",
      sorter: (a: ScenesConfigContentI, b: ScenesConfigContentI) =>
        a.flow_type - b.flow_type,
      render(row: ScenesConfigContentI) {
        return h("p", null, { default:() => {
          if(row.flow_type === 0)return "自研分流"
          if(row.flow_type === 3)return "apollo分流"
          return row.flow_type
        }});
      },
    },
    {
      title: "操作人",
      key: "operator",
      sorter: (a: ScenesConfigContentI, b: ScenesConfigContentI) =>
        a.operator.charCodeAt(0) - b.operator.charCodeAt(0),
    },
    {
      title: "创建时间",
      key: "createtime",
      sorter: (a: ScenesConfigContentI, b: ScenesConfigContentI) =>
        Date.parse(a.createtime) - Date.parse(b.createtime),
    },
    {
      title: "更新时间",
      key: "updatetime",
      sorter: (a: ScenesConfigContentI, b: ScenesConfigContentI) =>
        Date.parse(a.updatetime) - Date.parse(b.updatetime),
    },
    {
      title: "操作",
      key: "action",
      align: "center",
      fixed: "right",
      render(row: ScenesConfigContentI) {
        return h(NSpace,null, {default:() => [
          renderTooltip('修改',() => onReviseConfigData(row),24),
          renderTooltip('删除',() => onDeleteConfigData(row),24),
        ]});
      },
    },
  ] as any;
});
// hook

onMounted(() => {
  refresh();
});
</script>

<style lang="less" scoped>
.system-container {
  width: 100%;
  .system-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 5px auto;
    .select-right{
      width: 200px;
      .n-button {
        margin: 5px;
      }
    }
    .filter-form{
      flex: 1;
      .selects{
        display: inline-flex;
        justify-content: center;
        align-items: center;
        margin: 5px;
        label{
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
