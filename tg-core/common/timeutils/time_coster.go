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
	startTime	time.Time
	costMap		*sync.Map
}

const (
	totalCost = "totalCost"
)

func NewTimeCoster() *TimeCoster {
	t := new(TimeCoster)
	t.costMap	= &sync.Map{}
	t.StartCount()
	return t
}

func (t *TimeCoster) StartCount() {
	t.startTime = time.Now()
	t.costMap.Store(totalCost, t.startTime.UnixNano())
}

func (t *TimeCoster) StopCount() {
	st, ok := t.costMap.Load(totalCost)
	if !ok {
		t.costMap.Delete(totalCost)
		return
	}

	startTime := st.(int64)
	costTime := (time.Now().UnixNano() - startTime)/1000000
	t.costMap.Store(totalCost, costTime)
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
	costTime := (time.Now().UnixNano() - startTime)/1000000
	t.costMap.Store(key, costTime)
}

func (t *TimeCoster) GetSectionCount(key string) (int64, error) {
	ct,oks := t.costMap.Load(key)
	if !oks {
		return 0, fmt.Errorf("no time count info:%v", key)
	}

	costTime:=ct.(int64)
	return costTime, nil
}

func (t *TimeCoster) GetAllCounts() (map[string]int64) {
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
	if strings.HasSuffix(strs,"||") {
		strs = strs[0:len(strs)-2]
	}

	return strs
}

func (t *TimeCoster) SetSectionCount(key string, cost int64){
	t.costMap.Store(key, cost)
}

func (t *TimeCoster) GetStartTime() time.Time {
	return t.startTime
}

//目前没用，就不要了
/*func (t *TimeCoster) Merge(t1 *TimeCoster){
	if t1 == nil {
		return
	}

	t1.costMap.Range(func(k, v interface{}) bool {
		t.costMap.Store(k, v)
		return true
	})
}*/
