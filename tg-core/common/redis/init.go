package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"git.xiaojukeji.com/gobiz/config"
	"git.xiaojukeji.com/nuwa/golibs/redis"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/conf"
	"github.com/didi/tg-flow/tg-core/consts"
	"github.com/didi/tg-flow/tg-core/model"
	"log"
	"strings"
	"time"
)

var Handler *redis.Manager
var FeatureRedisHandler *redis.Manager
var Handlers map[string]*redis.Manager

const (
	SECTION               = "redis"
	FEATURE_REDIS_SECTION = "feature_redis"
	ErrNil                = "redigo: nil returned"
	//notFindKey            = "can not find key"
)

func InitRedisHandlers(sections []string) {
	if len(sections) <1 {
		return
	}

	handlers := make(map[string]*redis.Manager)
	for _, section := range sections {
		redisManager, err := NewManagerFromConf(conf.Handler, section)
		if err != nil || redisManager == nil {
			log.Fatal("Init redis client["+ section + "] error: ", err)
		}

		handlers[section] = redisManager
		//为兼容内置Handler，此处先特殊处理一下
		if section == SECTION {
			Handler = redisManager
		}else if section == FEATURE_REDIS_SECTION {
			FeatureRedisHandler = redisManager
		}
		log.Println("Init redis client["+ section + "] successful !!!")
	}
	Handlers = handlers
}

//redis.go支持的参数太少，包一下
func NewManagerFromConf(cfg config.Configer, sec string, opt ...redis.Option) (*redis.Manager, error) {
	var opts []redis.Option
	addrs, err := cfg.GetSetting(sec, "addrs")
	if err != nil {
		return nil, err
	}

	servers := strings.Split(addrs, ",")
	auth, err := cfg.GetSetting(sec, "auth")
	if err != nil {
		return nil, err
	}
	disfEnable, _ := cfg.GetBoolSetting(sec, "disf_enable")
	if disfEnable {
		sn, err := cfg.GetSetting(sec, "service_name")
		if err != nil {
			return nil, err
		}
		opts = append(opts, redis.EnableDisf())
		opts = append(opts, redis.DisfServiceName(sn))
	}

	poolSize, err := cfg.GetIntSetting(sec, "pool_size")
	if err == nil {
		opts = append(opts, redis.SetPoolSize(poolSize))
	}

	maxConn, err := cfg.GetIntSetting(sec, "max_conn")
	if err == nil {
		opts = append(opts, redis.SetMaxConn(maxConn))
	}

	connTimeout, err := cfg.GetIntSetting(sec, "conn_timeout")
	if err == nil {
		opts = append(opts, redis.SetConnectTimeout(time.Millisecond*time.Duration(connTimeout)))
	}

	readTimeout, err := cfg.GetIntSetting(sec, "read_timeout")
	if err == nil {
		opts = append(opts, redis.SetReadTimeout(time.Millisecond*time.Duration(readTimeout)))
	}

	writeTimeout, err := cfg.GetIntSetting(sec, "write_timeout")
	if err == nil {
		opts = append(opts, redis.SetWriteTimeout(time.Millisecond*time.Duration(writeTimeout)))
	}

	maxTryTimes, err := cfg.GetIntSetting(sec, "max_try_times")
	if err == nil {
		opts = append(opts, redis.SetMaxTryTimes(maxTryTimes))
	}

	opts = append(opts, opt...)

	return redis.NewManager(servers, auth, opts...)
}

