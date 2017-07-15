package core

import (
    "fmt"
    "github.com/pkg/errors"
    "reflect"
)

type Scope uint16

const (
    Prototype Scope = iota
    Singleton
)

type Initialization bool

const (
    Lazy = true
)

type Service interface {
}

type Provider interface {
    Name() string
    
    Get() interface{}
    
    Type() reflect.Type
}

type ApplicationContext interface {
    RegisterSingleton(interface{})
    
    RegisterProvider(*Provider)
    
    ContainsNamed(string)
    
    Contains(interface{}) bool
    
    GetByName(string) interface{}
    
    GetByType(reflect.Type) interface{}
    
    Scan(...interface{})
    
    /**
    
     */
    ByName(string) (interface{}, error)
    
    ByType(reflect.Type) (interface{}, error)
}

type DefaultApplicationContext struct {
    ApplicationContext
    
    objectGraph        *Graph
    namedSingletons    map[string]interface{}
    typedSingletons    map[reflect.Type]interface{}
    typeMappings       map[reflect.Type]interface{}
    namedInstantiators map[string]*Instantiator
}

func (c *DefaultApplicationContext) RegisterSingleton(instance interface{}) {
    t := reflect.TypeOf(instance)
    c.typeMappings[t] = instance
    c.typedSingletons[t] = instance
}

func (c *DefaultApplicationContext) Contains(i interface{}) bool {
    t := reflect.TypeOf(i)
    return c.typeMappings[t] != nil
}

func (c *DefaultApplicationContext) GetByType(t reflect.Type) interface{} {
    return c.typeMappings[t]
}

func (c *DefaultApplicationContext) registerSelf() {
    var cif *ApplicationContext
    ctype := reflect.TypeOf(cif).Elem()
    n := &Node{
        Key:     ctype,
        NameKey: "ApplicationContext",
        Value: &Instantiator{
            instance: c,
        },
    }
    c.objectGraph.Add(n, true)
}

func (c *DefaultApplicationContext) Scan(ifaces ...interface{}) {
    c.registerSelf()
    for _, i := range ifaces {
        log.Debugf("Scanning configuration %s", i)
        t := reflect.TypeOf(i)
        for m := 0; m < t.NumMethod(); m++ {
            method := t.Method(m)
            c.registerMethod(&method, i)
        }
        
    }
    for _, i := range ifaces {
        t := reflect.TypeOf(i)
        for m := 0; m < t.NumMethod(); m++ {
            method := t.Method(m)
            c.registerDependencies(&method, i)
        }
    }
    
    for _, v := range c.objectGraph.eager {
        v.Value.Construct()
    }
}

func (c *DefaultApplicationContext) ByName(name string) (interface{}, error) {
    if r := c.GetByName(name); r != nil {
        return r, nil
    }
    return nil, errors.New(fmt.Sprintf("No bean named '%s' is available in this application context ", name));
}

func (c *DefaultApplicationContext) GetByName(name string) interface{} {
    value := c.namedSingletons[name]
    if value == nil {
        result, exists := c.objectGraph.NamedNodes[name]
        if exists {
            return result.Value.Construct()
        }
        return nil
    }
    return value
}

func (c *DefaultApplicationContext) Get(t reflect.Type) interface{} {
    value := c.typedSingletons[t]
    if value == nil {
        return c.objectGraph.Nodes[t].Value.Construct()
    }
    return value
}

func (c *DefaultApplicationContext) registerDependencies(m *reflect.Method, host interface{}) {
    
    var current *Instantiator
    cnode := c.objectGraph.GetByName(m.Name)
    if cnode == nil && m.Type.NumOut() > 0 {
        cnode = c.objectGraph.Get(m.Type.Out(0))
        log.Debugf("Resolving dependencies for type %s", m.Type.Out(0).Name())
    } else {
        log.Debugf("Resolving dependencies for named object %s", m.Name)
    }
    
    if cnode == nil {
        //todo better error handling
        panic("no node with that name")
        
    }
    current = cnode.Value
    
    for ii := 1; ii < m.Type.NumIn(); ii++ {
        in := m.Type.In(ii)
        dep := c.objectGraph.Get(in)
        
        var (
            value        interface{}
            instantiator *Instantiator
        )
        
        if dep == nil {
            
            value = c.typedSingletons[in]
            if value == nil {
                error := errors.New(fmt.Sprintf(
                    "Error in dependency-graph resolution ."+
                            "Required dependency of type <%s> but could not find it", in))
                panic(error)
            }
        }
        if value != nil {
            instantiator = &Instantiator{
                host:         host,
                objectType:   in,
                methodName:   m.Name,
                factory:      m.Func,
                dependencies: nil,
                context:      c,
                eager:        true,
            }
        } else {
            instantiator = dep.Value
        }
        
        if instantiator == nil {
            panic("Whoops.  None of those were registered")
        } else {
            log.Debugf("\t Dependency:%s", instantiator.objectType)
            current.dependencies = append(current.dependencies, instantiator)
        }
    }
    
    if m.Type.NumOut() == 0 {
        c.objectGraph.eager[m.Name] = cnode
    }
}

func resolveDetails(m *reflect.Method) (reflect.Type, bool) {
    var t reflect.Type = nil
    var eager bool = true
    for i := 0; i < m.Type.NumOut(); i++ {
        p := m.Type.Out(i)
        if p.Kind() == reflect.Struct ||
                p.Kind() == reflect.Interface ||
                p.Kind() == reflect.Ptr {
            if t != nil {
                panic("only supports single wired type")
            }
            t = p
        }
        if p.Kind() == reflect.Bool {
            eager = false
        }
    }
    return t, eager
}

func (c *DefaultApplicationContext) registerMethod(m *reflect.Method, host interface{}) {
    
    var (
        eager bool
        out   reflect.Type
    )
    dependencies := make([]*Instantiator, 0)
    if m.Type.NumOut() > 0 {
        out, eager = resolveDetails(m)
    }
    node := &Node{
        Key:     out,
        NameKey: m.Name,
        Value: &Instantiator{
            host:         host,
            objectType:   out,
            methodName:   m.Name,
            factory:      m.Func,
            dependencies: dependencies,
            context:      c,
            eager:        eager,
        },
    }
    c.objectGraph.Add(node, eager)
}

func NewApplicationContext() ApplicationContext {
    
    return &DefaultApplicationContext{
        typeMappings:    make(map[reflect.Type]interface{}, 0),
        objectGraph:     NewGraph(),
        namedSingletons: make(map[string]interface{}, 0),
        typedSingletons: make(map[reflect.Type]interface{}, 0),
    }
}
