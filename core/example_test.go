package core

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

type SampleConfiguration struct {
}

type SampleConsumer struct {
	greeting string
}

func (s *SampleConsumer) Greet(name string) string {
	return fmt.Sprintf(s.greeting, name)
}

func (s *SampleConfiguration) GetGreeting() *string {
	u := "Hello %s"
	return &u
}

func (s *SampleConfiguration) HelloConsumer(hello *string) *SampleConsumer {
	return &SampleConsumer{greeting: *hello}
}

func TestGreetingWorks(t *testing.T) {

	ctx := NewApplicationContext()
	ctx.Scan(&SampleConfiguration{})

	consumer := ctx.GetByName("HelloConsumer").(*SampleConsumer)
	result := consumer.Greet("Bob")
	assert.Equal(t, result, "Hello Bob")
}
