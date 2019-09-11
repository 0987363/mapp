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

func (c *Context)run() {
	defer func() {
		if err := recover(); err != nil {
			log.Print("recv stack: ", err, string(debug.Stack()))
		}
	}()

	c.Groups[0](c)
}

func (c *Context)Clone() *Context{
	return &Context {
		Groups : c.Groups,
	}
}

func (c *Context)Run() {
	for {
		c.Clone().run()
		time.Sleep(time.Second * 5)
	}
}

func NewContext(groups HandlersChain) *Context {
	return &Context {
		Groups : groups,
	}
}

func (engine *Engine)Run() {
	for _, h := range engine.Handlers {
		c := NewContext(append(engine.Groups, h))
		go c.Run()
	}
	select{}
}

func (engine *Engine)Use(middleware ...HandlerFunc) {
	engine.Groups = append(engine.Groups, middleware...)
}

func (engine *Engine)Register(handlers ...HandlerFunc) {
	engine.Handlers = append(engine.Handlers, handlers...)
}

