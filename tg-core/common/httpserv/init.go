package httpserv

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/conf"
	"github.com/didi/tg-flow/tg-core/router"
	"git.xiaojukeji.com/nuwa/nuwa-go-httpclient"
	ngs "git.xiaojukeji.com/nuwa/nuwa-go-httpserver/v2"
	"log"
	"time"
)

var Handler ngs.HTTPServer

func init() {
	log.Println("tg-core httpserv init start...")
	Handler = ngs.New(conf.Handler)
	Handler.SetLogger(tlog.Handler)
	Handler.SetRouter(router.Handler)

	Handler.RegisterOnShutdown(func() {
	})
}

func SendHttpRequest(url string, timeOut int64, parmasMap interface{}) ([]byte, error) {
	_, resp, err := httpclient.PostJSON(context.TODO(), url, parmasMap, httpclient.SetTimeout(time.Millisecond*time.Duration(timeOut)))
	if err != nil { //重试2次
		for i := 0; i < 2; i++ {
			_, resp, err = httpclient.PostJSON(context.TODO(), url, parmasMap, httpclient.SetTimeout(time.Millisecond*time.Duration(timeOut)))
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}
