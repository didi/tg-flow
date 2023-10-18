<template>
  <n-form ref="formRef" :label-width="80" :model="formValue" :rules="rules">
    <n-form-item label="场景编号" path="id" v-if="addOrrevise">
      <n-input-number
        v-model:value="formValue.id"
        disabled
        placeholder="输入场景编号"
      />
    </n-form-item>
    <n-form-item label="场景编号" path="id" v-else>
      <n-input-number v-model:value="formValue.id" placeholder="输入场景编号" />
    </n-form-item>
    <n-form-item label="场景名称" path="name">
      <n-input v-model:value="formValue.name" placeholder="输入场景名称" />
    </n-form-item>

    <n-form-item label="系统名称" path="appname">
      <n-select
        v-model:value="formValue.appname"
        filterable
        clearable
        :options="selectData"
      />
    </n-form-item>
    <n-form-item label="实验名称" path="expname">
      <n-input v-model:value="formValue.expname" placeholder="输入实验名称" />
    </n-form-item>
    <n-form-item label="分流类型" path="flow_type">
      <n-select
        v-model:value="formValue.flow_type"
        clearable
        placeholder="输入分流类型"
        :options="flow_typeData"
      />
    </n-form-item>
    <n-form-item>
      <n-button
        type="info"
        class="successButton"
        @click="handleValidateClick"
      >
        完成
      </n-button>
    </n-form-item>
  </n-form>
</template>

<script lang="ts" setup>
import { defineProps, onMounted, PropType, ref } from "vue";
import {
  getScenesConfigAppname,
  addScenesConfig,
  updateScenesConfig,
} from "@/api/strategy";
import { FormInst, useMessage, FormItemRule } from "naive-ui";
import { valueAndLabelI } from "@/types/globalInterface";

// interface
interface SectionDataI {
  appname: string;
  id: number;
  name: string;
  expname: string;
  buckettype: number;
  flow_type: number;
}

// state
const props = defineProps({
  scenesConfigData: {
    type: Object as PropType<SectionDataI>,
    default: () => {},
  },
});

let addOrrevise = ref(false); // 是修改还是添加

const formRef = ref<FormInst | null>();

const message = useMessage();

const selectData = ref(
  [] as Array<{
    label: string;
    value: string;
  }>
);

const flow_typeData = [
  {
    value:'0',
    label:"自研分流"
  },
  {
    value:'3',
    label:"阿波罗"
  }
] as any[] ;

const formValue = ref({} as SectionDataI);

const rules = {
  name: {
    required: true,
    trigger: "blur",
    message: "请输入场景名称",
  },
  appname: {
    required: true,
    trigger: "blur",
    message: "请选择系统名称",
  },
  id: {
    required: true,
    trigger: "blur",
    validator: (rule: FormItemRule, value: string) => {
      return new Promise<void>((resolve, reject) => {
        if (!value || Number(value) < 0) {
          reject(Error("请输入正确的ID")); // reject with error message
        } else {
          resolve();
        }
      });
    },
  },
  expname: {
    required: true,
    trigger: "blur",
    message: "请输入实验名称",
  },
  flow_type: {
    required: true,
    trigger: "blur",
    validator: (rule: FormItemRule, value: string) => {
      return new Promise<void>((resolve, reject) => {
        if (value == "0" || Number(value) === 3) {
          resolve();
        } else {
          reject(Error("请输入正确的分流类型")); // reject with error message
        }
      });
    },
  },
};

const emits = defineEmits(["updateSuccess", "addSuccess"]);
// methods

const getSelects = async () => {
  let res = await getScenesConfigAppname({
    operator: window.userInfo.username,
  });
  if (typeof res === "string") {
    message.error(res);
    return null;
  } else {
    selectData.value = res as valueAndLabelI[];
  }
};

const handleValidateClick = (e: MouseEvent) => {
  e.preventDefault();
  const messageReactive = message.loading("Verifying", {
    duration: 0,
  });
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      if (addOrrevise.value) {
        let res = await updateScenesConfig({
          ...formValue.value,
          operator: window.userInfo.username,
          oldid: formValue.value.id,
        });
        if (res === true) {
          message.success("修改成功");
          emits("updateSuccess", true);
        } else {
          message.error(res ? res : "null");
        }
      } else {
        let res = await addScenesConfig({
          ...formValue.value,
          operator: window.userInfo.username,
        });
        if (typeof res === "boolean") {
          message.success("添加成功");
          emits("addSuccess", true);
        } else {
          message.error(res ? res : "null");
        }
      }
    } else {
      message.error("更新失败,请重新填写");
      console.log("errors", errors);
    }
    messageReactive.destroy();
  });
};

//watch && computed

// hook

onMounted(() => {
  getSelects();
  addOrrevise.value =
    Object.keys(props.scenesConfigData).length !== 0;
  if (addOrrevise.value) {
    formValue.value = props.scenesConfigData;
  } else {
    formValue.value = {
      id: -1,
      appname: "",
      name: "",
      expname: "",
      buckettype: 1,
      flow_type: 0,
    };
  }
});
</script>

<style lang="less" scoped>
.n-form {
  .n-form-item {
    // margin: 0;
    :deep(.n-form-item-blank) {
      display: block !important;
      width: 100%;
      :deep(.n-input-number) {
        width: 100%;
        :deep(.n-button) {
          width: 100%;
        }
      }
      .successButton {
        width: 30%;
      }
    }
    .n-input {
      width: 100%;
    }
  }
}
</style>
