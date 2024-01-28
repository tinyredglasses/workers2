package temp

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/tinyredglasses/workers2/cloudflare/d1"
	"github.com/tinyredglasses/workers2/internal/jsutil"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"syscall/js"
)

var RuntimeContext = js.Global().Get("context")

type Task func(ctx context.Context) error

var task Task

func runTemp(eventObj js.Value, runtimeCtxObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj, runtimeCtxObj)
	if err := task(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	runTempCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
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
				err := runTemp(eventObj, runtimeCtxObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("runTemp", runTempCallback)
}

//go:wasmimport workers ready
func ready()

// RunTemp sets the Task to be executed
func RunTemp(t Task) {
	task = t
	ready()
	select {}
}

func Abc() {

}

func CreateD1() (driver.Connector, error) {
	fmt.Printf("%+v", RuntimeContext)
	ctx := runtimecontext.New(context.Background(), js.Value{}, RuntimeContext)

	connector, err := d1.OpenConnector(ctx, "DB")
	if err != nil {
		return nil, err
	}

	return connector, nil
}
