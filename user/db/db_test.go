package db

import "testing"

func Test_mapper(t *testing.T) {
	s := mapper("helloWorld")
	expected := "hello_world"
	if s != expected{
		t.Errorf("Expected %s but %s", expected, s)
	}
}
