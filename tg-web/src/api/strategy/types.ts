export interface OperatorI {
  operator: string;
}
interface globalConfigBaseI {
  err_msg: string;
  tag: boolean;
  typenum: number;
}

export interface globalConfigBaseContentStringI extends globalConfigBaseI {
  content: string;
}
interface ApolloNamespaceI {
  id: number,
}

interface ApolloServiceI {
  id: number,
  appId: number,
}

interface ApolloResponseI {
  id: number,
  services: Array<ApolloServiceI>,
  namespace: ApolloNamespaceI,
  configStatus: string,
}

export interface GetApolloResI extends globalConfigBaseI {
  content: ApolloResponseI
}
export interface FlowChartParamsI extends OperatorI {
  workflowid?: number;
  scene_name?: string;
}

export interface FlowChartResI extends globalConfigBaseI {
  content: string | null
}

// nodes: {label: string; id: string ...}
export interface SystemConfigI extends globalConfigBaseI {
  content: Array<SystemConfigContentI>
}

export interface SystemConfigContentI extends SystemConfigUpdateI, createAndUpdate {
}

export interface SystemConfigUpdateI extends operatorAndId {
  app_name: string;
  node_name: string;
  old_id?: number;
  machine_room: string;
  git_url: string;
}

export interface operatorAndId {
  operator: string;
  id: number;
}

export interface ScenesConfigI extends globalConfigBaseI {
  content: Array<ScenesConfigContentI>
}

interface createAndUpdate {
  create_time: string;
  update_time: string;
}

export interface ScenesConfigContentI extends ScenesConfigUpdateI {
  appid: number;
  bifrost_config: string;
  createtime: string;
  namezh: string;
  updatetime: string;
}

// ==================== Modules ====================

export interface ModulesRequestI extends OperatorI {
  app_name: string,
  project_branch: string,
}

export interface ModulesContentI {
  name: string,
  type: string,
  desc: string,
  workflow_list: Array<number>
  scene_list: Array<string>
  create_time: string,
}

export interface ModulesResI extends globalConfigBaseI {
    content: Record<string, Array<ModulesContentI>>
}

export interface ModuleBranchRequestI extends OperatorI {
  app_name: string,
}

export interface ModuleBranchRespI extends globalConfigBaseI {
  content: Array<string>
}


export interface ScenesConfigUpdateI {
  oldid?: number;
  operator: string;
  id: number;
  name: string;
  appname: string;
  expname: string;
  buckettype: number;
  flow_type: number;
}

export interface globalConfigBaseContentNullI extends globalConfigBaseI {
  content: null;
}
export interface scenesConfigAppnameI extends globalConfigBaseI {
  content: Array<scenesConfigAppnameContentI>
}

export interface scenesConfigAppnameContentI {
  AppId: string;
  AppName: string;
}

export interface scenesConfigNameContentI {
  SceneId: string;
  SceneName: string;
}

export interface scenesConfigNameI extends globalConfigBaseI {
  content: Array<scenesConfigNameContentI>;
}
export interface getExperimentConfigParamsI extends OperatorI {
  dimension_id: string;
  scenename: string;
}

export interface getExperimentConfigDimensionResI extends globalConfigBaseI {
  content: Array<getExperimentConfigDimensionResContentI>;
}

export interface getExperimentConfigDimensionResContentI {
  DimensionId: string;
  DimensionName: string;
}
export interface experimentConfigResI extends globalConfigBaseI {
  content: Array<experimentConfigResContentI>
}

export interface experimentConfigResContentI extends OperatorI {
  configured: boolean;
  createtime: string;
  updatetime: string;
  defult: string;
  dimension_id: number;
  groupname: string;
  manual_slot_ids: string;
  modules: string;
  oldworkflowid: string;
  proportion: string;
  range1: string;
  range2: string;
  remark: string;
  scene_id: number;
  scenename: string;
  showmodules: string;
  workflowid: string;
}

export interface addOrUpdateWorkFlow extends deleteWorkflowI {
  oldworkflowid: string;
  modules: string;
  proportion: number;
  defult: number;
  remark: string;
  dataJson?: string;
}

export interface deleteWorkflowI extends experimentHistoryVersionParamsI {
  workflowid: string;
}

export interface experimentHistoryVersionI extends globalConfigBaseI {
  content: Array<experimentHistoryVersionResI>;
}

export interface experimentHistoryVersionParamsI extends OperatorI {
  scenename: string;
  dimension_id: number;
}

export interface experimentHistoryVersionResI extends OperatorI, experimentHistoryVersionResFilterI {
  scenename: string;
}

export interface experimentHistoryVersionResFilterI {
  dimension_id: number;
  version_create_time: string;
  version_id: number;
}

export interface experimentRollHistoryVersionParamsI extends OperatorI {
  version_id: number;
  dimension_id: number;
  scene_name: string;
}

export interface experimentWorkflowNodeTypeResI extends globalConfigBaseI {
  content: Array<experimentWorkflowNodeTypeResContentI>;
}

export interface experimentWorkflowNodeTypeResContentI {
  node_name: string[];
  node_type: string;
}
