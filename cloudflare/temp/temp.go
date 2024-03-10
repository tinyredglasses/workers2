package temp

import (
	"context"
	"fmt"
	"github.com/tinyredglasses/workers2/internal/jsutil"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"syscall/js"
)

type Task func(ctx context.Context) error

var task Task

func handleData(eventObj js.Value, runtimeCtxObj js.Value) error {
	fmt.Println("handleData")
	//fmt.Println("handleData1", eventObj, runtimeCtxObj)

	ctx := runtimecontext.New(context.Background(), eventObj, runtimeCtxObj)

	//wc := runtimeCtxObj.Get("client")
	//fmt.Println(runtimeCtxObj.Get("client").IsUndefined())
	//fmt.Println(runtimeCtxObj.Get("ctx").Get("client").IsUndefined())

	//v := runtimecontext.TryExtractRuntimeObj(ctx)
	//e := v.Get("env")
	//fmt.Println(e.IsUndefined())
	//fmt.Println("handleData2")

	if err := task(ctx); err != nil {
		return err
	}
	//fmt.Println("handleData3")

	return nil
}

func init() {
	handleDataCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
		fmt.Println("handleDataCallback0", args)

		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to handleData: %d", len(args)))
		}
		eventObj := args[0]
		fmt.Println("handleDataCallback1", eventObj)
		runtimeCtxObj := jsutil.RuntimeContext
		fmt.Println("handleDataCallback2", runtimeCtxObj)

		fsdf := js.Global().Get("JSON").Call("stringify", runtimeCtxObj)
		fmt.Println(fsdf)

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := handleData(eventObj, runtimeCtxObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("handleData", handleDataCallback)
}

//go:wasmimport workers ready
func ready()

// RunTemp sets the Task to be executed
func RunTemp(t Task) {
	task = t
	ready()
	select {}
}

//func CreateD1() (driver.Connector, error) {
//	fmt.Printf("%+v", RuntimeContext)
//	ctx := runtimecontext.New(context.Background(), js.Value{}, RuntimeContext)
//
//	connector, err := d1.OpenConnector(ctx, "DB")
//	if err != nil {
//		return nil, err
//	}
//
//	return connector, nil
//}
