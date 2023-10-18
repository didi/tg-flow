<template>
  <div class="x6-container">
    <div class="x6-content">
      <MTFlowChart 
        :nodeData="props.dataNodes"
        @click-node="onClickX6"
      />
    </div>
    <div class="x6-params" v-if="params">
      <p> <span class="left">ID:</span> <span class="right">{{params.id}}</span></p>
      <p> <span class="left">名称:</span> <span class="right">{{params.label}}</span></p>
      <p> <span class="left">类型:</span> <span class="right">{{params.type}}</span></p>
      <p> <span class="left">超时:</span> <span class="right">{{params.timeout}} ms</span></p>

      <div class="params-name" v-if="params.params.length > 0">params:</div>
      <div v-if="params.params.length > 0" style="max-height:300px;overflow-y: hidden;">
        <n-data-table
          style="box-sizing: border-box;padding: 5px;"
          max-height="200px"
          :columns="paramsList"
          :data="params.params"
          :bordered="true"
          scroll-x="200"
        />
      </div>
    </div>
  </div>
</template>

<script lang='ts' setup>
import { defineProps, ref,PropType } from 'vue';
import MTFlowChart from "@/components/g6/flow.vue";

const props = defineProps({
  dataNodes: {
    type: Object as PropType<any>,
    default: () => ({ nodes: [], edges: [] }),
  },
})

const params = ref(null as any)

const paramsList = [
  {
    title:"名称",
    key:"name",
    width: 100,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title:"值",
    key:"value",
    width: 100,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title:"类型",
    fixed: "right",
    width:"70",
    key:"type"
  },
]

const onClickX6 = (res: any) => {
  if(res){
    let tempLen = res.attrs.text.text.length;
    params.value = {
      type: res.store.data?.nodeType ? res.store.data?.nodeType : (tempLen > 3 ? 'task' : 'condition'),
      label: res.attrs.text.text,
      id: res.id,
      timeout: res.store.data?.timeout,
      params: res.store.data?.data?.params ? res.store.data?.data?.params : []
    }
    // console.log(params.value);
    
  }else params.value = null
}
</script>

<style lang='less' scoped>
.x6-container{
  width: 100%;
  height: 100%;
  position: relative;
  .x6-content{
    width: 100%;
    height: 100%;
  }
  .x6-params{
    position: absolute;
    right: 10px;
    top: 10px;
    width: 300px;
    height: auto;
    border: 1px solid gray;
    padding: 5px;
    background: white;
    box-shadow: 0 0 10px rgba(0,0,0,.2);
    border-radius: 5px;
    p{
      margin: 2px auto;
      display: flex;
      height: 24px;
      white-space: nowrap;
      .left{
        width: 60px;
        text-align: right;
        color: green;
      }
      .right{
        padding-left: 5px;
        flex: 1;
        overflow-x: scroll;
      }
    }
    .params-name{
      display: inline-block;
      width: 60px;
      text-align: right;
      color:green
    }
  }
}
::-webkit-scrollbar {
  width: 0;
}
</style>