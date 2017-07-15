package core

import "reflect"

type Instantiator struct {
    eager        bool
    methodName   string
    host         interface{}
    objectType   interface{}
    factory      reflect.Value
    dependencies []*Instantiator
    instance     interface{}
    context      *DefaultApplicationContext
}

func (i *Instantiator) recoverFromInjection() {
    if r := recover(); r != nil {
        log.Errorf("Injection failed on instantiator: {"+
                "eager = %t, "+
                "methodName = %s, "+
                "host = %v#, "+
                "object type = %s}",
            i.eager,
            i.methodName,
            i.host,
            i.objectType,
        )
        panic(r)
    }
}

func (i *Instantiator) Construct() interface{} {
    
    defer i.recoverFromInjection()
    
    if i.instance != nil {
        return i.instance
    }
    
    if i.context.namedSingletons[i.methodName] != nil {
        return i.context.namedSingletons[i.methodName]
    }
    
    if i.objectType != nil &&
            i.context.typedSingletons[i.objectType.(reflect.Type)] != nil {
        return i.context.typedSingletons[i.objectType.(reflect.Type)]
    }
    
    var result []reflect.Value
    if i.dependencies != nil && len(i.dependencies) > 0 {
        
        arguments := make([]reflect.Value, 0)
        arguments = append(arguments, reflect.ValueOf(i.host))
        for _, dependentInstantiator := range i.dependencies {
            arguments = append(
                arguments,
                reflect.ValueOf(dependentInstantiator.Construct()))
        }
        result = i.factory.Call(arguments)
    } else {
        result = i.factory.Call([]reflect.Value{reflect.ValueOf(i.host)})
    }
    
    scope := i.resolveScope(result)
    name := i.resolveName(result)
    if result != nil && len(result) > 0 {
        value := i.resolveValue(result).Interface()
        if scope == Singleton {
            i.context.namedSingletons[name] = value
            i.context.typedSingletons[reflect.TypeOf(value)] = value
        }
        i.firePostConstruct(value, i.objectType.(reflect.Type))
        return value
    }
    return nil
    
}

func (u *Instantiator) firePostConstruct(v interface{}, t reflect.Type) {
    m, b := t.MethodByName("PostConstruct")
    if b {
        m.Func.Call([]reflect.Value{reflect.ValueOf(v), reflect.ValueOf(u.context)})
    }
}

func (u *Instantiator) resolveValue(result []reflect.Value) reflect.Value {
    for _, v := range result {
        if v.Kind() == reflect.Interface || v.Kind() == reflect.Struct {
            return v
        }
        if v.Kind() == reflect.Ptr {
            return v
        }
    }
    panic("Illegal state")
}

func (u *Instantiator) resolveName(result []reflect.Value) string {
    for _, v := range result {
        if v.Kind() == reflect.String {
            return v.String()
        }
    }
    return u.methodName
}

func (u *Instantiator) resolveScope(result []reflect.Value) Scope {
    for _, v := range result {
        if v.Kind() == reflect.Uint16 {
            return Scope(v.Uint())
        }
    }
    return Singleton
}
