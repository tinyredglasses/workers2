package cloudflare

import (
	"context"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"syscall/js"

	"github.com/tinyredglasses/workers2/cloudflare/internal/cfruntimecontext"
)

// Getenv gets a value of an environment variable.
//   - https://developers.cloudflare.com/workers/platform/environment-variables/
//   - This function panics when a runtime context is not found.
func Getenv(ctx context.Context, name string) string {
	return cfruntimecontext.MustGetRuntimeContextEnv(ctx).Get(name).String()
}

// GetBinding gets a value of an environment binding.
//   - https://developers.cloudflare.com/workers/platform/bindings/about-service-bindings/
//   - This function panics when a runtime context is not found.
func GetBinding(ctx context.Context, name string) js.Value {
	return cfruntimecontext.MustGetRuntimeContextEnv(ctx).Get(name)
}

func GetWebsocketClient(ctx context.Context, name string) js.Value {
	return cfruntimecontext.MustGetRuntimeContextValue(ctx, "client")
}

func GetCtx(ctx context.Context, name string) js.Value {
	return cfruntimecontext.MustGetExecutionContext(ctx).Get(name)
}

func GetTriggerObject(ctx context.Context) js.Value {
	return runtimecontext.MustExtractTriggerObj(ctx)
}
