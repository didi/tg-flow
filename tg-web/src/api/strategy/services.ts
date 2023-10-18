import request from "@/api/request";
import {
  addOrUpdateWorkFlow,
  deleteWorkflowI,
  experimentConfigResI,
  experimentHistoryVersionI,
  experimentHistoryVersionParamsI,
  experimentRollHistoryVersionParamsI,
  experimentWorkflowNodeTypeResI,
  FlowChartParamsI,
  FlowChartResI, GetApolloResI,
  getExperimentConfigDimensionResI,
  getExperimentConfigParamsI,
  globalConfigBaseContentNullI,
  globalConfigBaseContentStringI,
  ModuleBranchRequestI, ModuleBranchRespI, ModulesRequestI, ModulesResI,
  operatorAndId,
  OperatorI,
  scenesConfigAppnameI,
  ScenesConfigI,
  scenesConfigNameI,
  ScenesConfigUpdateI,
  SystemConfigI,
  SystemConfigUpdateI,
} from "./types";
import {
  experimentAppname,
  filterData,
  scenesSelectAppName,
  scenesSelectName
} from "./utils";

export async function getWorkFlowChartByFlowIdV2(params: FlowChartParamsI) {
  const res = await request<FlowChartResI>(
    "/api/exportWorkFlow",
    params,
    "POST"
  );
  if (res) {
    if (res.content) {
      return JSON.parse(res.content);
    } else {
      return "";
    }
  } else {
    return "";
  }
}

// 拿取系统管理中的数据
export async function getSystemConfig(params: OperatorI) {
  const res = await request<SystemConfigI>(
    "/api/searchAppConfig",
    params,
    "POST"
  );
  if (res?.tag) {
    return res?.content;
  }
  return res?.err_msg;
}

// 更新系统管理中的数据
export async function updateSystemConfig(params: SystemConfigUpdateI) {
  const res = await request<SystemConfigI>(
    "/api/addOrUpdateAppConfig",
    params,
    "POST"
  );
  if (res?.tag) {
    return res;
  }
  return null;
}

// 删除系统管理中的数据
export async function deleteSystemConfig(params: operatorAndId) {
  const res = await request<SystemConfigI>(
    "/api/deleteAppConfig",
    params,
    "POST"
  );
  if (res?.tag) {
    return res;
  }
  return res?.err_msg;
}

// 导出系统管理中的部分数据
export async function exportSystemConfig(params: operatorAndId) {
  const res = await request<globalConfigBaseContentStringI>(
    "/api/exportAppConfig",
    params,
    "POST"
  );
  if (res?.tag) {
    return res;
  }
  return res?.err_msg;
}

// 提交系统管理中的数据到 Apollo
export async function commitSystemConfig(params: operatorAndId, is_kflower: boolean) {
  return await request<globalConfigBaseContentNullI>(
      "/api/submitAppConfig",
      {
        ...params,
        is_kflower,
      },
      "POST"
  )
}

// 查询系统对应的 Apollo 的状态
export async function getApolloInfo(params: operatorAndId, is_kflower: boolean) {
  return await request<GetApolloResI>(
      "/api/getApolloInfo",
      {
        ...params,
        is_kflower,
      },
      "POST"
  )
}


// 添加系统管理中的数据
export async function addSystemConfig(params: SystemConfigUpdateI) {
  const res = await request<SystemConfigI>(
    "/api/addOrUpdateAppConfig",
    params,
    "POST"
  );
  if (res?.tag) {
    return res;
  }
  return res?.err_msg;
}

// 拿到场景管理中的数据
export async function getScenesConfig(params: OperatorI) {
  const res = await request<ScenesConfigI>("/api/searchScene", params, "POST");
  if (res?.tag) {
    return res["content"];
  }
  return [];
}

// 获取模块数据
export async function getAppModules(params: ModulesRequestI) {
  const res = await request<ModulesResI>("/api/getAppModules", params, "POST");
  if (res?.tag) {
    return res["content"];
  }
    return {};
}

// 查询系统对应分支信息
export async function getGitBranches(params: ModuleBranchRequestI) {
    const res = await request<ModuleBranchRespI>("/api/getGitBranches", params, "POST");
    if (res?.tag) {
        return res["content"];
    }
    return [];
}

// 修改场景管理中的数据
export async function updateScenesConfig(params: ScenesConfigUpdateI) {
  const res = await request<globalConfigBaseContentNullI>(
    "/api/addOrUpdateScene",
    params,
    "POST"
  );
  if (res?.tag) {
    return true;
  }
  return res?.err_msg;
}

// 删除场景管理中的数据
export async function deleteScenesConfig(params: operatorAndId) {
  const res = await request<globalConfigBaseContentNullI>(
    "/api/deleteScene",
    params,
    "POST"
  );
  return !!res?.tag;

}

// 新增场景管理中的数据
export async function addScenesConfig(params: ScenesConfigUpdateI) {
  const res = await request<globalConfigBaseContentNullI>(
    "/api/addOrUpdateScene",
    params,
    "POST"
  );
  if (res?.tag) {
    return true;
  }
  return res?.err_msg;
}

