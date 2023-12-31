/**
 * This file is auto-generated by dirpcgen don't modify manaully
 *
 * Copyright (c) 2018 didichuxing.com, Inc. All Rights Reserved
 *
 * Generated-date: 2018-07-13
 */

package controller

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/common/utils"
	ngs "git.xiaojukeji.com/nuwa/nuwa-go-httpserver/v2"
	"net/http"
	"tg-service/common/logs"
	"tg-service/constant"
	"tg-service/idl"
	"tg-service/logic"
)

type RPCService struct {
	ngs.BaseController
}

var isNewAdmin bool

/**
 * 登录校验1
 */
func (h RPCService) CheckLogin(w http.ResponseWriter, r *http.Request) {
	defer utils.Recover(context.TODO(), nil, logs.DLTagSystemPanic, "logic_CheckLogin")
	responseInfo := &idl.ResponseInfo{}
	data := logic.ParsRequestParam(r)
	tlog.Handler.Infof(context.TODO(), logs.DLTagProcessLog, "etype=controller_CheckLogin||econtent=%v||err=", data)
	if data.UserCookie == "" {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "用户cookie为空，需要登录"
	} else {
		tag, err := logic.Login(r.FormValue("usercookie"), r.FormValue("passcookie"))
		if tag {
			responseInfo.Tag = true
			responseInfo.ErrMsg = "success"
		} else {
			responseInfo.Tag = false
			responseInfo.TypeNum = 1 //登录校验不通过
			responseInfo.ErrMsg = err
		}
	}
	logic.EchoJSON(w, r, responseInfo)
}

/**
 * 登录校验2
 */
func (h RPCService) Login(w http.ResponseWriter, r *http.Request) {
	defer utils.Recover(context.TODO(), nil, logs.DLTagSystemPanic, "logic_Login")
	responseInfo := &idl.ResponseInfo{}
	data := logic.ParsRequestParam(r)
	data.UserCookie = r.FormValue("username")
	tlog.Handler.Infof(context.TODO(), logs.DLTagProcessLog, "etype=controller_Login||econtent=%v||err=", data)
	tag, err := logic.Login(r.FormValue("username"), r.FormValue("password"))
	if tag {
		responseInfo.Tag = true
		responseInfo.ErrMsg = "success"
	} else {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err
	}
	logic.EchoJSON(w, r, responseInfo)
}

//退出系统
func (h RPCService) Logout(w http.ResponseWriter, r *http.Request) {
	logOutUrl := constant.LogoutUrl
	hashCache, _ := r.Cookie("__hash__cache")
	referer := r.Referer()
	if hashCache == nil || referer != "http://strategy-arch-platform.intra.xiaojukeji.com/" {
		logOutUrl = constant.OldLogoutUrl
	}
	responseInfo := &idl.ResponseInfo{
		Tag:    false,
		ErrMsg: logOutUrl,
	}
	logic.EchoJSON(w, r, responseInfo)
}

func ReturnOpFailMsg(w http.ResponseWriter, r *http.Request, errMsg string) {
	responseInfo := &idl.ResponseInfo{
		Tag:     false,
		TypeNum: 2, //权限校验不通过
		ErrMsg:  errMsg,
	}
	w.WriteHeader(http.StatusFound)
	logic.EchoJSON(w, r, responseInfo)
}

//如果登录校验通过，则统一返回信息
func ReturnLoginSuccessMsg(w http.ResponseWriter, r *http.Request, content interface{}) {
	responseInfo := &idl.ResponseInfo{
		Tag:     true,
		Content: content,
	}
	logic.EchoJSON(w, r, responseInfo)
}

//
//如果登录校验通过，则统一返回信息
func ReturnLoginSuccessMessage(w http.ResponseWriter, r *http.Request, content interface{}) {

	responseMsg := &idl.ResponseMsg{
		Code:    0,
		Message: "Sucess!",
		Data:    content,
	}

	logic.EchoToJSON(w, r, responseMsg)
}

func ReturnOpFailMessage(w http.ResponseWriter, r *http.Request, errMsg string) {
	responseMsg := &idl.ResponseMsg{
		Code:    1,
		Message: "权限校验不通过", //权限校验不通过
	}
	w.WriteHeader(http.StatusFound)
	logic.EchoToJSON(w, r, responseMsg)
}

//如果登录校验不通过，则统一返回提示信息
func ReturnLoginFailMessage(w http.ResponseWriter, r *http.Request) {
	responseMsg := &idl.ResponseMsg{
		Code:    2,
		Message: "登录校验不通过", //登录校验不通过
	}
	w.WriteHeader(http.StatusFound)
	logic.EchoToJSON(w, r, responseMsg)
}

func ReturnFailMessage(w http.ResponseWriter, r *http.Request, errMsg string) {
	responseMsg := &idl.ResponseInfo{
		Tag:    false,
		ErrMsg: errMsg,
	}
	w.WriteHeader(http.StatusFound)
	logic.EchoToJSON(w, r, responseMsg)
}
