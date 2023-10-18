<template>
  <n-spin :show="show">
    <template #icon>
      <n-icon>
        <Reload/>
      </n-icon>
    </template>
    <div class="system-container">
      <div class="system-header">
        <div class="filter-form">

        </div>
        <div class="select-right">
          <n-button size="small" @click="addSystemConfig">
            <template #icon>
              <n-icon>
                <Add/>
              </n-icon>
            </template>
            新增
          </n-button>
          <n-button size="small" @click="refresh">
            <template #icon>
              <n-icon>
                <Refresh/>
              </n-icon>
            </template>
            刷新
          </n-button>
        </div>
      </div>
      <n-data-table size="small"
          :columns="column"
          :data="data"
          :pagination="pagination"
          :bordered="true"
          :scroll-x="1600"
      />
    </div>
  </n-spin>
</template>

<script lang='ts' setup>
import {h, onMounted, ref, reactive} from 'vue'
import {getSystemConfig, SystemConfigContentI, exportSystemConfig, deleteSystemConfig} from '@/api/strategy'
import {DialogReactive, NButton, NSpace, useDialog, useMessage} from 'naive-ui';
import {DataTableColumns,} from 'naive-ui';
import {renderTooltip} from "@/utils/renderComponents";
import {Add, Refresh, Reload} from '@vicons/ionicons5';
import FormData from '@/views/strategy/system/formData.vue';
import JsZip from 'jszip';
import saveAs from "jszip/vendor/FileSaver.js";
import {ElLink} from "element-plus";
import SvgIcon from "@/components/svgIcon/index.vue";

//state
const column = reactive<DataTableColumns<any>>([
  {
    title: "系统编号",
    key: "id",
    sorter: (a: SystemConfigContentI, b: SystemConfigContentI) => a.id - b.id
  },
  {
    title: "系统名称",
    key: "app_name"
  },
  {
    title: "部署机房",
    key: "machine_room",
    sorter: (a: SystemConfigContentI, b: SystemConfigContentI) => a.machine_room.charCodeAt(0) - b.machine_room.charCodeAt(0)
  },
  {
    title: "节点名称",
    key: "node_name",
    sorter: (a: SystemConfigContentI, b: SystemConfigContentI) => a.node_name.charCodeAt(0) - b.node_name.charCodeAt(0)
  },
  {
    title: "操作人",
    key: "operator",
    sorter: (a: SystemConfigContentI, b: SystemConfigContentI) => a.operator.charCodeAt(0) - b.operator.charCodeAt(0)
  },
    // render as an el-link and icon
  {
    title: "Git 仓库",
    key: "git_url",
    render(row: SystemConfigContentI) {
      return h(ElLink, {
        href: row.git_url,
        underline: false,
        type: "primary"
      }, {
          default: () => h(NSpace, {
              size: "small"
          }, {
            default: () => {
                if (row.git_url === '') {
                    return h('span', '')
                } else {
                    return [
                        h(SvgIcon, {name: 'ele-Link'}),
                        h('span', row.git_url),
                    ]
                }
            }
          })
      })
    }
  },
  {
    title: "创建时间",
    key: "create_time",
    sorter: (a: SystemConfigContentI, b: SystemConfigContentI) => Date.parse(a.create_time) - Date.parse(b.create_time)
  },
  {
    title: "更新时间",
    key: "update_time",
    sorter: (a: SystemConfigContentI, b: SystemConfigContentI) => Date.parse(a.update_time) - Date.parse(b.update_time)
  },
  {
    title: "操作",
    key: "action",
    align: "center",
    fixed: "right",
    render(row: SystemConfigContentI) {
      return h(NSpace, null, {
        default: () => [
          renderTooltip('修改', () => onReviseConfigData(row), 24),
          renderTooltip('导出', () => exportConfigData(row), 24),
          renderTooltip('删除', () => onDeleteConfigData(row), 24),
        ]
      })
    }
  }
])

const pagination = {pageSize: 15}// 表格一页显示的数据条数

const dialog = useDialog();

const message = useMessage();

const show = ref(false)

let dialogCommit!: DialogReactive;
let dialogDelete!: DialogReactive;
let DialogExport!: DialogReactive;
let dialogAdd!: DialogReactive;

