package cron

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/tinyredglasses/workers2/internal/jsutil"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
)

type Task func(ctx context.Context) error

var scheduledTask Task

func runScheduler(eventObj js.Value, runtimeCtxObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj, runtimeCtxObj)
	if err := scheduledTask(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	runSchedulerCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
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
				err := runScheduler(eventObj, runtimeCtxObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("runScheduler", runSchedulerCallback)
}

//go:wasmimport workers ready
func ready()

// ScheduleTask sets the Task to be executed
func ScheduleTask(task Task) {
	scheduledTask = task
	ready()
	select {}
}
