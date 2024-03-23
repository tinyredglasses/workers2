package temp

import (
	"context"
	"fmt"
	"github.com/tinyredglasses/workers2/internal/jsutil"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"log/slog"
	"syscall/js"
)

func init() {
	slog.Info("init")

	handleDataCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {

		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to handleData: %d", len(args)))
		}
		eventObj := args[0]
		//fmt.Println("handleDataCallback1", eventObj)
		runtimeCtxObj := jsutil.RuntimeContext
		//fmt.Println("handleDataCallback2", runtimeCtxObj)

		//fsdf1 := js.Global().Get("JSON").Call("stringify", eventObj)
		//fsdf2 := js.Global().Get("JSON").Call("stringify", runtimeCtxObj)
		//
		//fmt.Println(fsdf1, fsdf2)

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

func main() {

}

var (
	handler MessageHandler
	//closeCh = make(chan struct{})
)

//go:wasmimport workers ready
func ready()

type MessageHandler interface {
	handle(ctx context.Context, reqObj js.Value)
}

func handleRequest(reqObj js.Value, runtimeCtxObj js.Value) {
	//req, err := jshttp.ToRequest(reqObj)
	//if err != nil {
	//	panic(err)
	//}
	//e := runtimeCtxObj.Get("env")
	//fmt.Println(e)
	//fmt.Printf("%+v\n", e)
	ctx := runtimecontext.New(context.Background(), reqObj, runtimeCtxObj)

	//req = req.WithContext(ctx)
	//reader, writer := io.Pipe()
	//w := &jshttp.ResponseWriter{
	//	HeaderValue: http.Header{},
	//	StatusCode:  http.StatusOK,
	//	Reader:      &appCloser{reader},
	//	Writer:      writer,
	//	ReadyCh:     make(chan struct{}),
	//}
	handler.handle(ctx, reqObj)
}

func HandleMessages(h MessageHandler) {
	handler = h
	ready()
	select {
	//case <-closeCh:
	}
}
