package web

import (
	"net/http"
	"regexp"
	"reflect"
	"strings"
	_"fmt"
)


//控制器处理
//@controller ： 包含多种动作，URL中的文件名自动映射到控制器的函数
//				 注意，是区分大小写的,默认映射到index函数
//				 如果是POST请求将映射到控制器“函数名+_post”的函数执行
func ProcessByController (
	controller interface{} ,
	 w http.ResponseWriter,
	  r *http.Request){
	  	
	var do string
	reg := regexp.MustCompile("/([^/]+)$")
	groups := reg.FindAllStringSubmatch(r.URL.Path,1)
	if len(groups) == 0 || len(groups[0]) == 0 {
		do = "Index"
	}else{
		do = groups[0][1]
		
		//去扩展名
		extIndex := strings.Index(do,".")
		if extIndex != -1 {
			do = do[0:extIndex]
		}
		
		//将第一个字符转为大写,这样才可以
		upperFirstLetter := strings.ToUpper(do[0:1])
		if upperFirstLetter != do[0:1] {
			do = upperFirstLetter+do[1:]
		}
	}
	
	if r.Method == "POST" {
		do += "_post"	
	}
	
	t := reflect.ValueOf(controller)
	params := []reflect.Value{reflect.ValueOf(w),reflect.ValueOf(r) }
	method := t.MethodByName(do)
	
	if !method.IsValid(){
		//fmt.Println(do)
		w.Write([]byte("No action Named:"+ strings.Replace(do,"_post","",1)))
	}else{
		method.Call(params)
	}
}

