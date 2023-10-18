<template>
  <div class="out-container">
    <div class="vis-wrapper">
      <div class="vis-content">
        <MTFlowChart
            :nodeData="props.dataNodes"
            :idPreFix="props.prefix"
            :lastId="props.suffix"
            @click-node="onClickNode"
            :editable="true"
            :nodeTypes="nodeType"
            @onDataChange="dataChange"
            class="vis-container"
        />
        <div class="vis-sider">
          <div class="sider-top">
            <slot></slot>
          </div>
          <div class="sider-content">
            <div class="config-item">
              <div class="item-label">类型</div>
              <p>{{ workflowConfig.type }}</p>
            </div>

            <div class="config-item">
              <div class="item-label">id</div>
              <p>{{ workflowConfig.id }}</p>
            </div>

            <div class="config-item" v-if="workflowConfig.type !== 'condition'">
              <div class="item-label">名称</div>
              <n-input
                size="small"
                  :on-blur="() => labelChange(workflowConfig.label)"
                  v-model:value="workflowConfig.label"
              />
            </div>

            <div 
              class="config-item"
              v-if="workflowConfig.type == 'condition'"
              >
              <div class="item-label">条件</div>
                <n-select
                  size="small"
                  v-model:value="workflowConfig.conditionSelectValue" 
                  :options="conditionSelectData"
                  clearable
                  @update:value="conditionChange"
                  />
            </div>

            <div class="config-item">
              <div class="item-label">超时</div>
              <n-popover trigger="hover" placement="bottom">
                <template #trigger>
                  <n-input-number
                    size="small"
                      :show-button="false"
                      :on-blur="() => timeoutChange(workflowConfig.timeout)"
                      v-model:value="workflowConfig.timeout"
                      style="width:65%; height:50%"
                  >
                    <template #suffix>
                      ms
                    </template>
                  </n-input-number>
                </template>
                <span>0 表示不限制时间</span>
              </n-popover>
            </div>

            <div class="config-item" v-if="workflowConfig.type == 'flow'">
              <div class="item-label">子流程id</div>
              <!-- <n-popover trigger="hover" placement="bottom">
                <template #trigger> -->
                  <n-input-number
                    size="small"
                      :show-button="false"
                      :on-blur="() => refidChange(workflowConfig.ref_workflow_id)"
                      v-model:value="workflowConfig.ref_workflow_id"
                      style="width:65%; height:50%"
                  >
                  </n-input-number>
                <!-- </template>
                <span>0 表示不限制时间</span>
              </n-popover> -->
            </div>

            <!-- 新增属性：超时异步回调 -->
            <div class="config-item" v-if="workflowConfig.timeout !== 0 && workflowConfig.timeout !== undefined">
              <div class="item-label">异步回调</div>
              <n-switch
                @update:value="timeoutAsyncChange"
                v-model:value="workflowConfig.timeout_async"
                style="margin-top: 4px;"
              />
            </div>

            <!-- 新增属性：是否动态超时 -->
            <div class="config-item" v-if="workflowConfig.timeout !== 0 && workflowConfig.timeout !== undefined">
              <div class="item-label">动态超时</div>
              <n-switch 
                @update:value="timeoutDynamicChange"
                v-model:value="workflowConfig.timeout_dynamic"
                style="margin-top: 4px;"
              />
            </div>

          </div>

          <div class="sider-bottom">
            <div class="config" v-for="(items, index_) in configArr">
              <div class="set-config">
                <div class="config-item" v-for="(item, index) in items">
                  <label for="name">{{ index }}</label>
                  <n-input
                    size="small"
                      id="name"
                      v-model:value="configArr[index_][index]"
                      v-if="index !== 'type'"
                  />
                  <n-select
                    size="small"
                      :options="typeOptions"
                      clearable
                      v-model:value="configArr[index_][index]"
                      v-if="index === 'type'"
                  />
                </div>
              </div>
              <div class="remove-config" @click="removeConfig(index_)">
                <n-icon>
                  <Close/>
                </n-icon>
              </div>
            </div>
            <n-button
                @click="addConfig"
                type="success"
                v-if="workflowConfig.id"
            >
              添加参数
            </n-button>
            <n-button
                @click="saveConfig"
                type="warning"
                v-if="workflowConfig.id"
            >
              保存参数
            </n-button>
          </div>
        </div>
      </div>
    </div>
    <div class="vis-button">
      <n-button type="success" class="success" @click="callbackSave">
        保存
      </n-button>
      <n-button
          type="warning"
          class="reset"
          @click="() => emits('resetWorkflow')"
      >取消
      </n-button
      >
    </div>
  </div>
