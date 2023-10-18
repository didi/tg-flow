/**
	Description:	存放apollo实验相关的配置信息,提供相关访问方法的封装
	Author:			dayunzhangyunfeng@didiglobal.com
	Date:			2020-08-31
**/

package model

import (
	"fmt"
	"go.intra.xiaojukeji.com/apollo/apollo-golang-sdk-v2"
	apolloModel "go.intra.xiaojukeji.com/apollo/apollo-golang-sdk-v2/model"
)

type ApolloConfig struct {
	ApolloUsers            map[string]*apolloModel.User
	dispatchUser           *apolloModel.User
	dispatchExperimentName string
}

func NewApolloConfig() *ApolloConfig {
	apolloUsers := make(map[string]*apolloModel.User)
	return &ApolloConfig{ApolloUsers: apolloUsers}
}

// SetDispatchUser 设置分流 User
func (this *ApolloConfig) SetDispatchUser(apolloUser *apolloModel.User) {
	this.dispatchUser = apolloUser
}

// SetDispatchExperimentName 设置分流实验名，如果分流实验名保存在数据库中，则不需要 SetDispatchExperimentName
func (this *ApolloConfig) SetDispatchExperimentName(experimentName string) {
	this.dispatchExperimentName = experimentName
}

func (this *ApolloConfig) GetDispatchExperimentName() string {
	return this.dispatchExperimentName
}

func (this *ApolloConfig) GetDispatchGroupName() (string, error) {

	if this.dispatchUser == nil || this.dispatchExperimentName == "" {
		return "", fmt.Errorf("dispatch info not correct, user=%v, exp_name=%v", this.dispatchUser, this.dispatchExperimentName)
	}

	toggle, err := apollo.FeatureToggle(this.dispatchExperimentName, this.dispatchUser)
	if err != nil {
		return "", err
	}

	if toggle.IsAllow() {
		return toggle.GetAssignment().GetGroupName(), nil
	} else {
		//The individual which is not allowed to enter the experiment using the default strategy
		return "", fmt.Errorf("the user is not permitted to enter the apollo experiment: %v", this.dispatchExperimentName)
	}
}

// SetApolloUser 已经弃用，建议自己采用 Apollo SDK 执行实验相关操作
func (this *ApolloConfig) SetApolloUser(experimentName string, apolloUser *apolloModel.User) {
	this.ApolloUsers[experimentName] = apolloUser
}

// GetGroupName 已经弃用，建议自己采用 Apollo SDK 执行实验相关操作
func (this *ApolloConfig) GetGroupName(experimentName string) (string, error) {
	toggleResult, err := apollo.FeatureToggle(experimentName, this.ApolloUsers[experimentName])
	if err != nil {
		return "", err
	}

	//If toggle.allow is true, this individual is allowed to enter the experiment
	if toggleResult.IsAllow() {
		assignment := toggleResult.GetAssignment()
		return assignment.GetGroupName(), nil
	} else {
		//The individual which is not allowed to enter the experiment using the default strategy
		return "", fmt.Errorf("the user is not permitted to enter the apollo experiment: %v", experimentName)
	}
}

// GetApolloParam 取指定实验的指定参数值，注意：此函数有可能由于缺省值的补漏而隐藏取参数值失败的错误,
// 已经弃用，建议自己采用 Apollo SDK 执行实验相关操作
func (this *ApolloConfig) GetApolloParam(experimentName, paramName, defaultParamValue string) (string, error) {
	apolloUser, ok := this.ApolloUsers[experimentName]
	if !ok {
		return defaultParamValue, fmt.Errorf("apollo user for experimentName:%v not initialzed!", experimentName)
	}

	toggleResult, err := apollo.FeatureToggle(experimentName, apolloUser)
	if err != nil {
		return defaultParamValue, err
	}

	//If toggle.allow is true, this individual is allowed to enter the experiment
	if toggleResult.IsAllow() {
		//The sample were divided into different groups and used different strategies.
		assignment := toggleResult.GetAssignment()
		return assignment.GetParameter(paramName, defaultParamValue), nil
	} else {
		return defaultParamValue, fmt.Errorf("user not allowed in apollo experiment:%v", experimentName)
	}
}

// GetApolloParams 取指定实验的全部参数并以map[string]string格式返回，如value不为string类型，则转为string
// 已经弃用，建议自己采用 Apollo SDK 执行实验相关操作
func (this *ApolloConfig) GetApolloParams(experimentName string) (map[string]string, error) {
	apolloUser, ok := this.ApolloUsers[experimentName]
	if !ok {
		return nil, fmt.Errorf("apollo user for experimentName:%v not initialzed!", experimentName)
	}

	toggleResult, err := apollo.FeatureToggle(experimentName, apolloUser)
	if err != nil {
		return nil, err
	}

	//If toggle.allow is true, this individual is allowed to enter the experiment
	if toggleResult.IsAllow() {
		//The sample were divided into different groups and used different strategies.
		assignment := toggleResult.GetAssignment()
		return assignment.GetParameters(), nil
	} else {
		return nil, fmt.Errorf("user not allowed in apollo experiment:%v", experimentName)
	}
}

// GetRawApolloParams 取指定实验的全部原始参数
// 已经弃用，建议自己采用 Apollo SDK 执行实验相关操作
func (this *ApolloConfig) GetRawApolloParams(experimentName string) (map[string]interface{}, error) {
	apolloUser, ok := this.ApolloUsers[experimentName]
	if !ok {
		return nil, fmt.Errorf("apollo user for experimentName:%v not initialzed!", experimentName)
	}

	toggleResult, err := apollo.FeatureToggle(experimentName, apolloUser)
	if err != nil {
		return nil, err
	}

	//If toggle.allow is true, this individual is allowed to enter the experiment
	if toggleResult.IsAllow() {
		//The sample were divided into different groups and used different strategies.
		assignment := toggleResult.GetAssignment()
		return assignment.GetRawParameters(), nil
	} else {
		return nil, fmt.Errorf("user not allowed in apollo experiment:%v", experimentName)
	}
}
