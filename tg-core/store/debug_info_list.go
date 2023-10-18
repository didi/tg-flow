/**
    @Description:
    @Author:zhouzichun
    @Date:2022/4/25
**/

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/model"
	"go.intra.xiaojukeji.com/apollo/apollo-golang-sdk-v2"
	apolloModel "go.intra.xiaojukeji.com/apollo/apollo-golang-sdk-v2/model"
	"strconv"
	"time"
)

type DebugResponse struct {
	UserId          string                       `json:"user_id"`
	AppId           int64                        `json:"app_id"`
	SceneId         int64                        `json:"scene_id"`
	FlowId          int64                        `json:"flow_id"`
	TimeStamp       int64                        `json:"time_stamp"`
	ActionResultMap map[string]map[string]string `json:"action_result_map"`
}

func Debug(sc *model.StrategyContext) error {
	redisKey := fmt.Sprintf("debug_info_%v_%v", sc.AppName, sc.Phone)
	debugResponse := &DebugResponse{
		AppId:           sc.AppId,
		SceneId:         sc.SceneId,
		FlowId:          sc.FlowId,
		TimeStamp:       time.Now().Unix(),
		ActionResultMap: sc.ActionResultMap,
	}
	debugResponseStr, err := json.Marshal(debugResponse)
	if err != nil {
		return err
	}
	var limitLength int64
	limitLength = 100
	user := apolloModel.NewUser(sc.Phone).With("phone", sc.Phone)
	apolloToggle, apolloErr := apollo.FeatureToggle("online_debug_config", user)
	if apolloErr != nil || !apolloToggle.IsAllow() {
		return apolloErr
	} else {
		limitLengthStr := apolloToggle.GetAssignment().GetParameter("limit_num", "100")
		limitLength, _ = strconv.ParseInt(limitLengthStr, 10, 64)
	}
	ctx := context.TODO()
	length, err := redis.Handler.LPush(ctx, redisKey, string(debugResponseStr))
	if err != nil {
		return err
	}
	_, err = redis.Handler.Expire(ctx, redisKey, 24*60*60)
	if err != nil {
		return err
	}
	if length > limitLength {
		_, err = redis.Handler.LTrim(ctx, redisKey, int(length-limitLength), int(length-1))
		if err != nil {
			return err
		}
	}
	return nil
}