</template>

<script lang="ts" setup>
import { defineProps, ref, PropType, defineEmits, onBeforeMount } from "vue";
import {
  experimentWorkflowNodeTypeResContentI,
  experimentConfigResContentI,
} from "@/api/strategy";
import MTFlowChart from "@/components/g6/flow.vue";
import {Close} from "@vicons/ionicons5";
import {useMessage} from "naive-ui";
import {typeOptions} from "@/utils";

const message = useMessage();

const props = defineProps({
  dataNodes: {
    type: Object as PropType<any>,
    default: () => ({nodes: [], edges: []}),
  },
  nodeType: {
    type: Object as PropType<Array<experimentWorkflowNodeTypeResContentI>>,
    default: () => [{node_name: [], node_type: ""}],
  },
  workflowId: {
    type: String,
    default: "",
  },
  prefix: {
    type: String,
    default: "",
  },
  suffix: {
    type: Number,
    default: 0,
  },
  operator: {
    type: String,
    default: "",
  },
  experimentData: {
    type: Object as PropType<experimentConfigResContentI>,
    default: () => {
    },
  },
});

const emits = defineEmits(["resetWorkflow", "updateSuccess", "addSuccess"]);

interface FormDataI {
  label: string;
  id: string;
  type: string;
  shape: string;
  timeout: any;
  timeout_async: boolean;
  timeout_dynamic: boolean;
  ref_workflow_id?: any;
  conditionSelectValue?: any;
}
const workflowConfig = ref({} as FormDataI);

const conditionSelectData = [
  {
    label: '大于',
    value: 'GT',
  },
  {
    label: '小于',
    value: 'LT',
  },
  {
    label: '等于',
    value: 'EQ',
  },
  {
    label: '大于等于',
    value: 'GE',
  },
  {
    label: '小于等于',
    value: 'LE',
  },
  {
    label: '不等于',
    value: 'NE',
  },
  {
    label: '枚举',
    value: 'ENUM',
  }
]

// let conditionSelectValue = ref('条件节点' as string)

const labelChange = (text: string) => {
  const item = tempNodes()
  const model = item.getModel()
  model.label = text
  model.data.label = text
  item.update(model)
};

const timeoutChange = (timeout: number) => {
  const item = tempNodes()
  const model = item.getModel()
  model.data.timeout = timeout
  item.update(model)
  if (timeout == 0){
    timeoutAsyncChange(false)
    timeoutDynamicChange(false)
  }
}

const timeoutAsyncChange = (timeout_async: boolean) => {
  const item = tempNodes()
  const model = item.getModel()
  model.data.timeout_async = timeout_async
  item.update(model)
}

const timeoutDynamicChange = (timeout_dynamic: boolean) => {
  const item = tempNodes()
  const model = item.getModel()
  model.data.timeout_dynamic = timeout_dynamic
  item.update(model)
}

const refidChange = (ref_workflow_id: number) => {
  const item = tempNodes()
  const model = item.getModel()
  model.data.ref_workflow_id = ref_workflow_id
  // model.label = ref_workflow_id
  item.update(model)
}

const configArr = ref([] as { name: string; value: string; type: string }[]);

let tempNodes: () => any;
let dataNodesJson: any;
let curClickNode: any
// let ref_workflow_id:any

const onClickNode = (node: any) => {
  const res = node.getModel()
  if(res){
    curClickNode = res
    configArr.value = [];
    tempNodes = () => node;
    console.log(res);
    // let tempLen = res.attrs.text.text.length;
    // workflowConfig.value.type = res.data?.nodeType ? res.data?.nodeType : (tempLen > 3 ? 'task' : 'condition');
    // workflowConfig.value.type = res.data?.type;
    workflowConfig.value.type = res.data?.type
    workflowConfig.value.shape = res.shape
    workflowConfig.value.id = res.id;
    // workflowConfig.value.ref_workflow_id = res.data?.data?.ref_workflow_id ? res.data.data?.ref_workflow_id : 0;
    // ref_workflow_id = workflowConfig.value.ref_workflow_id
    // workflowConfig.value.type == 'condition' ? conditionSelectValue.value = res.data.attrs.text.text : workflowConfig.value.label = res.attrs.text.text;
    workflowConfig.value.label = res.data.label
    workflowConfig.value.conditionSelectValue = res.data.conditionSelectValue
    workflowConfig.value.timeout = res.data?.conditionSelectValue ? res.data?.conditionSelectValue : 0
    workflowConfig.value.timeout = res.data?.timeout ? res.data?.timeout : 0
    workflowConfig.value.timeout_async = res.data?.timeout_async ? res.data?.timeout_async : false
    workflowConfig.value.timeout_dynamic = res.data?.timeout_dynamic ? res.data?.timeout_dynamic : false
    workflowConfig.value.ref_workflow_id = res.data?.ref_workflow_id ? res.data?.ref_workflow_id : 0
    
    if (res.data?.params) {
      res.data?.params.forEach((p: any) => {
        configArr.value.push(p);
      });
    }
  }
};

