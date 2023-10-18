<template>
  <n-form label-align="right" label-placement="left" label-width="120px">
    <n-form-item label="实验编号">
      <n-input :disabled='true' :value="id"/>
    </n-form-item>
    <n-form-item label="导入内容(Json)">
      <n-input type="textarea" v-model:value="json" placeholder="请输入json格式"/>
    </n-form-item>
    <div class="button">
      <n-button type="primary" @click="onClickCallBack">确认</n-button>
      <n-button @click="onClickReset">取消</n-button>
    </div>
  </n-form>
</template>

<script lang='ts' setup>
import { useMessage } from 'naive-ui';
import { defineProps, ref } from 'vue'
const props = defineProps({
  id:{
    type:String,
    default:''
  }
})

const message = useMessage()

const emits = defineEmits(['success','reset'])

const onClickCallBack = () => {
  if(json.value){
    emits('success',json.value,Number(props.id))
  }else{
    message.warning('请输入一点内容吧~')
  }
}

const onClickReset = () => emits('reset')

const json = ref('' as string)

</script>

<style lang='less' scoped>
.button{
  float: right;
  .n-button{
    margin:0 10px;
  }
}
</style>