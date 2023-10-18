<template>
  <div class="container" id="container">
    <div class="side-tools-wrapper" v-if="editable">
      <div class="node-type">
        <div class="title">流程</div>
        <div class="nodes">
          <div
            class="node"
            title="缺省节点"
            draggable="true"
            @dragend="addItem($event, 'task', '任务节点')"
          >
            <div class="id">ID：action-id</div>
            <div class="label">缺省节点</div>
          </div>
          <div
            class="node"
            title="超时节点"
            draggable="true"
            @dragend="addItem($event, 'timeout', '超时任务')"
          >
            <div class="id">ID：action-id</div>
            <div class="label">超时节点</div>
          </div>
          <div
            class="node diamond"
            title="条件节点"
            draggable="true"
            @dragend="addItem($event, 'condition', '条件')"
          >
            <div class="label">条件节点</div>
          </div>
          <div
            class="node circle"
            title="子流程"
            draggable="true"
            @dragend="addItem($event, 'flow', '子流程')"
          >
            <div class="label">子流程</div>
          </div>
        </div>
      </div>
      <div class="node-type" v-for="(type, index) in nodeTypes" :key="index">
        <div class="title">{{ type.node_type }}</div>
        <div class="nodes">
          <div
            class="node"
            v-for="name in type.node_name"
            :title="`${type.node_type}.${name}`"
            draggable="true"
            @dragend="addItem($event, 'task', `${type.node_type}.${name}`)"
          >
            <div class="id">ID：action-id</div>
            <div class="label">{{ name }}</div>
          </div>
        </div>
      </div>
    </div>
    <div id="g6-container" class="g6-container"></div>
    <div id="g6-toolbar" class="g6-toolbar"></div>
    <div id="g6-minimap" class="g6-minimap"></div>
  </div>
</template>
<script lang="ts">
import {defineComponent, onBeforeUnmount, onMounted} from "vue";
import {NCollapse, NCollapseItem} from "naive-ui";
import {renderFlow} from "./g6";

export default defineComponent({
  props: {
    nodeData: Object,
    editable: Boolean,
    idPrefix: String,
    lastId: Number,
    nodeTypes: Array,
  },
  components: {
    NCollapse,
    NCollapseItem
  },
  setup(props, { emit }) {
    let graph: any;
    let idCount = props.lastId || 1;
    onMounted(() => {
      graph = renderFlow(
        "g6-container",
        props.nodeData!.nodes!,
        props.nodeData!.edges!,
        props.nodeData!.combos!,
        props.editable
          ? {
              click: (data: any) => {
                emit("click-node", data);
              },
              datachange: (data: any) => {
                emit("onDataChange", data);
              }
            }
          : {},
        props.editable
      );
    });
    const getModelTpl = (type: string, name: string) => {
      const id = (props.idPrefix || "action-") + props.nodeData?.id + "-" + ++idCount;
        return {
          id: id,
          title: id,
          x:0,
          y:0,
          desc: name,
          data: {
              id: id,
              label: name,
              type: type,
              params: null,
              timeout: 0,
              ref_workflow_id: 0,
              timeout_async: false,
              timeout_dynamic: false,
          },
      };
    };
    const addItem = (e: any, type: string, name: string) => {
      // console.log(type);
      const point = graph.getPointByClient(e.x, e.y);
      const model = getModelTpl(type, name);
      model.x = point.x;
      model.y = point.y;
      graph.addItem("node", model, true);
    };
    onBeforeUnmount(() => {
      graph.destroy();
    })
    return {
      addItem,
    };
  },
});
</script>
<style scoped>
.container {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  flex-direction: row;
}
.side-tools-wrapper {
  width: 170px;
  height: 100%;
  border: 1px solid #dbdbdb;
  overflow: scroll;
}

.node-type {
  display: flex;
  flex-direction: column;
  /* text-align: center;
  margin-left: auto;
  margin-right: auto; */
}

.node-type .title {
  font-size: 14px;
  font-weight: 500;
  background: #dbdbdb;
}
.node-type .nodes {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  flex-wrap: wrap;
  /* text-align: center;
  margin-left: auto;
  margin-right: auto; */
}

.node-type .nodes .node {
  box-sizing: border-box;
  width: 120px;
  height: 44px;
  border: 3px solid #5B8FF9;
    user-select: none;
  cursor: move;
  /* text-align: center; */
    margin: 3px auto;
}

.node-type .nodes .diamond {
  width: 100px;
  height: 45px;
  background-color: #5B8FF9;
  display: flex;
  align-content: center;
  justify-content: center;
  align-items: center;
  clip-path: polygon(0px 50%,50% 0px,100% 50%,50% 100%,0px 50%);
  /* text-align: center; */
  margin: 3px auto;
}

.node-type .nodes .circle {
  width: 52px;
  height: 52px;
  background-color: #5B8FF9;
  border-radius: 100%;
  display: flex;
  /* align-content: center; */
  justify-content: center;
  align-items: center;
  text-align: center;
  margin: 3px auto;
}

.node-type .nodes .node .id {
  background: #5B8FF9;
  font-size: 10px;
  font-weight: 500;
  color: #fff;
  height: 20px;
  overflow: hidden;
  padding: 0 5px;
}

.node-type .nodes .node .label {
  font-size: 10px;
  color: rgba(0, 0, 0, 0.4);
  padding: 0 5px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.node-type .nodes .diamond .label, .node-type .nodes .circle .label {
  font-size: 5px;
  color: #fff;
}

.g6-container {
  width: 100%;
  height: 100%;
  border: 1px solid #dbdbdb;
}
.g6-minimap {
  position: absolute;
  bottom: 0;
  right: 0;
}
.g6-toolbar {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 250px;
}
</style>