//TODO:write to redis
func WriteRedis(sc *model.StrategyContext, redisKey string, info interface{}, etype string) {
	infoByte, err := json.Marshal(info)
	if err != nil {
		tlog.LogError(context.TODO(), sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.jsonMarshal", etype), fmt.Sprintf("info=%v", info), err)
		return
	}
	if _, err := Handler.Set(context.TODO(), redisKey, infoByte); err != nil {
		tlog.LogError(context.TODO(), sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteRedis", etype), fmt.Sprintf("info=%v", info), err)
		return
	}
}

//write to fusion
func SetFusion(ctx context.Context, sc *model.StrategyContext, redisKey string, info interface{}, etype string) {
	infoByte, err := json.Marshal(info)
	if err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.jsonMarshal", etype), fmt.Sprintf("info=%v", info), err)
		return
	}
	WriteFushion(ctx, sc, redisKey, infoByte, etype)
}

func WriteFushion(ctx context.Context, sc *model.StrategyContext, redisKey string, info []byte, etype string) {
	if _, err := Handler.Set(ctx, redisKey, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushion", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//write to fusion
func WriteFushionEx(ctx context.Context, sc *model.StrategyContext, redisKey string, info []byte, expireTime int, etype string) {
	if _, err := Handler.SetEx(ctx, redisKey, expireTime, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushionEx", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//get from fusion
func GetFushion(ctx context.Context, sc *model.StrategyContext, redisKey string, etype string) string {
	info, err := Handler.Get(ctx, redisKey)
	if err != nil && err.Error() != ErrNil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.GetFushion", etype), fmt.Sprintf("info=%v", info), err)
		return ""
	}
	return info
}

func WToFusion(ctx context.Context, sc *model.StrategyContext, redisKey string, info interface{}, etype string) {
	if _, err := Handler.Set(ctx, redisKey, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushion", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//write to fusion
func WToFusionEx(ctx context.Context, sc *model.StrategyContext, redisKey string, info interface{}, expireTime int, etype string) {
	if _, err := Handler.SetEx(ctx, redisKey, expireTime, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushionEx", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//复制以上方法，新增一个redis.handle参数
func HSetRedis(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, subKey string, info interface{}, etype string) {
	infoByte, err := json.Marshal(info)
	if err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.HSetRedisMarshal", etype), fmt.Sprintf("info=%v", info), err)
		return
	}
	if _, err = fusionHandler.HSet(ctx, redisKey, subKey, infoByte); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.HSetRedis", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//write to fusion
func SetRedis(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, info interface{}, etype string) {
	infoByte, err := json.Marshal(info)
	if err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.SetRedis", etype), fmt.Sprintf("info=%v", info), err)
		return
	}
	WriteFushionByte(ctx, sc, fusionHandler, redisKey, infoByte, etype)
}

func WriteFushionByte(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, info []byte, etype string) {
	if _, err := fusionHandler.Set(ctx, redisKey, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushionByte", etype), fmt.Sprintf("info=%v", info), err)
	}
}

func WriteFushionInterface(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, info interface{}, etype string) {
	if _, err := fusionHandler.Set(ctx, redisKey, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushionInterface", etype), fmt.Sprintf("info=%v", info), err)
	}
}

func SetRedisEx(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, info interface{}, expireTime int, etype string) {
	infoByte, err := json.Marshal(info)
	if err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.SetRedisEx", etype), fmt.Sprintf("info=%v", info), err)
		return
	}
	WriteFushionByteEx(ctx, sc, fusionHandler, redisKey, infoByte, expireTime, etype)
}

//write to fusion
func WriteFushionByteEx(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, info []byte, expireTime int, etype string) {
	if _, err := fusionHandler.SetEx(ctx, redisKey, expireTime, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushionByteEx", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//write to fusion
func WriteFushionInterfaceEx(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, info interface{}, expireTime int, etype string) {
	if _, err := fusionHandler.SetEx(ctx, redisKey, expireTime, info); err != nil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.WriteFushionInterfaceEx", etype), fmt.Sprintf("info=%v", info), err)
	}
}

//get from fusion
func GetRedis(ctx context.Context, sc *model.StrategyContext, fusionHandler *redis.Manager, redisKey string, etype string) string {
	info, err := fusionHandler.Get(ctx, redisKey)
	if err != nil && err.Error() != ErrNil {
		tlog.LogError(ctx, sc, consts.DLTagFushion, fmt.Sprintf("%v.redis.GetRedis", etype), fmt.Sprintf("info=%v", info), err)
		return ""
	}
	return info
}
