package core

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

func init() {

}

type SampleType struct {
	Name string
}

type AlltogetherValue struct {
	SampleType             SampleType
	DependentType          DependentType
	DependentDependentType DependentDependentType
}

func (s *SampleType) Constructor() {
	s.Name = "Frapper"
}

type DependentType struct {
	SampleType SampleType
}

type DependentDependentType struct {
	SampleType    SampleType
	DependentType DependentType
}

type Cfg struct {
}
type NamedStruct struct {
}

func (d *Cfg) CreateNamedStruct() (NamedStruct, Scope) {
	return NamedStruct{}, Singleton
}

func (d *Cfg) CreateAltogetherNow() AlltogetherValue {
	return AlltogetherValue{}
}

func (d *Cfg) CreateSampleType() SampleType {
	return SampleType{Name: "Frapper"}
}

func (d *Cfg) CreateDependentType(t SampleType) DependentType {
	return DependentType{SampleType: t}
}

func (d *Cfg) CreateDependentDependentType(t DependentType, u SampleType) DependentDependentType {
	return DependentDependentType{DependentType: t, SampleType: u}
}

type ApplicationContextSuite struct {
	suite.Suite
	context *DefaultApplicationContext
}

func (s *ApplicationContextSuite) SetupTest() {
	s.context = NewApplicationContext().(*DefaultApplicationContext)
}

func (s *ApplicationContextSuite) TestResolvingByNameProducesExpectedResults() {

	cfg := new(Cfg)
	s.context.Scan(cfg)
	result := s.context.GetByName("CreateSampleType").(SampleType)
	s.Equal(result.Name, "Frapper")
}

func (s *ApplicationContextSuite) TestRegisteringTypeResultsInTypeBeingRegistered() {

	sampleType := new(SampleType)
	s.context.RegisterSingleton(sampleType)
	s.True(s.context.Contains(sampleType))
}

func (s *ApplicationContextSuite) TestResolvingRegisteredTypeByTypeReturnsCorrectResults() {
	sampleType := new(SampleType)
	s.context.RegisterSingleton(sampleType)
	t := reflect.TypeOf(sampleType)
	s.Equal(s.context.GetByType(t), sampleType)
}

func (s *ApplicationContextSuite) TestScanningAConfigurationMustProduceTheCorrectNodeCount() {
	cfg := new(Cfg)
	s.context.Scan(cfg)
	s.Equal(len(s.context.objectGraph.Nodes), 6)
}

func (s *ApplicationContextSuite) TestScanningAConfigurationMustConstructASimpleDependentTypeCorrectly() {
	cfg := new(Cfg)
	s.context.Scan(cfg)
	result := s.context.Get(reflect.TypeOf(SampleType{})).(SampleType)
	s.Equal(result.Name, "Frapper")
}

func (s *ApplicationContextSuite) TestScanningAConfigurationMustInstantiateDependentTypes() {
	cfg := new(Cfg)
	s.context.Scan(cfg)
	dt := s.context.Get(reflect.TypeOf(DependentType{})).(DependentType)
	s.NotNil(dt)
	s.NotNil(dt.SampleType)
	s.Equal(dt.SampleType.Name, "Frapper")
}

func (s *ApplicationContextSuite) TestScanningDependenciesToMultipleLevelsPopulatesFields() {
	cfg := new(Cfg)
	s.context.Scan(cfg)
	dt := s.context.Get(reflect.TypeOf(DependentDependentType{})).(DependentDependentType)
	name := dt.DependentType.SampleType.Name
	s.NotNil(dt.SampleType)
	s.Equal(name, "Frapper")
}

func (s *ApplicationContextSuite) TestAutowiringWorks() {

	//cfg := new(Cfg);
	//s.context.Scan(cfg);
	//dt := s.context.Get(reflect.TypeOf(AlltogetherValue{})).(AlltogetherValue);
	//println(fmt.Sprintf("#v", dt));
}

func TestApplicationContextSuite(t *testing.T) {
	suite.Run(t, &ApplicationContextSuite{})
}
