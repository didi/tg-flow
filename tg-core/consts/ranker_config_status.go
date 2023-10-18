package consts 

const (
	RANKER_STATUS_UNKNOWN = -1			//未知
	RANKER_STATUS_TOONLINE  = 1			//模型待上线
	RANKER_STATUS_ONLINE_FAIL = 2		//模型上线失败
	RANKER_STATUS_ONLINE = 3			//模型上线成功
)

/**
1. 新增配置时，scene_id、algo_name、model_name（以后简称sam）均不许为空，初始状态统一为status=1
2. 更新配置时，sam改为非空的sam，则状态置为模型待加载；sam改为空的sam，则状态置为模型加载成功；
3. 删除配置时，必须满足当前状态为空闲时才允许删除（当ranker发现过期模型卸载成功且sam均为空时，会将状态置为空闲）
4. 更新配置时，不允许更新ip字段，如果要删除该ip对应的记录，需要先使之为空闲状态
**/