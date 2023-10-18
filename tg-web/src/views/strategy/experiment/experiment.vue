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
          <div class="select-out">
            <div class="select">
              <label for="systemId">系统:</label>
              <n-select
                size="small"
                id="systemId"
                style="width: 200px"
                filterable
                clearable
                v-model:value="formValue.app_id"
                :on-update:value="onAppidSuccess"
                :options="app_idData"
              />
            </div>
          </div>
          <div class="select-out">
            <div class="select">
              <label for="scenesName">场景名称:</label>
              <n-select
                size="small"
                id="scenesName"
                style="width: 300px"
                clearable
                filterable
                v-model:value="formValue.scene_id"
                :on-update:value="onSceneidSuccess"
                :options="scene_idData"
              />
            </div>
          </div>
          <div class="select-out">
            <n-button size="small" type="info" @click="initGetExperimentConfig">
              查询
            </n-button>
          </div>
        </div>
        <div class="select-out">
          <n-button size="small" @click="addScenesConfig">
            <template #icon>
              <n-icon><Add /></n-icon>
            </template>
            新增
          </n-button>
          <n-button size="small" @click="refresh">
            <template #icon>
              <n-icon><Refresh /></n-icon>
            </template>
            刷新
          </n-button>
          <n-button size="small" @click="onHistoryVersion">
            <template #icon>
              <n-icon><EllipsisVerticalOutline /></n-icon>
            </template>
            历史版本
          </n-button>
          <n-button size="small" @click="refreshNode">
            <template #icon>
              <n-icon><Refresh /></n-icon>
            </template>
            刷新节点配置
          </n-button>
        </div>
      </div>
      <n-data-table size="small"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :bordered="true"
        :scroll-x="1400"
      />
    </div>
  </n-spin>
</template>

<script lang="ts" setup>
import { h, onBeforeMount, ref } from "vue";
import {
  experimentConfigResContentI,
  getExperimentConfigAppname,
  getExperimentScenes,
  getExperimentConfig,
  getExperimentConfigDimension,
  getExperimentWorkflowNodeType,
  addOrUpdateWorkflow,
  deleteWorkflowI,
  deleteWorkflowData,
  exportexperimentConfig,
  getWorkFlowChartByFlowIdV2,
  refreshExperimentConfigNodeType,
  requireWorkflowConfig,
} from "@/api/strategy";
import { DialogReactive, NSpace, NIcon, useDialog, useMessage } from "naive-ui";
import {
  Add,
  Refresh,
  Reload,
  EllipsisVerticalOutline,
} from "@vicons/ionicons5";
import { renderTooltip} from '@/utils/renderComponents';
import { computed } from "@vue/reactivity";
import historyVersionVue from "./historyVersion.vue";
import visVue from "./vis.vue";
import requireConfigVue from './requireConfig.vue';
import updateFormVue from "./updateForm.vue";
import { valueAndLabelI } from "@/types/globalInterface";
import showX6paramsVue from "./showX6params.vue";
import {useRoute} from "vue-router";
// interface
interface selectModelConfigI {
  app_id: string;
  scene_id: string;
}

//state
let formValue = ref({
  app_id: "",
  scene_id: "",
} as selectModelConfigI);
// 获取页面基本数据的请求参数。 

let app_idData = ref([] as Array<valueAndLabelI>);
const initExperimentConfigAppname = async () => {
  const res = await getExperimentConfigAppname();
  if (res) {
    if (typeof res !== "string") {
      app_idData.value = res;
    }
  }
}; //初始化appname

let scene_idData = ref([] as Array<valueAndLabelI>);
const initExperimentConfigName = async (types: number) => {
  const sceneList = await getExperimentScenes(types);
  if (sceneList) {
    let scene: string = "";
    if (typeof sceneList !== "string") {
      sceneList.forEach((e) => {
        if (scene) {
          scene = `${scene},${e.value}`;
        } else {
          scene = `${e.value}`;
        }
      });
      scene_name.value = sceneList[0].label;
      scenes.value = scene;
      initScenes = scene;
      scene_idData.value = sceneList;
    }
  }
}; // 初始化场景名称

const pagination = { pageSize: 15 }; // 表格一页显示的数据条数

const dialog = useDialog();
const route = useRoute();
const message = useMessage();

