package gonet

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct{
	router *router
}

func New() *Engine{
	return &Engine{
		router: newRouter(),
	}
}

func (engine *Engine) addRoute(method ,pattern string,handler HandlerFunc){
	engine.router.addRouter(method,pattern,handler)
}

func (engine *Engine) Get(pattern string, handler HandlerFunc){
	engine.addRoute("GET",pattern,handler)
}

func (engine *Engine) POST(pattern string,handler HandlerFunc){
	engine.addRoute("POST",pattern,handler)
}

func (engine *Engine) Run(address string) error{
	return http.ListenAndServe(address,engine)
}

func(engine *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	cxt := newContext(w,r)
	engine.router.handle(cxt)
}


