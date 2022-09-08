package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type JsonToStruct struct {
}

//不支持空数组(array [])
//不支持结构体中存在不能导出的key  struct{name string `json:name`}
//不支持 format-data 传关联数组 （a[b][]=[1,2,3],a[c][]=[4,5,6]）
func (j *JsonToStruct) ParseJson(body []byte, v interface{}) error {
	var err error

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}

	tmpMap := map[string]interface{}{}
	err = json.Unmarshal(body, &tmpMap)
	if err != nil {
		return errors.New(fmt.Sprintf("json Unmarshal error err:%s,string:%s", err.Error(), string(body)))
	}
	rv = rv.Elem()
	log.Printf("json to struct map:%+v, kind:%+v", tmpMap, rv.Kind())
	if rv.Kind() == reflect.Struct {
		//struct
		return j.jsonToStruct(rv, tmpMap)
	} else {
		return errors.New(fmt.Sprintf("Invalid format this format:%+v", rv.Kind()))
	}

	return err
}

func (j *JsonToStruct) jsonToStruct(rv reflect.Value, mapInterface map[string]interface{}) error {
	var err error
	if rv.Kind() != reflect.Struct {
		return errors.New("must be pointer of struct")
	}

	numField := rv.Type().NumField()
	for n := 0; n < numField; n++ {
		structField := rv.Type().Field(n)
		//如果是 匿名结构体
		if structField.Anonymous == true {
			err = j.jsonToStruct(rv.Field(n), mapInterface)
			if err != nil {
				log.Println(err.Error())
				break
			}
			continue
		}

		//非匿名结构体
		fieldTag := j.getJsonTagName(structField.Tag.Get("json"))
		val, ok := mapInterface[fieldTag]
		if ok == false {
			continue
		}

		switch structField.Type.Kind() {
		case reflect.Int64, reflect.Int, reflect.Int32, reflect.Int8, reflect.Int16:
			tmpV, err := j.toInt64(val)
			if err != nil {
				log.Println(err.Error())
				break
			}
			rv.Field(n).SetInt(tmpV)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tmpV, err := j.toUInt64(val)
			if err != nil {
				log.Println(err.Error())
				break
			}
			rv.Field(n).SetUint(tmpV)
		case reflect.Float64, reflect.Float32:
			tmpV, err := j.toFloat64(val)
			if err != nil {
				log.Println(err.Error())
				break
			}
			rv.Field(n).SetFloat(tmpV)
		case reflect.String:
			tmpV, err := j.toString(val)
			if err != nil {
				log.Println(err.Error())
				break
			}
			rv.Field(n).SetString(tmpV)
		case reflect.Struct:
			err = j.jsonToStruct(rv.Field(n), val.(map[string]interface{}))
			if err != nil {
				log.Println(err.Error())
				break
			}
		case reflect.Slice:
			sliceType := rv.Field(n).Type().String()
			if sliceType == "[]uint64" {
				tmpSlice := make([]uint64, 0, len(val.([]interface{})))
				for _, v := range val.([]interface{}) {
					tmpV, err := j.toUInt64(v)
					if err != nil {
						log.Println(err.Error())
						break
					}
					tmpSlice = append(tmpSlice, tmpV)
				}
				rv.Field(n).Set(reflect.ValueOf(tmpSlice))
			} else if sliceType == "[]string" {
				tmpSlice := make([]string, 0, len(val.([]interface{})))
				for _, v := range val.([]interface{}) {
					tmpV, err := j.toString(v)
					if err != nil {
						log.Println(err.Error())
						break
					}
					tmpSlice = append(tmpSlice, tmpV)
				}
				rv.Field(n).Set(reflect.ValueOf(tmpSlice))
			}
		case reflect.Interface:
			rv.Field(n).Set(reflect.ValueOf(val))
		case reflect.Bool:
			rv.Field(n).Set(reflect.ValueOf(val))
		case reflect.Map:
			mapType := rv.Field(n).Type().String()
			if mapType == "map[string]string" {
				tmpStringMap := make(map[string]string, 0)
				for k, v := range val.(map[string]interface{}) {
					tmpString, err := j.toString(v)
					if err != nil {
						log.Println(err.Error())
						continue
					}
					tmpStringMap[k] = tmpString
				}
				rv.Field(n).Set(reflect.ValueOf(tmpStringMap))
			} else if mapType == "map[string]interface {}" {
				rv.Field(n).Set(reflect.ValueOf(val.(map[string]interface{})))
			}
		default:
			break
		}
	}
	return err
}

func (j *JsonToStruct) toInt64(v interface{}) (int64, error) {
	returnVal := int64(0)
	switch v.(type) {
	case string:
		tmpFV, err := strconv.ParseFloat(v.(string), 10)
		if err != nil {
			break
		}
		returnVal = int64(tmpFV)
	case float64:
		returnVal = int64(v.(float64))
	default:
		return returnVal, errors.New(fmt.Sprintf("type error type:%+v", reflect.TypeOf(v)))
	}
	return returnVal, nil
}

func (j *JsonToStruct) toUInt64(v interface{}) (uint64, error) {
	returnVal := uint64(0)
	switch v.(type) {
	case string:
		tmpFV, err := strconv.ParseUint(v.(string), 10, 64)
		if err != nil {
			break
		}
		returnVal = uint64(tmpFV)
	case float64:
		returnVal = uint64(v.(float64))
	default:
		return returnVal, errors.New(fmt.Sprintf("type error type:%+v", reflect.TypeOf(v)))
	}
	return returnVal, nil
}

func (j *JsonToStruct) toFloat64(v interface{}) (float64, error) {
	returnVal := float64(0)
	switch v.(type) {
	case string:
		val, err := strconv.ParseFloat(v.(string), 10)
		if err != nil {
			return returnVal, err
		}
		returnVal = val
	default:
		return returnVal, errors.New(fmt.Sprintf("type error type:%+v", reflect.TypeOf(v)))
	}
	return returnVal, nil
}

func (j *JsonToStruct) toString(v interface{}) (string, error) {
	returnVal := ""
	switch v.(type) {
	case string:
		returnVal = v.(string)
	case float64:
		returnVal = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	default:
		return returnVal, errors.New(fmt.Sprintf("type error type:%+v", reflect.TypeOf(v)))
	}
	return returnVal, nil
}

func (j *JsonToStruct) getJsonTagName(tagNameStr string) string {
	tagName := strings.Split(tagNameStr, ",")
	return tagName[0]
}
