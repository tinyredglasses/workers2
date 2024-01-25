package temp

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/tinyredglasses/workers2/cloudflare/d1"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"syscall/js"
)

var RuntimeContext = js.Global().Get("context")

func CreateD1() (driver.Connector, error) {
	fmt.Printf("%+v", RuntimeContext)
	ctx := runtimecontext.New(context.Background(), js.Value{}, RuntimeContext)

	connector, err := d1.OpenConnector(ctx, "DB")
	if err != nil {
		return nil, err
	}

	return connector, nil
}