const conditionChange = (value:string) => {
  const item = tempNodes()
  const model = item.getModel()
  model.data.conditionSelectValue = value
  model.label = value
  item.update(model)
}
// watch(conditionSelectValue, (newVal) => {
//   curClickNode.setAttrs({
//     label: {text: newVal}
//   })
// })

// const callbackSave = () => {
//   dataNodesJson.nodes.forEach((p: any) => {
//     p.params = p.data?.params;
//     p.type = p.data.type;
//     p.shape = enmuType[p.data.type];
    
//     if(p.type == 'condition') {
//       p["label"] = p.data.label;
//     }else {
//       p["label"] = p.label;
//     }
//     p.conditionSelectValue = p.data?.conditionSelectValue
//     p.ref_workflow_id = p.data?.ref_workflow_id
//     p.timeout = p.data?.timeout
//     p.timeout_async = p.data?.timeout_async == undefined ? false : p.data?.timeout_async
//     p.timeout_dynamic = p.data?.timeout_dynamic == undefined ? false : p.data?.timeout_dynamic
//     return p
//   });
//   console.log('dataNodesJson:');
//   console.log(dataNodesJson);
//   dataNodesJson.edges.forEach((p: any) => {
//     p.source = p.source;
//     p.target = p.target;
//     p["label"] = p.labels ? p.labels[0]?.attrs.label.text : null;
//     return p
//   });
//   if (props.operator === "添加") {
//     emits("addSuccess", JSON.stringify(dataNodesJson));
//   } else if (props.operator === "修改") {
//     emits("updateSuccess", JSON.stringify(dataNodesJson));
//   }
// };

const callbackSave = () => {
  if (dataNodesJson === undefined) {
    dataNodesJson = props.dataNodes
  } 
  // console.log(dataNodesJson) 
  // const actions: any = {}
  const flow_nodes: any[] = []
  const save_nodes: any[] = []
  dataNodesJson.nodes.forEach((p: any) => {
    if (!p.hasOwnProperty('comboId') || p.comboId == undefined){
      p.params = p.data?.params;
      switch(p.data.type) {
        case 'task': 
          p.type = 'rect';
          break;
        case 'timeout': 
          p.type = 'clock';
          break;
        case 'condition': 
          p.type = 'diamond';
          break;
        case 'flow': 
          p.type = 'circle';
          break;
      }
      if(p.data.type == 'condition') {
        p["label"] = p.data.conditionSelectValue;
      }else {
        p["label"] = p.data.label;
      }
      p.ref_workflow_id = p.data?.ref_workflow_id;
      p.timeout = p.data?.timeout;
      p.timeout_async = p.data?.timeout_async == undefined ? false : p.data?.timeout_async;
      p.timeout_dynamic = p.data?.timeout_dynamic == undefined ? false : p.data?.timeout_dynamic;
      p.location = p.x + "," + p.y;
      save_nodes.push(p);
    }else{
      flow_nodes.push(p);
    }
  });
  console.log('save:',save_nodes, 'flow:', flow_nodes);
  dataNodesJson.nodes = save_nodes;

  const flow_edges = [];
  const save_edges: any[] = [];
  dataNodesJson.edges.forEach((p: any) => {
    if (!p.hasOwnProperty('comboId') || p.comboId === undefined){
      // p.source = p.source;
      // p.target = p.target;
      p["label"] = p.labels ? p.labels[0]?.attrs.label.text : null;
      save_edges.push(p);
    }else{
      flow_edges.push(p);
    }
  });
  dataNodesJson.edges = save_edges;
  
  // for (let k in dataNodesJson.nodes) {
  //   const node = dataNodesJson.nodes[k]
  //   const action: any = {}
  //   action.action_type = node.data.type;
  //   action.params = node.data?.params;
  //   action.action_id = node.id;
  //   action.action_name = node.label;
  //   action.ref_workflow_id = node.data?.ref_workflow_id;
  //   action.timeout = node.data?.timeout;
  //   action.timeout_async = node.data?.timeout_async == undefined ? false : node.data?.timeout_async;
  //   action.timeout_dynamic = node.data?.timeout_dynamic == undefined ? false : node.data?.timeout_dynamic;
  //   action.next_action_ids = []
  //   action.next_conditions = []
  //   for (let j in dataNodesJson.edges) {
  //     const edge = dataNodesJson.edges[j]
  //     if (edge.source === node.id) {
  //       action.next_action_ids.push(edge.target)
  //       if (action.action_type === 'condition') {
  //         action.next_conditions.push(edge.label)
  //       }
  //     }
  //   }    
  //   actions[action.action_id] = action
  // }
    // console.log('actions:');
  // console.log(actions);

  console.log('dataNodesJson:');
  console.log(dataNodesJson);

  if (props.operator === "添加") {
    emits("addSuccess", JSON.stringify(dataNodesJson));
  } else if (props.operator === "修改") {
    emits("updateSuccess", JSON.stringify(dataNodesJson));
    // emits("updateSuccess", JSON.stringify(actions));
  }
};

