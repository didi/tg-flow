<template>
  <div class="update-form">
    <div class="item-form"> 
      <div class = "item-label">实验编号</div>
      <n-input size="small" v-model:value="data.workflowid" placeholder="可不填，自动获取"/>
    </div>
    <div class="item-form"> 
      <div class = "item-label">分组名称</div>
      <n-input size="small" type="textarea" autosize v-model:value="data.groupname"/>
    </div>
    <div class="item-form"> 
      <div class = "item-label">主从</div>
      <n-input size="small" v-model:value="data.defult" disabled/>
    </div>
    <div class="item-form"> 
      <div class = "item-label">备注</div>
      <n-input size="small" type="textarea" autosize v-model:value="data.remark"/>
    </div>
  </div>
</template>

<script lang='ts' setup>
import { defineProps, PropType} from 'vue';
import { experimentConfigResContentI } from '@/api/strategy';

const props = defineProps({
  message:{
    type:Object as PropType<experimentConfigResContentI>,
    required: true,
  }
})

const data = props.message as experimentConfigResContentI

</script>

<style lang='less' scoped>
.update-form{
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow-y: scroll;
  .item-form {
    background: rgba(0, 0, 0, 0.1);
    margin-top: 3px;
    margin-bottom: 3px;
    display: flex;
    line-height: 30px;
    padding: 0 5px;
    box-sizing: border-box;

    .item-label {
      display: inline-block;
      width: 80px;
      text-align: start;
      margin: 0 5px;
      box-sizing: border-box;
      line-height: 30px;
    }
  }
}
</style>