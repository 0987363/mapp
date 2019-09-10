package mapp

import (
	"runtime/debug"
	"log"
	"time"
	"reflect"
	"runtime"
	"strings"
)

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Context struct {
	Groups HandlersChain
	index    int8

	Keys     map[string]interface{}
}

type Engine struct {
	Groups HandlersChain
	Handlers HandlersChain
}

func New() *Engine {
	return &Engine{}
}

func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *Context) Method() (string) {
	f := c.Groups[len(c.Groups) - 1]
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	ss := strings.Split(name, ".")
	return ss[len(ss) - 1]
}

func (c *Context) Next() {
	c.index++
	s := int8(len(c.Groups))
	for ; c.index < s; c.index++ {
		c.Groups[c.index](c)
	}
}

func (c *Context)Main(groups HandlersChain) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("recv stack: ", err, string(debug.Stack()))
		}
	}()

	c.index = 0
	c.Groups = groups

	c.Groups[0](c)
}

func (c *Context)Run(groups HandlersChain) {
	for {
		c.Main(groups)

		time.Sleep(time.Second * 5)
	}
}

func (engine *Engine)Run() {
	for _, h := range engine.Handlers {
		var c Context
		go c.Run(append(engine.Groups, h))
	}
	select{}
}

func (engine *Engine)Use(middleware ...HandlerFunc) {
	engine.Groups = append(engine.Groups, middleware...)
}

func (engine *Engine)Register(handlers ...HandlerFunc) {
	engine.Handlers = append(engine.Handlers, handlers...)
}

