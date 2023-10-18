package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"time"
	"strings"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

/**
	生成32位MD5
**/
func MD5(text string) string {
   ctx := md5.New()
   ctx.Write([]byte(text))
   return hex.EncodeToString(ctx.Sum(nil))
}

func GetCurrentTimeString() string {
	return time.Now().Format(timeLayout)
}

/**
	将原始的time格式转换成标准格式，如：
	原始格式：2020-08-17T14:36:31+08:00  ===> 转换后的标准格式 ：2020-08-17 14:36:31
**/
func FormatDateTime(updateTime string) string {
	length := len(updateTime)
	if length > 19 {
		updateTime =  updateTime[:19]
	}
	
	updateTime = strings.ReplaceAll(updateTime, "T", " ")
	
	return updateTime
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func SubStr(str string, start int, end int) (string, error) {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		err := errors.New("SubStr start is wrong")
		return "", err
	}
	if end < 0 || end > length {
		err := errors.New("SubStr end is wrong.")
		return "", err
	}
	return string(rs[start:end]), nil
}

//将int64数组，转化为string字符串
func Int64ArrayTOString(array []int64) string {
	resultStr := ""
	if array == nil || len(array) <= 0 {
		return resultStr
	}
	for _, data := range array {
		if len(resultStr) <= 0 {
			resultStr += strconv.FormatInt(data, 10)
		} else {
			resultStr += "," + strconv.FormatInt(data, 10)
		}
	}
	return resultStr
}

//map[string]string,转化为string
func StringMapToString(strMap map[string]string) string {
	resultStr := ""
	if strMap == nil || len(strMap) <= 0 {
		return resultStr
	}
	for key, value := range strMap {
		if len(resultStr) <= 0 {
			resultStr += key + ":" + value
		} else {
			resultStr += "," + key + ":" + value
		}
	}
	return resultStr
}

//map[int64]int64,转化为string
func Int64MapToString(int64Map map[int64]int64) string {
	resultStr := ""
	if int64Map == nil || len(int64Map) <= 0 {
		return resultStr
	}
	for key, value := range int64Map {
		if len(resultStr) <= 0 {
			resultStr += strconv.FormatInt(key, 10) + ":" + strconv.FormatInt(value, 10)
		} else {
			resultStr += "," + strconv.FormatInt(key, 10) + ":" + strconv.FormatInt(value, 10)
		}
	}
	return resultStr
}

//map[int64]int,转化为string
func IntMapToString(intMap map[int64]int) string {
	resultStr := ""
	if intMap == nil || len(intMap) <= 0 {
		return resultStr
	}
	for key, value := range intMap {
		if len(resultStr) <= 0 {
			resultStr += strconv.FormatInt(key, 10) + ":" + strconv.Itoa(value)
		} else {
			resultStr += "," + strconv.FormatInt(key, 10) + ":" + strconv.Itoa(value)
		}
	}
	return resultStr
}

//map[int64]map[string]float64,转化为string
func Float64MapToString(float64Map map[int64]map[string]float64) string {
	resultStr := ""
	if float64Map == nil || len(float64Map) <= 0 {
		return resultStr
	}
	for key, hourMap := range float64Map {
		hourStr := ""
		if hourMap == nil || len(hourMap) <= 0 {
			hourStr += "nil"
		} else {
			for hour, value := range hourMap {
				if len(hourStr) <= 0 {
					hourStr += hour + ":" + strconv.FormatFloat(value, 'f', -1, 64)
				} else {
					hourStr += "," + hour + ":" + strconv.FormatFloat(value, 'f', -1, 64)
				}
			}
		}
		hourStr = strconv.FormatInt(key, 10) + ":{" + hourStr + "}"
		if len(resultStr) <= 0 {
			resultStr += hourStr
		} else {
			resultStr += "," + hourStr
		}
	}
	return resultStr
}

func ToString(itr interface{}) string {
	b, err := json.Marshal(itr)
	if err != nil {
		return fmt.Sprintf("%+v", itr)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", itr)
	}
	return out.String()
}
