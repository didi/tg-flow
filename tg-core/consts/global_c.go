package consts

import ()

const (
	FEATURE_VALUE_DEFAULT_FLOAT = -1 //TODO 全部工程用NUMBER替换后删除

	FEATURE_VALUE_NEGATIVE_ONE = "-1" //TODO 全部工程用STRING替换后删除

	FEATURE_VALUE_DEFAULT_NUMBER = -1

	FEATURE_VALUE_DEFAULT_STRING = "-1"

	LOCALCACHE_TTL = 1800 //unit:s

	USER_LOCALCACHE_TTL = 3600 //unit:s

	FLOW_BY_ONLINE_RANDOM = 0 //在线随机分流
	FLOW_BY_ONLINE_SALT   = 1 //在线salt分流
	FLOW_BY_OFFLINE       = 2 //离线分流
	FLOW_BY_APOLLO        = 3 //apollo分流
	
	DefaultWorkflowVersion = "1.0"

	//strategycontext中set的key常量
	RANKER_IP = "ranker_ip"

)

//fusion中按在线salt分流和离线分流数据的表名
var SALT_ONLINE_KEY string
var OFFLINE_KEY string
