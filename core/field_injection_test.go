package core

import "testing"

type InjectableThing struct {
	Name string `Inject:MyBean`
}

func TestFieldInjection(t *testing.T) {

}
