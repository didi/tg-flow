package logic

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/constant"
	"tg-service/idl"
)

//登录校验
func Login(username string, password string) (bool, string) {
	if constant.Env == "test" {
		return true, "login success."
	}
	//将请求密码，转为md5值，进行校验
	//	md5PassWord := MD5(password)
	md5PassWord := ""
	sqlStr := "select user_name,pass_word,role_id from user_info where user_name=? and pass_word = ?"
	tag, err := selectUserInfo(sqlStr, username, md5PassWord)
	return tag, err
}

//数据库查询
func selectUserInfo(sqlStr string, userName string, passWord string) (bool, string) {
	stmt, err := mysql.Handler.Prepare(sqlStr)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagProcessLog, "etype=logic_selectUserInfo||econtent=sql:%v||err=%v", sqlStr, err)
		return false, err.Error()
	}
	defer stmt.Close()

	rows, err := stmt.Query(userName, passWord)
	defer rows.Close()

	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagProcessLog, "etype=logic_SelectUserInfo_Query||econtent=sql:%v||err=%v", sqlStr, err)
		return false, err.Error()
	}

	if rows.Err() != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagProcessLog, "etype=logic_UserInforowsErr||econtent=||err=%v", err)
		return false, err.Error()
	}

	tag := false
	for rows.Next() {
		var userInfo = new(idl.UserInfo)
		err := rows.Scan(&userInfo.UserName, &userInfo.PassWord, &userInfo.RoleId)
		if err != nil {
			tlog.Handler.Errorf(context.TODO(), logs.DLTagProcessLog, "etype=logic_UserInforowsScan||econtent=rows:%v||err=%v", rows, err)
			return false, err.Error()
		}
		tag = true
	}

	msg := "login success."
	if !tag {
		msg = "userName or passWord is error."
	}
	return tag, msg
}
