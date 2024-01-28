package temp

import (
	"context"
	"fmt"
	"github.com/tinyredglasses/workers2/internal/jsutil"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"syscall/js"
)

var RuntimeContext = js.Global().Get("context")

type Task func(ctx context.Context) error

var callTask Task

func runCall(eventObj js.Value, runtimeCtxObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj, runtimeCtxObj)
	if err := callTask(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	runCallCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to runScheduler: %d", len(args)))
		}
		eventObj := args[0]
		runtimeCtxObj := jsutil.RuntimeContext
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := runCall(eventObj, runtimeCtxObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("runCall", runCallCallback)
}

//go:wasmimport workers ready
func ready()

// ScheduleTask sets the Task to be executed
func RunTask(t Task) {
	callTask = t
	ready()
	select {}
}
