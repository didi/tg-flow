package utils

import (
	"context"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/model"
	"sync"
)
//////////////////////////以下都是带删除
func Recover(ctx context.Context, sc *model.StrategyContext, tag string, subType string) {
	if err := recover(); err != nil {
		tlog.LogError(ctx, sc, tag, subType, "Recover system panic", fmt.Errorf("%v", err))
	}
}

func RecoverThread(ctx context.Context, sc *model.StrategyContext, tag string, subType string, ch chan int) {
	if err := recover(); err != nil {
		tlog.LogError(ctx, sc, tag, subType, "RecoverThread system panic", fmt.Errorf("%v", err))
		ch <- 1
	}
}

func RecoverThreadByWg(ctx context.Context, sc *model.StrategyContext, tag string, subType string, wg *sync.WaitGroup) {
	if err := recover(); err != nil {
		tlog.LogError(ctx, sc, tag, subType, "RecoverThreadByWg system panic", fmt.Errorf("%v", err))
		wg.Done()
	}
}

////////////////////////////////////以上都是待删除//////////////////

func RecoverPanic(ctx context.Context, tag string) {
	if err := recover(); err != nil {
		tlog.ErrorCount(ctx, tag, fmt.Sprintf("Recover system panic : %v", err))
	}
}