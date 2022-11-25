package object

import "testing"

func TestStringHashKey(t *testing.T) {
	h1 := &String{Value: "Hello World"}
	h2 := &String{Value: "Hello World"}

	df1 := &String{Value: "My name is johnny"}
	df2 := &String{Value: "My name is johnny"}

	if h1.HashKey() != h2.HashKey() {
		t.Errorf("matching strings have different hash keys")
	}

	if df1.HashKey() != df2.HashKey() {
		t.Errorf("matching strings have different hash keys")
	}

	if h1.HashKey() == df1.HashKey() {
		t.Errorf("non matching strings have same hash keys")
	}
}
