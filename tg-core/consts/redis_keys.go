package consts

import ()

/**
	策略服务公用业务相关的key，一律以"strategy_"作为前缀
**/
const (
	//TODO ZYF待删除
	StrategySceneAndModuleMap      = "strategy_sceneAndModule_map"
	StrategyWorkflowMap            = "strategy_workflow_map"
	StrategyDimensionMap           = "strategy_dimension_map"

	//算法模型索引key
	RedisKeyAlgoModelConfig = "strategy_algo_model_config"
	RedisKeyRankerConfigMap = "strategy_rediskey_ranker_config_map"

	//机器ip
	RedisKeyMachineIp = "rediskey_machine_ip"

	//systemConfig信息
	RedisKeySystemConf = "strategy_system_conf"

	//特征数据类型前缀
	FeatureDataType = "feature_data_type_scene_%v"

	//召回配置
	RedisKeyRecallConfig = "strategy_recall_config"
)