let dialogRevise!: DialogReactive;
let dialogAdd!: DialogReactive;
let dialogHistory!: DialogReactive;
let dialogRequire!: DialogReactive;

let scenes = ref("");
let tempScenes = ""; // 存的是，当系统下拉框选完后，获取的全部场景名称名字。
let dimension_id = ref("-1");// 默认就是-1
let scene_name = ref(""); // 选全部时，则为tempScenes，选具体场景时，则是tempScenesValue
let tempScenesValue = ""; // 场景名字的临时值，因为当场景选“全部”的时候，并不能保存到场景名称的全部值。
let initScenes = ""; // 初始化时，scenes
let tempRow: any = ref("");// 点复制的时候，需要取一个作为
let tempWorkflowId: string = ""; // 将workflowid 存起来
let prefix: string = ""; // 流程图前缀
let suffix: number[] = []; // 流程图后缀

const data = ref([] as Array<experimentConfigResContentI>);

//methods
const show = ref(false); // spin

/**
 * 通过路由跳转过来时，设置初始化过滤条件
 */
const initGetExperimentConfigWithRoute = async () => {

    if (route.query.appName) {
        let appData = app_idData.value.find(
            (item) => item.label === route.query.appName
        )
        if (appData) {
            formValue.value.app_id = route.query.appName as string;
            await onAppidSuccess(appData.value, appData)
        }
    }

    await initGetExperimentConfig();

    if (route.query.workflowId) {
        const workflowId = route.query.workflowId as string;
        data.value = data.value.filter(
            (item) => item.workflowid === workflowId
        )
    }

}; // 初始化config的函数


const refresh = async () => {
  show.value = true;
  await initGetExperimentConfigWithRoute();
  if (data.value) {
    message.success("刷新成功");
  } else {
    message.error("刷新失败");
  }
  show.value = false;
}; // 触发刷新

const refreshNode = async () => {
  const res = await refreshExperimentConfigNodeType({
    operator: window.userInfo.username,
    systemName: app_idData.value.reduce((pre, cur) => {
      pre = `${pre},${cur.label}`;
      return pre;
    }, ""),
  });
  if (res) {
    if (res === true) {
      message.success("节点更新成功");
    } else {
      message.error(`节点更新失败,${res}`);
    }
  } else message.error("节点更新失败");
}; // 刷新节点配置

const dialogError = () => {
  dialog.error({
    title: "错误",
    content: "请在具体场景下进行该操作",
    positiveText: "OK",
  });
}; // 公共错误函数，三个地方都需要用

const onHistoryVersion = async () => {
  if (onClickError()) {
    dialogHistory = dialog.info({
      showIcon:false,
      content: () =>
        h(historyVersionVue as Object, {
          scenename: scenes.value,
          dimension_id: Number(dimension_id.value),
          onBackHistory: backHistory,
        }),
    });
  } else {
    dialogError();
  }
}; // 通过具体场景来发送请求

const backHistory = (res: boolean) => {
  if (res) {
    message.success("回退版本成功");
    refresh();
  } else {
    message.error("回退版本失败");
  }
  dialogHistory.destroy();
}; // 回滚某个版本成功

const onCopyConfigData = async (row: experimentConfigResContentI) => {
  if (onClickError()) {
    let res = await addOrUpdateWorkflow({
      operator: window.userInfo.username,
      scenename: row.scenename,
      dimension_id: row.dimension_id,
      oldworkflowid: row.oldworkflowid,
      workflowid: "",
      modules: row.modules,
      proportion: 1,
      defult: 0,
      remark: row.remark,
    });
    if (res === true) {
      message.success("复制成功");
      await refresh();
    } else {
      message.error(`复制失败,${res}`);
    }
  } else {
    dialogError();
  }
}; // 复制功能

const onImportConfigData = (row: experimentConfigResContentI) => {
  dialog.info({
    title: "导出",
    showIcon:false,
    content: "你确定要导出数据吗？",
    positiveText: "确定",
    negativeText: "取消",
    onPositiveClick: () => {
      rightExport(row.workflowid);
    },
  });
}; // 导出功能

const onRequireConfigData = (id: string) => {
  dialogRequire = dialog.info({
    title:'导入流程图',
    showIcon:false,
    content:() => h(requireConfigVue,{
      id,
      onReset:reset,
      onSuccess:success,
    }),

  })
}; // 导入功能

