package jsutil

import (
	"fmt"
	"syscall/js"
	"time"
)

var (
	RuntimeContext      = js.Global().Get("context")
	Binding             = js.Global().Get("context").Get("binding")
	ObjectClass         = js.Global().Get("Object")
	PromiseClass        = js.Global().Get("Promise")
	RequestClass        = js.Global().Get("Request")
	ResponseClass       = js.Global().Get("Response")
	HeadersClass        = js.Global().Get("Headers")
	ArrayClass          = js.Global().Get("Array")
	Uint8ArrayClass     = js.Global().Get("Uint8Array")
	ErrorClass          = js.Global().Get("Error")
	ReadableStreamClass = js.Global().Get("ReadableStream")
	DateClass           = js.Global().Get("Date")
	Null                = js.ValueOf(nil)
)

func NewObject() js.Value {
	return ObjectClass.New()
}

func NewUint8Array(size int) js.Value {
	return Uint8ArrayClass.New(size)
}

func NewPromise(fn js.Func) js.Value {
	return PromiseClass.New(fn)
}

// ArrayFrom calls Array.from to given argument and returns result Array.
func ArrayFrom(v js.Value) js.Value {
	return ArrayClass.Call("from", v)
}

func AwaitPromise(promiseVal js.Value) (js.Value, error) {
	fmt.Println("AwaitPromise", "1")

	resultCh := make(chan js.Value)
	errCh := make(chan error)
	var then, catch js.Func
	fmt.Println("AwaitPromise", "2")

	then = js.FuncOf(func(_ js.Value, args []js.Value) any {
		fmt.Println("AwaitPromise", "6.1")
		defer then.Release()
		fmt.Println("AwaitPromise", "6.2")
		result := args[0]
		fmt.Println("AwaitPromise", "6.3")
		resultCh <- result
		fmt.Println("AwaitPromise", "6.4")
		return js.Undefined()
	})
	fmt.Println("AwaitPromise", "3")

	catch = js.FuncOf(func(_ js.Value, args []js.Value) any {
		fmt.Println("AwaitPromise", "7")
		defer catch.Release()
		result := args[0]
		errCh <- fmt.Errorf("failed on promise: %s", result.Call("toString").String())
		return js.Undefined()
	})
	fmt.Println("AwaitPromise", "4")

	promiseVal.Call("then", then).Call("catch", catch)
	fmt.Println("AwaitPromise", "5.1")

	select {
	case result := <-resultCh:
		fmt.Println("AwaitPromise", "5.2")
		return result, nil
	case err := <-errCh:
		fmt.Println("AwaitPromise", "5.3")
		return js.Value{}, err
	}
}

// StrRecordToMap converts JavaScript side's Record<string, string> into map[string]string.
func StrRecordToMap(v js.Value) map[string]string {
	entries := ObjectClass.Call("entries", v)
	entriesLen := entries.Get("length").Int()
	result := make(map[string]string, entriesLen)
	for i := 0; i < entriesLen; i++ {
		entry := entries.Index(i)
		key := entry.Index(0).String()
		value := entry.Index(1).String()
		result[key] = value
	}
	return result
}

// MaybeString returns string value of given JavaScript value or returns nil if the value is undefined.
func MaybeString(v js.Value) string {
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

// MaybeDate returns time.Time value of given JavaScript Date value or returns nil if the value is undefined.
func MaybeDate(v js.Value) (time.Time, error) {
	if v.IsUndefined() {
		return time.Time{}, nil
	}
	return DateToTime(v)
}

// DateToTime converts JavaScript side's Data object into time.Time.
func DateToTime(v js.Value) (time.Time, error) {
	milli := v.Call("getTime").Float()
	return time.UnixMilli(int64(milli)), nil
}

// TimeToDate converts Go side's time.Time into Date object.
func TimeToDate(t time.Time) js.Value {
	return DateClass.New(t.UnixMilli())
}
