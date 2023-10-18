/**
    @Description:
    @Author:zhouzichun
    @Date:2023/5/7
**/

package conf

import (
	"github.com/didi/tg-flow/tg-core/conf"
	"github.com/didi/tg-flow/tg-core/router"
	"log"
	"tg-service/constant"
	"tg-service/controller"
)

func init() {
	log.Println("point-admin router init...")
	svc := new(controller.RPCService)
	router.Handler.POST("/checklogin", svc.CheckLogin)
	router.Handler.POST("/login", svc.Login)
	router.Handler.POST("/logout", svc.Logout)

	//workflow相关操作（新）
	router.Handler.POST("/searchWorkFlow", svc.SearchWorkFlow)
	router.Handler.POST("/addOrUpdateWorkFlow", svc.AddOrUpdateWorkFlow)
	router.Handler.POST("/deleteWorkFlow", svc.DeleteWorkFlow)
	router.Handler.POST("/exportWorkFlow", svc.ExportWorkFlow)
	router.Handler.POST("/getWorkFlowChart", svc.GetWorkFlowChart)
	router.Handler.POST("/saveWorkFlowChart", svc.SaveWorkFlowChart)
	router.Handler.POST("/importWorkFlow", svc.ImportWorkFlow)
}

func InitConfig() {
	//初始化:获取配置文件中sso服务接口地址
	var err error
	if constant.Env, err = conf.Handler.GetSetting("env", "type"); err != nil {
		log.Println("Fail to init env!", err)
		return
	}
	log.Println("point-admin config init finished")
}
