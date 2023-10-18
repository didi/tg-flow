<template>
  <n-form
    ref="formRef"
    :label-width="80"
    :model="formValue"
    :rules="rules"
  >
    <n-form-item label="系统编号" path="id" v-if="addOrrevise">
      <n-input-number v-model:value="formValue.id" disabled placeholder="输入id" />
    </n-form-item>
    <n-form-item label="系统编号" path="id" v-else>
      <n-input-number v-model:value="formValue.id" placeholder="输入id" />
    </n-form-item>
    <n-form-item label="系统名称" path="app_name">
      <n-input v-model:value="formValue.app_name" placeholder="输入系统名称" />
    </n-form-item>
    <n-form-item label="部署机房" path="machine_room">
      <n-input
        v-model:value="formValue.machine_room"
        placeholder="输入机房名称"
      />
    </n-form-item>
    <n-form-item label="节点名称" path="node_name">
      <n-input v-model:value="formValue.node_name" placeholder="输入节点名称" />
    </n-form-item>
    <n-form-item label="Git 仓库" path="git_url">
        <n-input v-model:value="formValue.git_url" placeholder="请输入仓库地址" />
    </n-form-item>
    <n-form-item>
      <n-button type="info" class="successButton" @click="handleValidateClick">
        保存
      </n-button>
    </n-form-item>
  </n-form>
</template>

<script lang='ts' setup>
import { defineProps, onMounted, PropType, ref } from 'vue'
import { addSystemConfig,updateSystemConfig } from '@/api/strategy';
import { FormInst, FormItemRule, useMessage } from 'naive-ui'

// interface 
interface SectionDataI {
  app_name: string;
  id: number;
  machine_room: string;
  node_name: string;
  git_url: string;
}

// state
const props = defineProps({
  systemConfigData: {
    type:Object as PropType<SectionDataI>,
    default:() => {}
  }
})

let addOrrevise = ref(false); // 是修改还是添加

const formRef = ref<FormInst | null>()

const message = useMessage()

const formValue = ref({} as SectionDataI)

const rules =  {
      machine_room:{
        required: true,
          trigger: 'blur',
          validator: (rule: FormItemRule, value: string) => {
            return new Promise<void>((resolve, reject) => {
              if (value.replaceAll(" ","").charAt(0) === ',' || value.replaceAll(" ","").charAt(value.replaceAll(" ","").length - 1) === ',' || !value) {
                reject(Error('首位不能是逗号')) // reject with error message
              } else {
                resolve()
              }
            })
          }
      },
      node_name:{
        required: true,
        trigger: 'blur',
        message: '请输入节点名称',
      },
      app_name:{
        required: true,
        trigger: 'blur',
        message:"请输入系统名称"
      },
      id:{
        required: true,
        trigger: 'blur',
        validator: (rule: FormItemRule, value: string) => {
          return new Promise<void>((resolve, reject) => {
              if (!value || Number(value) < 0) {
                reject(Error('请输入正确的ID')) // reject with error message
              } else {
                resolve()
              }
          })
        }
      },
  }

const emits = defineEmits(["updateSuccess","addSuccess"])
// methods

const handleValidateClick =  (e: MouseEvent) => {
    e.preventDefault()
    const messageReactive = message.loading('Verifying', {
      duration: 0
    })
    formRef.value?.validate(async (errors) => {
      if (!errors) {
        if(addOrrevise.value){
          let res = await updateSystemConfig({...formValue.value,operator:window.userInfo.username,old_id:formValue.value.id})
          if(res){
            message.success('修改成功')
            emits('updateSuccess',true)
          }else{
            message.error("修改失败")
          }
        }else{
          let res = await addSystemConfig({...formValue.value,operator:window.userInfo.username})
          if(typeof res === "object"){
            message.success('添加成功')
            emits('addSuccess',true)
          }else{
            message.error(res? res : "null")
          }
        }
      } else {
        message.error("更新失败,请重新填写")
        console.log('errors', errors)
      }
      messageReactive.destroy()
    })
  }

//watch && computed

// hook

onMounted(() => {
  addOrrevise.value = Object.keys(props.systemConfigData).length !== 0
  if(addOrrevise.value){
    formValue.value = props.systemConfigData
  }else{
    formValue.value = {
      id:-1,
      app_name:"",
      node_name:"",
      machine_room:"",
      git_url: ""
    }
  }
  

})
</script>

<style lang='less' scoped>
.n-form{
  .n-form-item{
    // margin: 0;
    :deep(.n-form-item-blank){
      display: block !important;
      width: 100%;
      :deep(.n-input-number) {
      width: 100%;
        :deep(.n-button){
          width: 100%;
        }
      }
      .successButton{
          width: 30%;
        }
    }
    .n-input{
      width: 100%;
    }
    
  }
}

</style>