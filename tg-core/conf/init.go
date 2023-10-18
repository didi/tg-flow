package conf

import (
	"git.xiaojukeji.com/gobiz/config"
	"github.com/didi/tg-flow/tg-core/common/path"
	"git.xiaojukeji.com/nuwa/go-monitor"
	"go.intra.xiaojukeji.com/platform-ha/onekey-degrade_sdk_go/degrade"
	"log"
	_ "net/http/pprof"
)

var Handler config.Configer

//center服务url地址
var CENTER_URL string
var CENTER_HTTP_TIMEOUT int64

//911限流预案名称
var RateLimitName string

//911降级预案名称
var DownGradeName string

//系统运行环境：正式、预发、线下
var ENV string

func init() {
	log.Println("tg-core conf init start, path.Root=" + path.Root)
	var err error
	Handler, err = config.New(path.Root + "/conf/app.conf")
	if err != nil {
		log.Println("Configer init error:", err, "(本地环境及单元测试可忽略该错误)")
		// 适配读取配置文件 解决本地单元测试中无法正确读取配置路径的问题
		/*Handler, err = config.New(path.ConfLocal + "/conf/app.conf")
		if err != nil {
			log.Fatal("读取本地配置失败 =>", path.ConfLocal)
		}*/
	}

	initRateLimitName()

	initDownGradeName()

	initEnv()

	nuwaMonitor()
}

//初始化系统911限流预案名称
func initRateLimitName() {
	var err error
	if RateLimitName, err = Handler.GetSetting("ratelimit", "name"); err != nil {
		log.Println("Fail to get ratelimit name!", err)
		return
	}
}

//初始化系统911降级预案名称
func initDownGradeName() {
	var downGradeName string
	var err error
	//获取降级预案配置
	if downGradeName, err = Handler.GetSetting("downgrade", "name"); err != nil {
		log.Println("Fail to get downgrade name!", err)
		return
	} else {
		DownGradeName = downGradeName
	}

	//初始化911sdk
	configMeta := degrade.NewConfigMeta()
	configMeta.Add(DownGradeName, degrade.SWITCH_CONFIG)
	if err := degrade.Init(configMeta); err != nil {
		log.Println("Fail to init 911 sdk!", err)
		return
	}
}

//性能监控
func nuwaMonitor() {
	var isopen bool
	var port string
	var err error
	log.Println("tg-core nuwaMonitor init start...")
	if isopen, err = Handler.GetBoolSetting("pprof", "isopen"); err != nil {
		log.Println("Fail to get nuwaMonitor isopen config!", err)
		return
	}
	if port, err = Handler.GetSetting("pprof", "port"); err != nil {
		log.Println("Fail to get nuwaMonitor port config!", err)
		return
	}
	if isopen {
		go monitor.Start(":"+port, monitor.AllPlugin)
	}
}

//初始化系统运行环境变量
func initEnv() {
	var err error
	if ENV, err = Handler.GetSetting("env", "name"); err != nil {
		return
	}
	log.Println("tg-core env init end")
}
