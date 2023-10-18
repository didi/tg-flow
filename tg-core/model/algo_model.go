package model

import (
	"time"
)

//算法模型配置索引
type AlgoModelIndex struct {
	IndexVersion string              `json:"index_version"`
	IndexMap     map[string][]string `json:"index_map"`
	UpdateTime   time.Time           `json:"updateTime"`
}

//索引更新后，向center回报状态的返回结果:ErrNo=0,表示操作成功
type CenterServerReportResponseInfo struct {
	ErrNo  int32  `thrift:"err_no,1,required" json:"err_no"`
	ErrMsg string `thrift:"err_msg,2,required" json:"err_msg"`
}
