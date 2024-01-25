package workers

import (
	"context"
	"database/sql/driver"
	"github.com/tinyredglasses/workers2/cloudflare/d1"
	"github.com/tinyredglasses/workers2/internal/runtimecontext"
	"syscall/js"
)

func CreateD1(runtimeCtxObj js.Value) (driver.Connector, error) {

	ctx := runtimecontext.New(context.Background(), js.Value{}, runtimeCtxObj)

	connector, err := d1.OpenConnector(ctx, "DB")
	if err != nil {
		return nil, err
	}

	return connector, nil
}
