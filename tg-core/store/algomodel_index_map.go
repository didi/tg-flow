package store

import (
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/consts"
	"github.com/didi/tg-flow/tg-core/model"
)

var AlgoModelIndexObj *model.AlgoModelIndex = &model.AlgoModelIndex{}

func LoadAlgoModelIndex() error {
	algoModelIndexByte, err := redis.Handler.Get(context.TODO(), consts.RedisKeyAlgoModelConfig)
	if err != nil && err.Error() != redis.ErrNil {
		return err
	}

	var algoModelIndexTemp *model.AlgoModelIndex
	err = json.Unmarshal([]byte(algoModelIndexByte), &algoModelIndexTemp)
	if err != nil {
		return err
	}
	AlgoModelIndexObj.IndexVersion = algoModelIndexTemp.IndexVersion
	AlgoModelIndexObj.IndexMap = algoModelIndexTemp.IndexMap
	AlgoModelIndexObj.UpdateTime = algoModelIndexTemp.UpdateTime
	return nil
}
