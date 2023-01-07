package resync

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Once is an object that will perform exactly one action.
// A Once must not be copied after first use.
type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	done uint32
	m    sync.Mutex
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once.
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

// DoFunc calls the specified method if and only if it is being called for the
// first time for this instance of Once, given the instance on which the method is
// implemented and its parameters.
func (o *Once) DoFunc(m interface{}, methodName string, param ...interface{}) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doFuncSlow(m, methodName, param...)
	}
}

func (o *Once) doFuncSlow(m interface{}, methodName string, param ...interface{}) (result []reflect.Value) {
	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)

		method := reflect.ValueOf(m).MethodByName(methodName)
		in := make([]reflect.Value, method.Type().NumIn())

		for i := 0; i < method.Type().NumIn(); i++ {
			in[i] = reflect.ValueOf(param[i])
		}

		result = method.Call(in)
	}

	return
}

// DoReturn calls the specified method if and only if it is being called for the
// first time for this instance of Once, given the instance on which the method is
// implemented and its parameters. It returns a channel of type interface to send the
// return values of the called method.
func (o *Once) DoReturn(m interface{}, methodName string, param ...interface{}) <-chan interface{} {
	if atomic.LoadUint32(&o.done) == 0 {
		result := o.doFuncSlow(m, methodName, param...)
		ch := make(chan interface{}, len(result))
		defer close(ch)

		for i := range result {
			ch <- result[i].Interface()
		}

		return ch
	}

	return nil
}

// Reset resets the once object so that it can be re-used.
func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()

	atomic.StoreUint32(&o.done, 0)
}
