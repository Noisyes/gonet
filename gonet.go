package gonet

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct{
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct{
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup
	engine *Engine
}

func New() *Engine{
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups,newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method ,comp string,handler HandlerFunc){
	pattern := g.prefix + comp
	log.Printf("Route %4s = %s\n",method,pattern)
	g.engine.addRoute(method,pattern,handler)
}

func (g *RouterGroup) GET(pattern string,handler HandlerFunc){
	g.addRoute("GET",pattern,handler)
}

func(g *RouterGroup) POST(pattern string,handler HandlerFunc){
	g.addRoute("POST",pattern,handler)
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

func(group *RouterGroup) Use(middlewares ... HandlerFunc){
	group.middlewares = append(group.middlewares,middlewares...)
}

func(engine *Engine) ServeHTTP(w http.ResponseWriter,r *http.Request){
	var middlewares []HandlerFunc
	for _,group := range engine.groups{
		if strings.HasPrefix(r.URL.Path,group.prefix){
			middlewares = append(middlewares,group.middlewares...)
		}
	}
	c := newContext(w,r)
	c.handlers = middlewares
	engine.router.handle(c)
}


