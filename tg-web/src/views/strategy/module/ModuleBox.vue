<template>
    <n-spin :show="show">
        <template #icon>
            <n-icon>
                <Reload/>
            </n-icon>
        </template>

        <div class="module-box">
            <el-collapse v-model="activeName" accordion>
                <el-collapse-item v-for="(modules, group) in data" :name="group">
                    <!--suppress VueUnrecognizedSlot -->
                    <template #title>
                        <span style="margin: 20px">{{ group }}</span>
                    </template>
                    <div id="col-modules">
                        <n-grid x-gap="12" :cols="3" y-gap="12">
                            <n-gi v-for="module in modules" :key="module.name">
                                <n-card :title="module.name"
                                        :segmented="{content: true}"
                                >
                                    <template #header>
                                        <span class="card-header">{{ module.name }}</span>
                                    </template>
                                    <template #header-extra>
                                        <n-tag round :bordered="false">{{ module.type }}</n-tag>
                                    </template>
                                    <n-descriptions :column="1">
                                        <n-descriptions-item label="描述">
                                            <span>{{ formatDesc(module.desc) }}</span>
                                        </n-descriptions-item>

                                        <n-descriptions-item label="创建时间">
                                            <span>{{ module.create_time }}</span>
                                        </n-descriptions-item>
                                    </n-descriptions>

                                    <template #footer>
                                        <div class="ref-tag">
                                            <n-tag type="success" :bordered="false">
                                                被引用的场景
                                            </n-tag>
                                            <el-link v-for="sceneid in module.scene_list" type="primary"
                                                     @click="clickScene(sceneid)">
                                                <SvgIcon name="ele-Guide"/>
                                                {{ sceneid }}
                                            </el-link>
                                        </div>
                                        <div class="ref-tag">
                                            <n-tag type="info" :bordered="false">
                                                被引用的流程
                                            </n-tag>
                                            <el-link v-for="wid in module.workflow_list" type="primary"
                                                     @click="clickWorkflow(wid)">
                                                <SvgIcon name="ele-Share"/>
                                                <span> </span>
                                                {{ wid }}
                                            </el-link>
                                        </div>
                                    </template>
                                </n-card>
                            </n-gi>
                        </n-grid>
                    </div>
                </el-collapse-item>
            </el-collapse>

        </div>
    </n-spin>

</template>

<script setup lang="ts">
import {onMounted, ref, watch} from "vue";
import SvgIcon from "@/components/svgIcon/index.vue";
import {getAppModules, ModulesContentI} from "@/api/strategy";
import {useRouter} from "vue-router";
import {useMessage} from "naive-ui";
import {Reload} from "@vicons/ionicons5";


// ===================== props ======================

const props = defineProps({
    // prop for app name
    appName: {
        type: String,
        default: "",
    },
    // prop for branch name
    branchName: {
        type: String,
        default: "master",
    },
});

// ====================== data ======================

// show for loading spin
const show = ref(false);

// current appName
const activeName = ref("1")
const data = ref({} as Record<string, Array<ModulesContentI>>);

// ====================== methods ======================

// message box
const message = useMessage();
const router = useRouter();

const clickScene = (sceneName: string) => {
    // route to scene page
    router.push({
        path: "/strategy-scenes",
        query: {
            sceneName: sceneName,
            appName: props.appName,
        },
    })
}

const clickWorkflow = (workflowId: number) => {
    router.push({
        path: "/strategy-experiment",
        query: {
            workflowId: workflowId,
            appName: props.appName,
        },
    })
}

const formatDesc = (desc: string) => {
    // if desc is empty, use defaut desc
    if (desc === "") {
        return "这个模块暂时没有描述";
    }
}

const queryAppModules = async () => {
    show.value = true;
    data.value = await getAppModules({
        operator: window.userInfo.username,
        app_name: props.appName,
        project_branch: props.branchName,
    });
    show.value = false;
};

// ====================== watch ======================

watch(() => props.appName, () => {
    queryAppModules();
});
watch(() => props.branchName, () => {
    queryAppModules();
});

// ====================== onMounted ======================
onMounted(() => {
    queryAppModules();
});


</script>

<style lang="less" scoped>

.module-box {
  height: 300px;
}

.n-card {
  max-width: 800px;
  height: auto;

  .ref-tag {
    margin: 10px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

#col-modules {
  margin: 40px;
}

.el-link {
  margin-right: 8px;
  margin-left: 8px;
}

.el-link .el-icon--right.el-icon {
  vertical-align: text-bottom;
}

</style>
