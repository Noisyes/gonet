package gonet

import (
	"log"
	"time"
)

func Logger()HandlerFunc{
	return func(c *Context){
		t := time.Now()
		log.Printf("[%d] %s in %v",c.StatusCode,c.Req.RequestURI,time.Since(t))
	}
}
