package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"reflect"
	"sort"
	"strconv"
)

/**
 * 开始计时
 * return time
 */
func Start() time.Time {
	t := time.Now()
	return t
}

/**
 * 结束计时
 * return ms  转化成毫秒 保留三位小数
 */
func End(t time.Time) float64 {
	return time.Since(t).Seconds() * 1000
	//return fmt.Sprintf("%.3f", s)
}

//map to key=val&key=val 按key正序排列
func MapSortKeyToString(data map[string]interface{}) string {
	keySlice := make([]string, 0)
	dataStr := ""
	if len(data) <= 0 {
		return ""
	}
	for k, _ := range data {
		keySlice = append(keySlice, k)
	}
	sort.Strings(keySlice)
	for _, k := range keySlice {
		v := data[k]
		tmpS := ""
		typeOf := reflect.TypeOf(v)
		if typeOf != nil {
			switch typeOf.String() {
			case "string":
				tmpS = v.(string)
			case "int64":
				tmpS = strconv.FormatInt(v.(int64), 10)
			case "uint64":
				tmpS = strconv.FormatUint(v.(uint64), 10)
			case "float64":
				tmpS = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			case "uint":
				tmpS = strconv.FormatUint(uint64(v.(uint)), 10)
			case "int":
				tmpS = strconv.Itoa(v.(int))
			default:
				tmpS = ""
			}
		} else {
			tmpS = "null"
		}
		dataStr += k + "=" + tmpS + "&"
	}
	dataStr = dataStr[:len(dataStr)-1]

	return dataStr
}

func UniqueId() string {
	return uuid.New().String()
}

//map to key=val&key=val 按key正序排列
func MapSortKeyToJsonString(data map[string]interface{}) string {
	keyMap := make([]string, 0)
	dataStr := ""
	if len(data) <= 0 {
		return ""
	}
	for k, _ := range data {
		keyMap = append(keyMap, k)
	}
	sort.Strings(keyMap)
	for _, k := range keyMap {
		v := data[k]
		tmpB, _ := json.Marshal(v)
		tmpS := string(tmpB)
		dataStr += k + "=" + tmpS + "&"
	}
	dataStr = dataStr[:len(dataStr)-1]
	fmt.Println("MapSortKeyToJsonString", dataStr)
	return dataStr
}
