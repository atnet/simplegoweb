package web

import (
	"net/http"
	"html/template"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"
	_"fmt"
)

//呈现模板
func RenderTemplate(w http.ResponseWriter, tplPath string,data interface{}) {
	t,err := template.ParseFiles(tplPath)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	t.Execute(w,data)
}

//转换到实体
func ParseFormToEntity(values map[string][]string,instance interface{}){
	refVal := reflect.ValueOf(instance).Elem()
	//类型装换参见：http://www.kankanews.com/ICkengine/archives/19245.shtml
	//for i:=0 ; i< refVal.NumField(); i++ {
	//	prop := refVal.Field(i)
	for k,v := range values{
		field := refVal.FieldByName(k)
		if field.IsValid() {
			//
			//var x interface{} = 1
			//y:= x.(type)
			//
			strVal := v[0]
			
			
			switch field.Type().Kind() {
				case reflect.String:
					field.Set(reflect.ValueOf(strVal))
					break;
				case reflect.Int:
					val,err := strconv.Atoi(strVal)
					if err == nil {
						field.Set(reflect.ValueOf(val))
					}
					break
				case reflect.Bool:
					val := strings.ToLower(strVal) == "true" || strVal =="1"
					field.Set(reflect.ValueOf(val))
					break;
					
				//接口类型
				case reflect.Interface:
					if reflect.TypeOf(time.Now()) == field.Type(){
						t,err := time.Parse("2006-01-02 15:04:05", strVal)
						if err == nil {
							field.Set(reflect.ValueOf(t))
						}
					}
					break;
			}
		}
	}
	//fmt.Println(instance)
}

//操作Json结果
type JsonProcessResult struct{
	Result 	bool 			`json:"result"`
	Code	int				`json:"code"`
	Data	interface{}		`json:"data"`
	Message string 			`json:"message"`
}

//序列化
func (r JsonProcessResult) Marshal()string{
	json,_ := json.Marshal(r)
	return string(json)
}
