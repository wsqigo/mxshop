package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"time"
)

type structFieldInfo struct {
	jsonTag    string
	typeKind   reflect.Kind
	typeString string
}

func getStructFieldsInfo(rv reflect.Type) []structFieldInfo {
	info := make([]structFieldInfo, 0)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// 暂时只处理一层struct的简单情况
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if f.Anonymous {
			if ft.Kind() == reflect.Struct {
				subInfo := getStructFieldsInfo(ft)
				info = append(info, subInfo...)
				continue
			}
		}

		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = f.Name
		}

		info = append(info, structFieldInfo{
			jsonTag:    jsonTag,
			typeKind:   ft.Kind(),
			typeString: ft.String(),
		})
	}

	return info
}

// 思路：
// 1. 把queryParam中结构体v有使用的字段也加入data，支持Unmarshal到结构体
// 2. 对时间类型：尝试把纳秒数转换为时间
func PrettyUnmarshal(data []byte, queryParam map[string]string, v any) error {
	rv := reflect.TypeOf(v)
	fieldInfos := getStructFieldsInfo(rv)
	timeFields := make([]string, 0)
	var err error

	for _, fInfo := range fieldInfos {
		if fInfo.typeString == "time.Time" {
			timeFields = append(timeFields, fInfo.jsonTag)
		}

		if queryParam == nil {
			continue
		}

		if param, exist := queryParam[fInfo.jsonTag]; exist {
			// 检查请求体 data 中是否存在
			_, dataType, _, _ := jsonparser.Get(data, fInfo.jsonTag)
			fmt.Println(dataType.String())
			if dataType != jsonparser.NotExist {
				return errors.New("参数" + fInfo.jsonTag + "重复")
			}

			setValue := "\"" + param + "\""
			// 数值类型的不能加双引号
			if fInfo.typeKind >= reflect.Bool && fInfo.typeKind <= reflect.Float64 {
				setValue = param
			} else if fInfo.typeString == "time.Time" {
				// 对于时间类型的可能传递的是纳秒数值，也需要处理
				if _, err = strconv.ParseInt(param, 10, 64); err == nil {
					setValue = param
				}
			}

			data, err = jsonparser.Set(data, []byte(setValue), fInfo.jsonTag)
			if err != nil {
				return errors.New("添加参数" + fInfo.jsonTag + "失败")
			}
		}
	}

	if len(timeFields) > 0 {
		cJson := map[string]any{}
		err = json.Unmarshal(data, &cJson)
		if err != nil {
			return err
		}

		for _, tf := range timeFields {
			sv, ok := cJson[tf]
			if !ok {
				continue
			}

			if fv, ok := sv.(float64); ok {
				iv := int64(fv)
				t := time.Unix(0, iv)
				if t.Year() < 1999 || t.Year() > 2050 {
					return errors.New("date out range")
				}
			}
		}
	}

	// 支持把纳秒时间数转化为时间、支持把string转化为float
	return jsoniter.Unmarshal(data, v)
}

func ShouldQueryParam(ctx *gin.Context, obj any) error {
	bt := []byte("{}")
	requestQuery := map[string]string{}
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 1 {
			return errors.New("request query length larger than 1")
		}
		requestQuery[k] = v[0]
	}

	return PrettyUnmarshal(bt, requestQuery, obj)
}
