package web

import (
	"net/http"
	"regexp"
)

//路由映射
type RouteMap struct {
	RouteCollection map[string]func(http.ResponseWriter, *http.Request)
}

//添加路由
func (this *RouteMap) Add(
	urlPattern string,
	requestFunc func(http.ResponseWriter, *http.Request)) {

	if this.RouteCollection == nil {
		this.RouteCollection =
			make(map[string]func(http.ResponseWriter, *http.Request))
	}
	_, exists := this.RouteCollection[urlPattern]

	if !exists {
		this.RouteCollection[urlPattern] = requestFunc
	}
}
	
//处理请求 
func (this *RouteMap) HandleRequest(w http.ResponseWriter, r *http.Request) {
	handleMapRoute(w, r, this.RouteCollection)
}

//处理路由请求
func handleMapRoute(
	w http.ResponseWriter,
	r *http.Request,
	routes map[string]func(http.ResponseWriter, *http.Request)) {
	path := r.URL.Path
	var isHandled bool = false

	for k, v := range routes {
		matched, err := regexp.Match(k, []byte(path))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if matched && v != nil {
			isHandled = true
			v(w, r)
			break
		}
	}

	if !isHandled {
		w.Write([]byte("404 Not found!"))
	}
}
 
