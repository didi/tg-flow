package inneraction

import (
	"context"
	"github.com/didi/tg-flow/tg-core/model"
	"github.com/didi/tg-flow/tg-core/wfengine"
)

type TimeoutAction struct {
	wfengine.ModelBase
}

func (t TimeoutAction) DoAction(ctx context.Context, sc *model.StrategyContext) interface{} {
	//nothing to do
	//fmt.Println(time.Now(),"you can do nothing")
	return nil
}

func (m TimeoutAction) OnTimeout(context.Context, *model.StrategyContext) {
	//this is a sample, you can do sth when timeout happen
	//fmt.Println("default OnTimeout function, nothing to do !!!")
}