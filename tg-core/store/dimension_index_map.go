package store

import (
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/consts"
	"github.com/didi/tg-flow/tg-core/model"
)

var DimensionIndexObj *model.DimensionIndex = &model.DimensionIndex{}

func LoadDimensionIndex() error {
	dimensionIndexString, err := redis.Handler.Get(context.TODO(), consts.StrategyDimensionMap)
	if err != nil && err.Error() != redis.ErrNil {
		return err
	}

	var dimensionIndexTemp *model.DimensionIndex
	err = json.Unmarshal([]byte(dimensionIndexString), &dimensionIndexTemp)
	if err != nil {
		return err
	}

	DimensionIndexObj.DimensionMap = dimensionIndexTemp.DimensionMap
	DimensionIndexObj.UpdateTime = dimensionIndexTemp.UpdateTime

	return nil
}