// 拿到场景管理中的系统名称数据
export async function getScenesConfigAppname(params: OperatorI) {
  const res = await request<scenesConfigAppnameI>(
    "/api/autoGetAPPNameMenu",
    params,
    "POST"
  );
  if (res?.tag) {
    return scenesSelectAppName(res["content"]);
  }
  return res?.err_msg;
}

//拿到场景管理中的场景名称数据
export async function getScenesConfigName(params: OperatorI) {
  const res = await request<scenesConfigNameI>(
    "/api/autoGetSceneMenu",
    params,
    "POST"
  );
  if (res?.tag) {
    return scenesSelectName(res["content"]);
  }
  return res?.err_msg;
}

// 拿到实验管理中的数据
export async function getExperimentConfig(params: getExperimentConfigParamsI) {
  const res = await request<experimentConfigResI>(
    "/api/searchWorkFlow",
    params,
    "POST"
  );
  if (res) {
    if (res.tag) {
      return res["content"];
    } else return res.err_msg;
  } else {
    return null;
  }
}

// 拿到实验管理中的scenes
export async function getExperimentScenes(types: number) {
  const res = await request<scenesConfigNameI>(
    "/api/autoGetSceneMenu",
    {
      types,
    },
    "POST"
  );
  if (res?.tag) {
    return scenesSelectName(res["content"]);
  }
  return res?.err_msg;
}

// 拿到实验管理中的Appname
export async function getExperimentConfigAppname() {
  const res = await request<scenesConfigAppnameI>(
    "/api/autoGetAPPNameMenu",
    {},
    "POST"
  );
  if (res?.tag) {
    return experimentAppname(res["content"]);
  }
  return res?.err_msg;
}

// 拿到实验管理中的DimensionMenu
export async function getExperimentConfigDimension(scene_name: string) {
  const res = await request<getExperimentConfigDimensionResI>(
    "/api/autoGetDimensionMenu",
    { scene_name },
    "POST"
  );
  if (res?.tag) {
    return res["content"];
  }
  return res?.err_msg;
}

// 刷新实验管理中的NodeType
export async function refreshExperimentConfigNodeType(params:any) {
  const res = await request<globalConfigBaseContentNullI>(
    "/api/setNodeTypeList",
    params,
    "POST"
  );
  if (res?.tag) {
    return true;
  }
  return res?.err_msg;
}

// 复制实验管理中具体某个场景下的某一项
export async function addOrUpdateWorkflow(params: addOrUpdateWorkFlow) {
  const res = await request<globalConfigBaseContentNullI>(
    "/api/addOrUpdateWorkFlow",
    params,
    "POST"
  );
  if (res) {
    if (res.tag) {
      return res.tag;
    } else {
      return res.err_msg;
    }
  } else {
    return null;
  }
}

export async function updateWorkflowV2(params: addOrUpdateWorkFlow) {
  const res = await request<globalConfigBaseContentNullI>(
    "/api/importWorkFlow",
    params,
    "POST"
  );
  if (res) {
    if (res.tag) {
      return res.tag;
    } else {
      return res.err_msg;
    }
  } else {
    return null;
  }
}

// 删除实验管理中具体某个场景下的某一项
export async function deleteWorkflowData(params:deleteWorkflowI) {
  const res = await request<globalConfigBaseContentNullI>("/api/deleteWorkFlow",params,"POST")
  if (res) {
    if (res.tag) {
      return res.tag;
    } else {
      return res.err_msg;
    }
  } else {
    return null;
  }
}

// 拿到某个场景下的历史版本
export async function getScenesHistoryVersion(params:experimentHistoryVersionParamsI) {
  const res = await request<experimentHistoryVersionI>("/api/searchWorkflowVersion",params,"POST")
  if (res) {
    if (res.tag) {
      return filterData(res["content"]);
    } else return res.err_msg;
  } else {
    return null;
  }
}

// 回滚到某个场景下的历史版本
export async function backToHistoryVersion(params:experimentRollHistoryVersionParamsI) {
  const res = await request<globalConfigBaseContentNullI>("/api/rollBackWorkFlow",params,"POST")
  if (res) {
    if (res.tag) {
      return res.tag;
    } else {
      return res.err_msg;
    }
  } else {
    return null;
  }
}

// 导出实验管理中的某条数据
export async function exportexperimentConfig(params:{operator:string,workflowid:string}){
  const res = await request<globalConfigBaseContentStringI>("/api/exportWorkFlow",params,"POST")
  if (res?.tag) {
    return res['content'];
  }
  return res?.err_msg;
}

// 导入实验管理workflow Json数据
export async function requireWorkflowConfig(params:{workflowid:number,dataJson:string}){
  const res = await request<FlowChartResI>('/api/importWorkFlow',params,"POST")
  if (res?.tag) {
    return true;
  }
  return res?.err_msg;
}

//拿到实验管理中的nodetype
export async function getExperimentWorkflowNodeType(params:FlowChartParamsI){
  const res = await request<experimentWorkflowNodeTypeResI>("/api/getNodeTypeList",params,"POST")
  if(res){
    return res['content']
  }else{
    return null
  }
}
