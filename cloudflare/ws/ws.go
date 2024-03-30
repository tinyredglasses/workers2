package ws

import (
	"context"
	"fmt"
	"github.com/tinyredglasses/workers2/internal/jsutil"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"log/slog"
	"syscall/js"
)

var (
	messageHandler MessageHandler
	logger         = slog.With("package", "ws")

	sender Sender
)

type MessageHandler interface {
	Handle(ctx context.Context, reqObj js.Value)
}

type MessageHandlerCreator func(ctx context.Context) MessageHandler

func init() {
	logger.Info("init")

	//outerRuntimeCtxObj := jsutil.RuntimeContext
	//ctx := runtimecontext.New(context.Background(), js.Value{}, outerRuntimeCtxObj)
	//websocketClient := cloudflare.GetWebsocketClient(ctx, "")
	//sender = Sender{websocketClient: websocketClient}

	handleDataCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {

		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to handleData: %d", len(args)))
		}
		eventObj := args[0]
		runtimeCtxObj := jsutil.RuntimeContext

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

func handleData(event js.Value, runtimeCtx js.Value) error {
	logger.Info("handleData")
	ctx := runtimecontext.New(context.Background(), event, runtimeCtx)

	messageHandler.Handle(ctx, event)
	return nil
}

//go:wasmimport workers ready
func ready()

func Handle(mhc MessageHandlerCreator) {
	logger.Info("Handle")
	ctx := runtimecontext.New(context.Background(), js.Value{}, jsutil.RuntimeContext)
	messageHandler = mhc(ctx)
	ready()
	select {}
}
