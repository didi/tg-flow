package dispatcher

import (
	"context"
	"fmt"
	"git.xiaojukeji.com/gobiz/logger"
	ctxtrace "git.xiaojukeji.com/lego/context-go"
	"github.com/didi/tg-flow/tg-core/common/utils"
	"strings"
	"time"
)

func WriteCustomLog(ctx context.Context, publicKey string, filePrefix string, pairs map[string]interface{}, tagName string) {
	defer utils.RecoverPanic(ctx, tagName)
	mergeLog(pairs)
	var kvs []string
	kvs = append(kvs, publicKey)
	kvs = append(kvs, "timestamp="+time.Now().Format("2006-01-02 15:04:05"))
	kvs = append(kvs, ctxtrace.FormatCtx(ctx))
	for k, v := range pairs {
		if v == nil {
			v = "NULL"
		}
		kvs = append(kvs, k+"="+fmt.Sprint(v))
	}
	//记本地日志
	logger.Track(filePrefix, strings.Join(kvs, "||"))
}