const reset = () => dialogRequire.destroy()

const success = async (json: string, workflowid: number) => {
  const res = await requireWorkflowConfig({dataJson:json,workflowid})
  if(typeof res === 'boolean'){
    message.success('导入成功')
    dialogRequire.destroy()
    await refresh()
  }else{
    message.warning(res || '导入失败')
  }
}

const rightExport = async (workflowid: string) => {
  const res = await exportexperimentConfig({
    operator: window.userInfo.username,
    workflowid,
  });
  if (res) {
    let data: any = JSON.parse(res);
    let content = JSON.stringify(data["flow_charts"], null, "\t");
    let filename =
      "workflow-" +
      data["scene_id"] +
      "-" +
      data["id"] +
      "-" +
      data["scene_name"] +
      ".json";
    let blob = new Blob([content], { type: "text/json" });
    let a = document.createElement("a");
    a.download = filename;
    a.href = window.URL.createObjectURL(blob); //生成的URL
    a.dataset.downloadurl = ["text/json", a.download, a.href].join(":");
    a.click();
    window.URL.revokeObjectURL(a.href); // 需要手动释放URL 因为此时这个文件还存在内存中
  } else {
    message.error(`导出失败${res}`);
  }
}; // 确认导出

const types = ref(0); // appId的映射
const onAppidSuccess = async (e: string, f: any) => {
  types.value = Number(e);
  await initExperimentConfigName(types.value);
  formValue.value.app_id = f.label;
  formValue.value.scene_id = scene_idData.value[0].label;
  tempScenes = scenes.value;
  tempScenesValue = '';
}; // 系统名字选择

const onSceneidSuccess = async (e: string, f: any) => {
  if (e === "全部") {
    if (tempScenes) {
      scenes.value = tempScenes;
    } else {
      scenes.value = initScenes;
    }
    tempScenesValue = "";
  } else {
    tempScenesValue = f.label;
  }
  formValue.value.scene_id = f.label;
  scene_name.value = f.label;
}; // 场景选择

const onClickError = (): boolean => {
  return !!(scenes.value !== "全部" && !scenes.value.includes(",") && scenes.value);
}; // 复制删除历史版本新增点击之后的公共判断

const addScenesConfig = async () => {
  if (onClickError()) {
    for (let i in tempRow.value) {
      tempRow.value[i] = "";
    }
    tempRow.value["defult"] = "从";
    tempRow.value['proportion'] = '0%';
    tempRow.value['scenename'] = data.value[0].scenename;
    tempRow.value['dimension_id'] = data.value[0].dimension_id;
    const nodeType = await getExperimentWorkflowNodeType({
      operator: window.userInfo.username,
      scene_name: scene_name.value,
    });
    dialogAdd = dialog.info({
      style: { width: "100%", height: "100vh" },
      title: "添加",
      showIcon:false,
      content: () =>
        h(
          visVue as Object,
          {
            operator: "添加",
            nodeType,
            prefix: `action-${data.value[0].workflowid}`,
            onAddSuccess: addSuccess,
          },
          {
            default: () =>
              h(updateFormVue as object, {
                message: tempRow.value,
              }),
          }
        ),
    });
  } else {
    dialogError();
  }
}; // 添加函数

const addSuccess = async (res: string) => {
  if (res) {
    const data = await addOrUpdateWorkflow({
      ...tempRow.value,
      dataJson: res,
      operator: window.userInfo.username,
    });
    if (data === true) {
      message.success("添加成功");
      await refresh();
    } else {
      message.error(`添加错误-${data}`);
    }
  } else {
    message.error("添加失败");
  }
  dialogAdd.destroy();
}; // 添加成功的回调函数

const showWorkflow = async (row: any) => {
  // const data = await byWorkflowIdToDataNodes(row.workflowid);
  let data = await getComboDataNodes(row.workflowid);
  dialog.info({
    style: { width: "100%", height: "100vh" },
    title: `场景编号:${row.scene_id} -> 流量占比:${row.proportion} -> 主/从:${row.defult}`,
    showIcon: false,
    content: () => {
      return h(showX6paramsVue as Object, {
        dataNodes: data,
        style: {
          width: "100%",
          height: "calc(100vh - 100px)",
        },
      });
    },
  });
}; // 通过workflow展示流程图

