package model

import (
	"time"
)

type DowngradeStrategyInfo struct {
	LevelList  []*DowngradeLevel `json:"levels"`
	SceneIds   []int64           `json:"scene_ids"`
	UpdateTime time.Time
}

type DowngradeLevel struct {
	Level  int   `json:"level"`
	Rate   int64 `json:"rate"`
	Status bool  `json:"status"`
}
