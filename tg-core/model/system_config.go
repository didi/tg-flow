package model

import (
	"time"
)

type SystemConfigInfo struct {
	Id         int64     `json:"id"`
	Type       int       `json:"type"`
	KeyName    string    `json:"keyname"`
	Content    string    `json:"content"`
	UpdateTime time.Time `json:"updateTime"`
}
