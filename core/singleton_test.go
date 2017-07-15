package core

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

type InjectedSingleton struct {

}

type InjectedSingletonDependent struct {
    singleton *InjectedSingleton
}


type SingletonConfiguration struct {

}

func (s *SingletonConfiguration) Inject(
        i *InjectedSingleton,
) *InjectedSingletonDependent {

    return &InjectedSingletonDependent{singleton: i}
}


func TestSingletonsAreInjected(t *testing.T) {
    ctx := NewApplicationContext()
    ctx.RegisterSingleton(&InjectedSingleton{})
    ctx.Scan(&SingletonConfiguration{})
    singleton := ctx.GetByName("Inject").(*InjectedSingletonDependent)
    assert.NotNil(t, singleton)
    assert.NotNil(t, singleton.singleton)
}

