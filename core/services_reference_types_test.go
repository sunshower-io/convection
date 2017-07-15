package core

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

var count = 0

type DependsOnContext struct {
	context ApplicationContext
}

type ServiceRef struct {
	*ServiceDependencyRef
}

type ServiceDependencyRef struct {
	Name string
}

type SingletonRef struct {
	*ServiceDependencyRef
}

type ThirdRef struct {
}

type RefCfg struct {
}

func (r *RefCfg) CreateThirdRef() *ThirdRef {
	return &ThirdRef{}
}

func (r *RefCfg) CreateApplicationContextDependent(context ApplicationContext) *DependsOnContext {
	return &DependsOnContext{
		context: context,
	}
}

func (r *RefCfg) CreateSingletonRef(sev *ServiceDependencyRef, d *ServiceRef, r1 *ThirdRef) (*SingletonRef, Scope) {
	count++
	return &SingletonRef{
		sev,
	}, Singleton
}

func (r *RefCfg) CreateServiceRef(dep *ServiceDependencyRef) *ServiceRef {
	return new(ServiceRef)
}

func (r *RefCfg) CreateServiceDependency() *ServiceDependencyRef {
	cfg := new(ServiceDependencyRef)
	cfg.Name = "original"
	return cfg
}

type ServiceRefContextSuite struct {
	suite.Suite
	context *DefaultApplicationContext
}

func (s *ServiceRefContextSuite) SetupTest() {
	s.context = NewApplicationContext().(*DefaultApplicationContext)
}

func (s *ServiceRefContextSuite) TestDependentOfApplicationContextIsInjectedCorrectly() {
	s.context.Scan(&RefCfg{})
	r := s.context.Get(reflect.TypeOf(&DependsOnContext{}))
	s.NotNil(r)
	f := r.(*DependsOnContext)
	s.NotNil(f.context)
}

func(s *ServiceRefContextSuite) TestResolvingNonExistingReferenceMustReturnError() {
	
	s.context.Scan(&RefCfg{})
	_, err := s.context.ByName("NotHere =)")
	s.NotNil(err, "ByName should've returned an error for a non-existant object")
}

func (s *ServiceRefContextSuite) TestSingletonMustBeInjectedCorrectly() {
	s.context.Scan(&RefCfg{})
	val := s.context.Get(reflect.TypeOf(&SingletonRef{})).(*SingletonRef)
	val1 := s.context.GetByName("CreateSingletonRef").(*SingletonRef)
	s.NotNil(val)
	s.True(val1 == val)
	s.Equal(count, 1)
	count = 0

}

func (s *ServiceRefContextSuite) TestServiceRefIsInjected() {
	s.context.Scan(&RefCfg{})
	val := s.context.Get(reflect.TypeOf(&ServiceDependencyRef{})).(*ServiceDependencyRef)
	s.NotNil(val)
	count = 0
}

func (s *ServiceRefContextSuite) TestServiceRefIsInjectedByName() {
	s.context.Scan(&RefCfg{})
	val := s.context.GetByName("CreateServiceDependency").(*ServiceDependencyRef)
	s.NotNil(val)

	count = 0
}

func TestServiceRefSuite(t *testing.T) {
	suite.Run(t, &ServiceRefContextSuite{})
}