let enmuType:any = {
    'condition': 'diamond',
    'flow': 'circle',
    'task': 'rect',
    'timeout': 'clock',
}

const byWorkflowIdToDataNodes = async (workflowId: string): Promise<any> => {
  let data = await getWorkFlowChartByFlowIdV2({
    operator: window.userInfo.username,
    workflowid: Number(workflowId),
  });
  if (data === '') {
    data = {};
  } 
  suffix = [];
  if (!data.edges) {
    data.edges = [];
  }
  if (!data.nodes) {
    data.nodes = [];
  }
  if(data.flow_charts?.actions){
    for (let k in data.flow_charts.actions) {
      const n = data.flow_charts.actions[k]
      const node: any = {}
      node.width = 200
      node.height = 80
      node.label = n.action_name
      node.title = n.action_name
      node.id = n.action_id
      node.type = n.action_type
      node.params = n.params
      node.shape = n.action_type ? enmuType[n.action_type] : 'dag-node'
      node.timeout = n.timeout
      node.timeout_async = n.timeout_async
      node.timeout_dynamic = n.timeout_dynamic
      node.ref_workflow_id = n.ref_workflow_id
      node.conditionSelectValue = n.action_type == 'condition' ? n.action_name : undefined
      node.x = parseInt(n.location.split(",")[0])
      node.y = parseInt(n.location.split(",")[1])

      node.data = {}
      node.data.label = n.action_name
      node.data.id = n.action_id
      node.data.type = n.action_type
      node.data.shape = n.action_type ? enmuType[n.action_type] : 'dag-node'
      node.data.timeout = n.timeout
      node.data.timeout_async = n.timeout_async
      node.data.timeout_dynamic = n.timeout_dynamic
      node.data.ref_workflow_id = n.ref_workflow_id
      node.data.conditionSelectValue = n.action_type == 'condition' ? n.action_name : undefined

      data.nodes.push(node)

      if(n.action_type === 'condition') {
        for(let i in n.next_action_ids){
          const edge: any = {}
          edge.id = n.action_id + '_' + n.next_action_ids[i]
          edge.type = 'line'
          edge.source = n.action_id
          edge.target =  n.next_action_ids[i]
          edge.label = n.next_conditions[i]
          data.edges.push(edge)
        }
        node.params = n.params
      } else {
          for(let key in n.next_action_ids){
            const edge: any = {}
            const ac_id = n.next_action_ids[key]
            edge.id = n.action_id + '_' + ac_id
            edge.type = 'line'
            edge.source = n.action_id
            edge.target = ac_id
            data.edges.push(edge)
          }
      }
    }

    if (data.nodes.length > 0) {
      let sliceIdx = 8 + data.nodes[0].id.split("-")[1].length;

      data.nodes = data.nodes
        .map((n: any) => {
          suffix.push(n.id.slice(sliceIdx));
          return {
            ...n,
            title: {text: n.label },
          };
        })
        .sort((a: any, b: any) => {
          return a.id.slice(sliceIdx) - b.id.slice(sliceIdx);
        });
    }
  }
  // console.log(data)
  return data;
}; // 请求dataNodes

const getComboDataNodes = async (workflowId: string) => {
  let data = await byWorkflowIdToDataNodes(workflowId);
  if (!data.combos) {
    data.combos = [];
  }
  for (let k in data.nodes) {
    const node = data.nodes[k];
    if(node.type == 'flow'){
      data = await addComboDataNodes(data, node, '');
    }
  }
  return data;
}

