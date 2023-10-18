package model

import (
	"time"
)

type Dimension struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	SceneId int64  `json:"sceneId"`
	//第一层key：维度(如"city");value：值的数组
	ContentMap map[string][]string `json:"contentMap"`
	UpdateTime time.Time           `json:"updateTime"`
}

type DimensionIndex struct {
	//一级key->场景id;二级key->维度类型(如"city");三级key->维度值;value->维度id
	DimensionMap map[int64]map[string]map[string]int64 `json:"dimensionMap"`
	UpdateTime   time.Time                             `json:"updateTime"`
}
