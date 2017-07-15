#  Convection 

Convection is a lightweight, easy-to-use yet powerful dependency-injection framework for the Go programming language.  
Currently, only constructor-function injection is supported (as field reflection in Go is iffy at best)



## Usage


```go 

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

func(s *SampleConfiguration) HelloConsumer(hello *string) *SampleConsumer {
    return &SampleConsumer{greeting : *hello}
}


func TestGreetingWorks(t *testing.T)  {
    
    ctx := NewApplicationContext()
    ctx.Scan(&SampleConfiguration{})
    
    consumer := ctx.GetByName("HelloConsumer").(*SampleConsumer)
    result := consumer.Greet("Bob")
    assert.Equal(t,  result, "Hello Bob")
}

```

### Features

- Unlimited configurations supported
- Extremely lightweight 
- Cyclic dependency detection
- No external dependencies 

### Limitations


- Only function injection is supported
- Container objects should be references.  Prototype injection isn't quite there yet