const addComboDataNodes = async (data: any, node: any, parentId: string) => {
  // let res = Object.assign({}, data)
  // let res = Object.create(data);
  const combo: any = {};
  combo.id = 'combo_' + node.id;
  // combo.id = node.id;
  combo.collapsed = true;
  combo.x = node.x;
  combo.y = node.y;
  if (parentId !== '') {
    combo.parentId = parentId;
  }

  // res.combos.push(combo);

  let flowData = await byWorkflowIdToDataNodes(node.data.ref_workflow_id);
  if (flowData.nodes.length !== 0) {
    data.combos.push(combo);
  }
  // console.log(flowData)
  let left = flowData.nodes[0]?.x ? flowData.nodes[0].x : node.x;
  let right = flowData.nodes[0]?.x ? flowData.nodes[0].x : node.x;
  let top = flowData.nodes[0]?.y ? flowData.nodes[0].y : node.y;
  let bottom = flowData.nodes[0]?.y ? flowData.nodes[0].y : node.y;
  flowData.nodes.forEach((n: { x: number; y: number; }) => {
    if(n.x < left){
      left = n.x;
    }else if(n.x > right){
      right = n.x;
    }
    if(n.y < top){
      top = n.y;
    }else if(n.y > bottom){
      bottom = n.y;
    }
  })
  let x_avg = (left + right) / 2;
  let y_avg = (top + bottom) / 2;

  flowData.nodes?.forEach((n: { comboId: string; x: number; y: number; type: string; }) => {
    n.comboId = 'combo_' + node.id;
    // n.comboId = node.id;
    n.x = node.x + n.x - x_avg - 70;
    n.y = node.y + n.y - y_avg + 18;
    data.nodes.push(n);
    // res.nodes.push(n);
    
    if(n.type == 'flow') {
      addComboDataNodes(data, n, n.comboId);
    }
  })
  flowData.edges.forEach((e: { comboId: any; }) =>{
    e.comboId = node.id;
    data.edges.push(e);
  })

  // const flow_edge : any = {}
  // flow_edge.id = node.id + '_' + 'combo_' + node.id;
  // flow_edge.type = 'line'
  // flow_edge.source = node.id
  // flow_edge.target = 'combo_' + node.id;
  // flow_edge.comboId = node.id;
  // data.edges.push(flow_edge)
  return data;
}


const onUpdateConfigData = async (row: experimentConfigResContentI) => {
  // let data = await byWorkflowIdToDataNodes(row.workflowid);
  let data = await getComboDataNodes(row.workflowid);
  tempRow.value = row;
  tempWorkflowId = row.workflowid;
  prefix = `action-${tempWorkflowId}`;
  const nodeType = await getExperimentWorkflowNodeType({
    operator: window.userInfo.username,
    workflowid: Number(row.workflowid),
  });
  if (Number.isNaN(Number(suffix.sort((a, b) => a - b)[suffix.length - 1]))) // 如果suffix中最后一位不是数字，则说明为空。
    suffix[0] = 0;
  if (nodeType) {
    dialogRevise = dialog.info({
      title: `场景编号:${row.scene_id} -> 流量占比:${row.proportion} -> 主/从:${row.defult}`,
      style: { width: "100vw", height: "100vh", padding: "5px", margin: 0, maxWidth: "100vw"},
      maskClosable: false,
      showIcon: false,
      content: () =>
        h(
          visVue as Object,
          {
            dataNodes: data,
            nodeType: nodeType,
            operator: "修改",
            workflowId: row.workflowid,
            prefix,
            suffix: Number(suffix[suffix.length - 1]),
            onUpdateSuccess: updateSuccess,
            onResetWorkflow: resetWorkflow,
          },
          {
            default: () =>
              h(updateFormVue, {
                message: row,
              }),
          }
        ),
    });
  } else {
    if (data) {
      message.error(`nodeType获取失败${nodeType}`);
    } else {
      message.error(`nodes获取失败${nodeType}`);
    }
  }
}; // 修改systemconfig

const resetWorkflow = () =>  dialogRevise.destroy() // 关闭流程图弹窗

const updateSuccess = async (res: string) => {
  if (res) {
    const data = await addOrUpdateWorkflow({
      ...tempRow.value,
      dataJson: res,
      oldworkflowid: tempWorkflowId,
      operator: window.userInfo.username,
    });
    if (data === true) {
      message.success("修改成功");
      await refresh();
    } else {
      message.error(`修改错误-${data}`);
    }
  } else {
    message.error("修改失败");
  }
  dialogRevise.destroy();
}; // 修改systemconfig 之后 FormData的回调函数

const onDeleteConfigData = (row: experimentConfigResContentI) => {
  if (onClickError()) {
    dialog.warning({
      title: "警告",
      content: "你确定要删除吗？",
      positiveText: "确定",
      negativeText: "取消",
      onPositiveClick: () => {
        onRightDelete({
          operator: window.userInfo.username,
          workflowid: row.workflowid,
          scenename: row.scenename,
          dimension_id: row.dimension_id,
        });
      },
    });
  } else {
    dialogError();
  }
}; // 删除config的弹窗函数

