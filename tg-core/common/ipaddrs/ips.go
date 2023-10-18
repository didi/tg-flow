package ipaddrs

import (
	"context"
	"fmt"
	"git.xiaojukeji.com/gobiz/utils"
	"net"
	"net/http"
	"strings"
)

//ZYF 以下从nuwa-go-httpserver/middleware/trace.go里面拷贝过来的,便于实现GetClientIp
type ctxKey struct {
	name string
}

var (
	// RequestTimeKey ...
	RequestTimeKey = ctxKey{"requestTime"}
	// TraceRecordKey ...
	TraceRecordKey = ctxKey{"traceRecord"}
	// ExtraInfoReqOut ...
	ExtraInfoReqOut = ctxKey{"extraReqOut"}
)

// TraceRecord ...
type TraceRecord struct {
	TraceId      string
	SpanId       string
	HintCode     string
	HintContent  string
	URL          string
	Params       string
	Method       string
	Host         string
	From         string
	FormatString string
}
//ZYF 以上从nuwa-go-httpserver/middleware/trace.go里面拷贝过来的，便于实现GetClientIp，跟nuwa/trace/trace.go中的Trace定义高度相似，field是后者的子集.

func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("get local ip is fail,addrs is %v", addrs)
}

func GetClientIp(ctx context.Context) string {
	rec, ok := ctx.Value(TraceRecordKey).(TraceRecord)
	var clientIp string
	if ok {
		clientIp = formatIp(rec.From)
	} else {
		clientIp = "0.0.0.0"
	}
	return clientIp
}

func GetClientAddr(r *http.Request) string {
	remoteAddr := utils.GetClientAddr(r)
	return formatIp(remoteAddr)
}

func formatIp(remoteAddr string) string {
	if remoteAddr == "::1" {
		return "127.0.0.1"
	}

	strs := strings.Split(remoteAddr, ":")
	if strs != nil && len(strs) > 0 {
		remoteAddr = strs[0]
	}

	return remoteAddr
}

