package gonet

import (
	"log"
	"net/http"
	"strings"
)

type router struct{
	handlers map[string] HandlerFunc
	roots map[string]*node
}

func newRouter() *router{
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots: make(map[string]*node),
	}
}

func parsePattern(pattern string) []string{
	vs := strings.Split(pattern,"/")
	parts := make([]string,0)
	for _,item := range vs{
		if item != ""{
			parts = append(parts,item)
			if item[0] == '*'{
				break
			}
		}
	}
	log.Println(parts)
	return parts
}



func (r *router) addRouter(method,pattern string,handler HandlerFunc){
	log.Printf("Route %4s - %s\n",method,pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_,ok := r.roots[method]
	if !ok{
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern,parts,0)
	r.handlers[key] = handler
}

func(r *router) getRoute(method,path string)(*node,map[string]string){
	searchPath := parsePattern(path)
	params := make(map[string]string)
	root,ok := r.roots[method]
	if !ok{
		return nil,nil
	}
	n := root.search(searchPath,0)
	if n != nil{
		parts := parsePattern(n.pattern)
		for index,part := range parts{
			if part[0]==':'{
				params[part[1:]] = searchPath[index]
			}
			if part[0]=='*'&&len(part)>1{
				params[part[1:]] = strings.Join(searchPath[index:],"/")
				break
			}
		}
		return n,params
	}
	return nil,nil
}

func (r *router) handle(c *Context){
	n,params := r.getRoute(c.Method,c.Path)
	if n != nil{
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers,r.handlers[key])
	}else{
		c.handlers = append(c.handlers,func(c *Context){
			c.String(http.StatusNotFound,"404 NOT FOUND : %s\n",c.Path)
		})
	}
	c.Next()
}