const onRightDelete = async (params: deleteWorkflowI) => {
  let res = await deleteWorkflowData(params);
  if (res === true) {
    message.success("删除成功");
    await initGetExperimentConfig();
  } else {
    message.success(`删除失败,${res}`);
  }
}; // 确认config删除的函数

const initGetExperimentConfig = async () => {
  if (tempScenesValue) {
    scenes.value = tempScenesValue;
  }
  await initExperimentDimension(scene_name.value);
  let res = await getExperimentConfig({
    operator: window.userInfo.username,
    dimension_id: dimension_id.value,
    scenename: scenes.value,
  });
  if (res) {
    if (typeof res === "string") {
      message.error("查询失败");
    } else {
      data.value = res;
      if(res.length > 0){
        tempRow.value = JSON.parse(JSON.stringify(res[0]));
      }
    }
  } else {
    message.error("后端错误");
  }
}; // 初始化config的函数

const initExperimentDimension = async (scene_name: string) => {
  const res = await getExperimentConfigDimension(scene_name);
  if (res) {
    if (typeof res === "string") {
      message.error(res);
    } else {
      dimension_id.value = res[0].DimensionId;
    }
  }
}; // 初始化Dimension函数

// watch && computed

let columns = computed(() => {
  return [
    {
      title: "流程ID",
      key: "workflowid",
      sorter: (
        a: experimentConfigResContentI,
        b: experimentConfigResContentI
      ) => Number(a.workflowid) - Number(b.workflowid),
    },
    {
      title: "场景编号",
      key: "scene_id",
    },
    {
      title: "场景名称",
      key: "scenename",
    },
    {
      title: "策略组合",
      key: "流程图",
      render(row: experimentConfigResContentI) {
        return h("div", null, [
          h(
            "span",
            { onClick: () => showWorkflow(row), style: { cursor: "pointer",color:"green" } },
            "流程图"
          ),
          h(
            "span",
            { style: { color: "red" } },
            { default: () => (row.showmodules ? null : " 未配置") }
          ),
        ]);
      },
    },
    {
      title: "主/从",
      key: "defult",
    },
    {
      title: "分组名称",
      key: "groupname",
    },
    {
      title: "备注",
      key: "remark",
    },
    {
      title: "更新时间",
      key: "updatetime",
      sorter: (
        a: experimentConfigResContentI,
        b: experimentConfigResContentI
      ) => Date.parse(a.updatetime) - Date.parse(b.updatetime),
    },
    {
      title: "操作人",
      key: "operator",
    },
    {
      title: "操作",
      key: "action",
      align: "center",
      fixed: "right",
      render(row: experimentConfigResContentI) {
        return h(NSpace,null,{default:() => [
          renderTooltip('修改',() => onUpdateConfigData(row),24),
          renderTooltip('导出',() => onImportConfigData(row),24),
          renderTooltip('复制',() => onCopyConfigData(row),24),
          renderTooltip('导入',() => onRequireConfigData(row.workflowid),24),
          row.defult === "从" ? renderTooltip('删除',() => onDeleteConfigData(row),24) : null,
        ] });
      },
    },
  ] as any;
});

// hook
onBeforeMount(async () => {
  await initExperimentConfigName(types.value);
  await initExperimentConfigAppname();
  await refresh();
  if (!formValue.value.app_id) {
    formValue.value.app_id = app_idData.value[0].label;
  }
  if (!formValue.value.scene_id) {
    formValue.value.scene_id = scene_idData.value[0].label;
  }
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
    .n-button {
      margin: 5px;
    }
    .select-out {
      display: inline-block;
      vertical-align: middle;
      .select {
        display: flex;
        justify-content: center;
        align-items: center;
        margin: auto 5px;
        .n-input {
          width: 150px;
          margin: auto 5px;
          display: inline-block;
        }
        .n-select {
          width: 150px;
          margin: auto 5px;
          display: inline-block;
        }
        label {
          height: 50px;
          width: auto;
          line-height: 50px;
          text-align: center;
        }
      }
    }
  }
}
:deep(.n-button) {
  margin: auto 5px;
  text-align: left;
}
:deep(.n-dialog) {
  padding: 5px !important;
  max-width: 100vw !important;
  min-width: 100vw !important;
}
</style>