const dataChange = (data: any) => {
  data.edges?.forEach((p: any) => {
    if (p.labels?.length > 1) {
      message.error("边的参数不能大于一个哦~");
    }
  });

  dataNodesJson = data;
};

const addConfig = () => {
  message.warning("添加完参数记得点保存哦");
  configArr.value.push({
    name: "",
    value: "",
    type: "int",
  });
};

const saveConfig = () => {
  configArr.value.forEach((p) => {
    if (!p.name || !p.type || !p.value) {
      message.error("参数信息不能为空");
      return;
    }
  });
  const item = tempNodes()
  const model = item.getModel()
  model.data.params = configArr.value
  item.update(model)
};

const removeConfig = (index: number) => {
  configArr.value.splice(index, 1);
};

onBeforeMount(() => {
  console.log(props.dataNodes,'-----');
});
</script>

<style lang="less" scoped>
.out-container {
  width: 100%;
  height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;

  .vis-wrapper {
    width: 100%;
    height: 100%;
    overflow-y: hidden;
    display: flex;
    box-sizing: border-box;
    margin-right: 5px;

    .vis-content {
      flex: 1;
      display: flex;
      width: 100%;
      height: 100%;
      padding: 5px 10px;
      box-sizing: border-box;

      .vis-container {
        border: 1px solid gray;
        border-right: none;
        flex: 1;
      }

      .vis-sider {
        width: 210px;
        height: 100%;
        border: 1px solid gray;
        display: flex;
        flex-direction: column;
        font-size: 12px;

        .sider-top {
          width: 100%;
          padding: 5px;
          height: 160px;
          box-sizing: border-box;
          overflow-y: scroll;
          border-bottom: 1px solid gray;
          margin-bottom: 5px;
        }

        .sider-content {
          margin-top: 3px;
          margin-bottom: 3px;
          padding: 5px;
          border-top: 1px solid gray;
          width: 100%;
          border-bottom: 1px solid gray;

          .config-item {
            margin-top: 3px;
            margin-bottom: 3px;
            background: rgba(0, 0, 0, 0.1);
            display: flex;
            line-height: 30px;
            padding: 0 5px;
            box-sizing: border-box;

            .item-label {
              display: inline-block;
              width: 60px;
              text-align: start;
              margin: 0 5px;
              box-sizing: border-box;
              line-height: 30px;
            }

            p {
              flex: 1;
              margin: 0;
              color: #09a4f4;
              text-align: left;
            }

            :deep(.n-input) {
              flex: 1;
              height: 30px;
              width: 170px;
            }
            :deep(.n-select) {
              flex: 1;
              max-height: 30px !important;
              width: 170px;
            }
          }
        }

        .sider-bottom {
          text-align: center;
          border-top: 1px solid gray;
          width: 100%;
          flex: 1;
          margin: 7px auto;
          padding-top: 10px;
          overflow-y: scroll;

          .config {
            flex: 1;
            margin: 10px auto;
            display: flex;
            width: 100%;
            align-items: center;

            .config-item {
              flex: 1;
              display: flex;
              padding: 0 5px;
              box-sizing: border-box;

              label {
                width: 60px;
                text-align: right;
                padding: 0 5px;
                box-sizing: border-box;
              }

              :deep(.n-input) {
                height: 30px;
                margin: 2px auto;
              }
            }

            .remove-config {
              width: 30px;
            }
          }

          :deep(.n-button) {
            margin: 0 5px;
          }
        }
      }
    }
  }

  .vis-button {
    height: 0;
    text-align: right;

    :deep(.n-button) {
      margin: 5px 10px;
    }
  }
}

::-webkit-scrollbar {
  width: 0;
}
</style>
