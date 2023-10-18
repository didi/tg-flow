package model

import (

)

type RecallConfig struct {
	Id			int64				`json:"id"`
	SceneId		int64				`json:"scene_id"`
	RecallRule	[]RecallRuleType	`json:"recall_rule"`
}

type RecallRuleType struct {
	RecallType	string				`json:"recall_type"`
	TopN		int64				`json:"top_n"`
}