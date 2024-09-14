/**
function:	分段耗时统计工具
Author:		dayunzhangyunfeng@didiglobal.com
Date:		2021-06-29
*/

package timeutils

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"
)

type TimeCoster struct {
	startTime time.Time
	timeUnit  string
	costMap   *sync.Map
}

const (
	TotalCost          = "totalCost"
	TimeUnitSecond     = "s"
	TimeUnitMillSecond = "ms"
	TimeUnitNanoSecond = "ns"
)

// 老的，即将被下面NewTimeCosterUnit替换
func NewTimeCoster() *TimeCoster {
	t := new(TimeCoster)
	t.timeUnit = TimeUnitMillSecond
	t.costMap = &sync.Map{}
	t.StartCount()
	return t
}

/*
*

	 timeUnit三种取值：
		s :	秒
		ms：毫秒,缺省
		ns: 纳秒
*/
func NewTimeCosterUnit(timeUnit string) *TimeCoster {
	t := new(TimeCoster)
	//不想panc, 兼容一下
	if timeUnit != TimeUnitNanoSecond && timeUnit != TimeUnitMillSecond && timeUnit != TimeUnitSecond {
		//todo add warning
		timeUnit = TimeUnitMillSecond
	}
	t.timeUnit = timeUnit
	t.costMap = &sync.Map{}
	t.StartCount()
	return t
}

func (t *TimeCoster) StartCount() {
	t.startTime = time.Now()
	t.costMap.Store(TotalCost, t.startTime.UnixNano())
}

func (t *TimeCoster) StopCount() {
	st, ok := t.costMap.Load(TotalCost)
	if !ok {
		t.costMap.Delete(TotalCost)
		return
	}

	startTime := st.(int64)
	costTime := time.Now().UnixNano() - startTime
	if t.timeUnit == TimeUnitMillSecond {
		costTime = costTime / 1000000
	} else if t.timeUnit == TimeUnitSecond {
		costTime = costTime / 1000000000
	}
	t.costMap.Store(TotalCost, costTime)
}

func (t *TimeCoster) StartSectionCount(key string) {
	t.costMap.Store(key, time.Now().UnixNano())
}

func (t *TimeCoster) StopSectionCount(key string) {
	st, ok := t.costMap.Load(key)
	if !ok {
		t.costMap.Delete(key)
		return
	}

	startTime := st.(int64)
	costTime := time.Now().UnixNano() - startTime
	if t.timeUnit == TimeUnitMillSecond {
		costTime = costTime / 1000000
	} else if t.timeUnit == TimeUnitSecond {
		costTime = costTime / 1000000000
	}

	t.costMap.Store(key, costTime)
}

func (t *TimeCoster) GetSectionCount(key string) (int64, error) {
	ct, oks := t.costMap.Load(key)
	if !oks {
		return 0, fmt.Errorf("no time count info:%v", key)
	}

	costTime := ct.(int64)
	return costTime, nil
}

func (t *TimeCoster) GetAllCounts() map[string]int64 {
	kv := make(map[string]int64)
	t.costMap.Range(func(k, v interface{}) bool {
		key := fmt.Sprintf("%v", k)
		val := v.(int64)
		kv[key] = val
		return true
	})

	return kv
}

func (t *TimeCoster) ToCountString() string {
	var buffer bytes.Buffer
	t.costMap.Range(func(k, v interface{}) bool {
		key := fmt.Sprintf("%v", k)
		val := fmt.Sprintf("%v", v)
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(val)
		buffer.WriteString("||")
		return true
	})

	strs := buffer.String()
	if strings.HasSuffix(strs, "||") {
		strs = strs[0 : len(strs)-2]
	}

	return strs
}

func (t *TimeCoster) SetSectionCount(key string, cost int64) {
	t.costMap.Store(key, cost)
}

func (t *TimeCoster) GetStartTime() time.Time {
	return t.startTime
}

func (t *TimeCoster) Merge(t1 *TimeCoster) {
	if t1 == nil {
		return
	}

	if t.timeUnit != t1.timeUnit {
		//TODO 不允许，先忽略吧
		return
	}

	t1.costMap.Range(func(k, v interface{}) bool {
		if k != TotalCost {
			t.costMap.Store(k, v)
		}
		return true
	})
}
