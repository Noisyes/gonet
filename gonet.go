package gonet

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter,response *http.Request)

type Engine struct{
	router map[string]HandlerFunc
}

func New() *Engine{
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

func (engine *Engine) addRoute(method ,pattern string,handler HandlerFunc){
	key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) Get(pattern string, handler HandlerFunc){
	engine.addRoute("GET",pattern,handler)
}

func (engine *Engine) Run(address string) error{
	return http.ListenAndServe(address,engine)
}

func(engine *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	key := r.Method +"-"+r.URL.Path
	if handler , ok := engine.router[key]; ok{
		handler(w,r)
	}else{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found"))
	}
}


