package resync

import (
	"fmt"
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func (o *one) Add(d int) {
	*o += one(d)
}

func (o *one) GetResult() int {
	*o++

	fmt.Println(*o)

	return int(*o)
}

func TestOnce_Do(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10

	for i := 0; i < N; i++ {
		go runDo(t, once, o, c)
	}

	for i := 0; i < N; i++ {
		<-c
	}

	if *o != 1 {
		t.Errorf("once failed outside runDo: %d is not 1", *o)
	}
}

func runDo(t *testing.T, once *Once, o *one, c chan bool) {
	t.Helper()

	once.Do(func() { o.Increment() })

	if v := *o; v != 1 {
		t.Errorf("once failed inside runDo: %d is not 1", v)
	}

	c <- true
}

func TestOnce_DoFunc(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10

	for i := 0; i < N; i++ {
		go runDoFunc(t, once, o, c, 1)
	}

	for i := 0; i < N; i++ {
		<-c
	}

	if *o != 1 {
		t.Errorf("once failed outside runDoFunc: %d is not 1", *o)
	}
}

func runDoFunc(t *testing.T, once *Once, o *one, c chan bool, param ...interface{}) {
	t.Helper()

	once.DoFunc(o, "Add", param...)

	if v := *o; v != 1 {
		t.Errorf("once failed inside runDoFunc: %d is not 1", v)
	}

	c <- true
}

func TestOnce_DoReturn(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10

	for i := 0; i < N; i++ {
		go runDoReturn(t, once, o, c)
	}

	for i := 0; i < N; i++ {
		<-c
	}

	if *o != 1 {
		t.Errorf("once failed outside runDoFunc: %d is not 1", *o)
	}
}

func runDoReturn(t *testing.T, once *Once, o *one, c chan bool) {
	t.Helper()

	ch := once.DoReturn(o, "GetResult")

	if v := *o; v != 1 {
		t.Errorf("once failed inside runDoFunc: %d is not 1", v)
	}

	if r, ok := (<-ch).(int); ok && r != 1 {
		t.Errorf("once failed inside runDoFunc: retunred result %d is not 1", r)
	}

	c <- true
}

func TestOnce_Reset(t *testing.T) {
	o := new(one)
	once := new(Once)

	once.Do(func() { o.Increment() })
	once.Reset()
	once.Do(func() { o.Increment() })

	if *o != 2 {
		t.Errorf("once failed outside runDo: %d is not 2", *o)
	}
}
