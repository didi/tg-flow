package utils

import (
	"context"
	"fmt"
	"github.com/didi/tg-flow/common/tlog"
)

func Recover(ctx context.Context, tag string) {
	if err := recover(); err != nil {
		tlog.Handler.ErrorCount(ctx, tag, fmt.Sprintf("Recover system panic : %v", err))
	}
}
