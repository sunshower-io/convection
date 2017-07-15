package core

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MultipleContextsSuite struct {
	suite.Suite
	context *DefaultApplicationContext
}

func (s *MultipleContextsSuite) SetupTest() {
	s.context = NewApplicationContext().(*DefaultApplicationContext)
}

type MCfg struct {
}

var j int = 0

type DependsOnMultiple struct {
	*ServiceDependencyRef
	DependentDependentType
}

type PrototypeDependsOnSingleton struct {
	*DependsOnMultiple
}

var called = false

func (i *PrototypeDependsOnSingleton) PostConstruct(ctx ApplicationContext) {
	println("Called!")
	called = true
}


func (c *MCfg) PrototypeDependsOnMultiplesingleton(d *DependsOnMultiple) *PrototypeDependsOnSingleton {
	return &PrototypeDependsOnSingleton{d}
}

var dependsOn * DependsOnMultiple

func (c *MCfg) ConfiguredWithoutReturning(d *DependsOnMultiple) {
	dependsOn = d
}

func (c *MCfg) CreateMultiple(s *ServiceDependencyRef, d DependentDependentType) (*DependsOnMultiple, Scope) {
	j++
	return &DependsOnMultiple{
		s,
		d,
	}, Singleton
}

func (c *MCfg) CreateServiceDependency() *ServiceDependencyRef {
	dep := new(ServiceDependencyRef)
	dep.Name = "Overridden"
	return dep
}

func (m *MultipleContextsSuite) TestNodeRedefine() {
	m.context.objectGraph.Add(&Node{
		Key:     "helloIface",
		NameKey: "helloIface2",
		Value: &Instantiator{
			methodName: "Frapper",
		},
	}, true)

	m.context.objectGraph.Add(&Node{
		Key:     "helloIface",
		NameKey: "helloIface2",
		Value: &Instantiator{
			methodName: "Coolbeans",
		},
	}, true)

	m.Equal(
		m.context.objectGraph.NamedNodes["helloIface2"].Value.methodName,
		"Coolbeans",
	)
}

func (m *MultipleContextsSuite) TestMultipleContextsAreResolvedCorrect() {
	m.context.Scan(&Cfg{}, &RefCfg{}, &MCfg{})
	for i := 0; i < 20; i++ {
		m.context.GetByName("CreateMultiple")
	}
	res := m.context.GetByName("CreateMultiple").(*DependsOnMultiple)
	m.Equal(res.DependentDependentType.DependentType.SampleType.Name, "Frapper")
	m.Equal(j, 2)
	j = 0
}

func (m *MultipleContextsSuite) TestOverridesWork() {
	m.context.Scan(&Cfg{}, &RefCfg{}, &MCfg{})
	res := m.context.GetByName("CreateServiceDependency").(*ServiceDependencyRef)
	m.Equal(res.Name, "Overridden")
}

func (m *MultipleContextsSuite) TestPostConstructCalled() {
	m.context.Scan(&Cfg{}, &RefCfg{}, &MCfg{})
	m.True(called)

}

func (m *MultipleContextsSuite) TestFunctionReturningNothingIsConfigured() {
	m.context.Scan(&Cfg{}, &RefCfg{}, &MCfg{})
	m.NotNil(dependsOn)
	
}

func (m *MultipleContextsSuite) TestDependencyWorks() {
	m.context.Scan(&Cfg{}, &RefCfg{}, &MCfg{})
	v := m.context.GetByName("PrototypeDependsOnMultiplesingleton").(*PrototypeDependsOnSingleton)
	m.Equal(j, 1)
	m.Equal(v.DependsOnMultiple.ServiceDependencyRef.Name, "Overridden")
	j = 0

}

func TestMultipleCfg(t *testing.T) {
	suite.Run(t, &MultipleContextsSuite{})
}