const data = ref([] as Array<SystemConfigContentI>)

//methods

const refresh = async () => {
  show.value = true
  await initGetSystemConfig()
  if (data.value) {
    message.success('刷新成功')
    show.value = false
  } else {
    message.error('刷新失败')
    show.value = false
  }
}

const addSystemConfig = () => {
  dialogAdd = dialog.info({
    title: "添加",
    showIcon: false,
    style: "width: 70vw",
    content: () => h(FormData as Object, {
      systemConfigData: {},
      onAddSuccess: addSuccess
    })
  })
}

const addSuccess = (row: boolean | null) => {
  if (row) {
    initGetSystemConfig();
    console.log('添加成功');
  } else {
    console.log('添加失败');
  }
  dialogAdd.destroy();
}

const onReviseConfigData = (row: SystemConfigContentI) => {
  dialogCommit = dialog.info({
    title: "修改",
    showIcon: false,
    style: "width: 70vw",
    content: () => h(FormData, {
      systemConfigData: {
        app_name: row.app_name,
        id: row.id,
        machine_room: row.machine_room,
        node_name: row.node_name,
        git_url: row.git_url,
      },
      onUpdateSuccess: updateSuccess,
    }),
  })
} // 修改systemconfig

const updateSuccess = (row: boolean | null) => {
  if (row) {
    initGetSystemConfig();
    console.log('修改成功');
  } else {
    console.log('修改失败');
  }
  dialogCommit.destroy();

} // 修改systemconfig 之后 FormData的回调函数

const onDeleteConfigData = (row: SystemConfigContentI) => {
  dialogDelete = dialog.warning({
    title: '警告',
    content: '你确定要删除吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      onRightDelete(row.id);
    },
  })
}// 删除config的弹窗函数

const onRightDelete = async (row: number) => {
  let res = await deleteSystemConfig({
    id: row,
    operator: window.userInfo.username
  })
  if (typeof res === "object") {
    message.success("删除成功")
    await initGetSystemConfig();
  } else {
    message.success(res as string)
  }
} // 确认config删除的函数

const exportConfigData = (row: SystemConfigContentI) => {
  DialogExport = dialog.success({
    title: '导出',
    content: '你确定要导出数据吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      rightExport(row.id);
    },
  })
} // 导出数据的弹窗函数

const rightExport = async (row: number) => {
  let res = await exportSystemConfig({
    id: row,
    operator: window.userInfo.username
  })
  if (typeof res === "object") {
    let timeStamp = Date.parse(String(new Date()));
    let data = JSON.parse(res.content)
    let zip = new JsZip();
    let promises = []
    let promise = new Promise((res) => {
      let file_name = "scene.json"
      zip.file(file_name, JSON.stringify(JSON.parse(data["scene"]), null, "\t"))
      res(null)
    })
    promises.push(promise)
    let workflowList = eval(data["workflowList"])
    await workflowList.forEach((item: any) => {
      let workflow = JSON.parse(item)
      let promise = new Promise((res) => {
        let sceneName = workflow["scene_name"]
        let interfaceName = sceneName.split("-")[0]
        let file_name = `${interfaceName}/workflow-${workflow["scene_id"]}-${workflow["id"]}-${sceneName}.json`
        zip.file(file_name, JSON.stringify(workflow["flow_charts"], null, "\t"))
        res(null)
      })
      promises.push(promise)
    });
    Promise.all(promises).then(() => {
      zip.generateAsync({
        type: "blob"
      }).then((content) => {
        saveAs(content, `scene_${timeStamp}.zip`)
      })
    })
    message.success("导出成功")

  } else {
    message.error("导出失败")
  }
} // 确认导出数据的好函数

const initGetSystemConfig = async () => {
  let res = await (getSystemConfig({operator: window.userInfo.username}))
  if (typeof res === 'object') {
    data.value = res;
  }
} // 初始化config的函数

// watch && computed

// hook

onMounted(() => {
  initGetSystemConfig();
})

</script>

<style lang='less' scoped>
.system-container {
  width: 100%;

  .system-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 5px auto;
  }

  .n-button {
    margin: 5px;
  }
}

:deep(.n-button) {
  margin: auto 5px auto;
  text-align: left;
}

</style>
