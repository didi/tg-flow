<template>
  <n-data-table
    :columns="columns"
    :data="data"
    :pagination="pagination"
    :bordered="false"
  />
</template>

<script lang='ts' setup>
import { defineProps, onBeforeMount,defineEmits, ref,h } from 'vue'
import {getScenesHistoryVersion,experimentHistoryVersionResFilterI,backToHistoryVersion} from '@/api/strategy';
import { useDialog, useMessage } from 'naive-ui';

const props = defineProps({
  scenename:{
    type:String,
    default:''
  },
  dimension_id:{
    type:Number,
    default:0
  }
})

const emits = defineEmits(['backHistory'])

const message = useMessage()

const dialog = useDialog();

const columns = [
  {
    title:"版本号",
    key:"version_id",
  },
  {
    title:"创建时间",
    key:"version_create_time",
  },
  {
    title:"操作",
    key:"operator",
    render:(row:experimentHistoryVersionResFilterI) =>{
     return h('p',{
        onClick:() => toBackVersion(row),
        style:{
          cursor:"pointer",
          margin:0
        }
      },
      {
        default:() => '回滚至该版本'
      })
    }
  }
]

const pagination = { pageSize: 15 };

const data = ref([] as Array<experimentHistoryVersionResFilterI>)

const toBackVersion =async (row:experimentHistoryVersionResFilterI) => {
  dialog.info({
    title:"提示",
    showIcon: false,
    content:"确认要回滚版本吗？",
    positiveText:"确认",
    onPositiveClick: async() => {
      const res = await backToHistoryVersion({
        operator:window.userInfo.username,
        scene_name:props.scenename,
        dimension_id:props.dimension_id,
        version_id:row.version_id
      })
      if(res === true){
        emits('backHistory',res)
      }else{
        message.error(`回滚失败,${res}`)
      }
    },
    negativeText:"取消"
   })
  
}

const initVersion = async () => {
  const res = await getScenesHistoryVersion({
      operator: window.userInfo.username,
      scenename: props.scenename,
      dimension_id: props.dimension_id,
    });
  if (res) {
    if (typeof res === "string") {
      message.error(res);
    } else {
      data.value = res
    }
  }
}

onBeforeMount(async() => {
  await initVersion();
})
</script>

<style lang='less' scoped>

</style>