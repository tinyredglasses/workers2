package d1

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/tinyredglasses/workers2/internal/jsutil"
)

type stmt struct {
	stmtObj js.Value
}

var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)
)

func (s *stmt) Close() error {
	// do nothing
	return nil
}

// NumInput is not supported and always returns -1.
func (s *stmt) NumInput() int {
	return -1
}

func (s *stmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, errors.New("d1: Exec is deprecated and not implemented")
}

// ExecContext executes prepared statement.
// Given []drier.NamedValue's `Name` field will be ignored because Cloudflare D1 client doesn't support it.
func (s *stmt) ExecContext(_ context.Context, args []driver.NamedValue) (driver.Result, error) {
	argValues := make([]any, len(args))
	for i, arg := range args {
		argValues[i] = arg.Value
	}
	resultPromise := s.stmtObj.Call("bind", argValues...).Call("run")
	resultObj, err := jsutil.AwaitPromise(resultPromise)
	if err != nil {
		return nil, err
	}
	return &result{
		resultObj: resultObj,
	}, nil
}

func (s *stmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, errors.New("d1: Query is deprecated and not implemented")
}

func (s *stmt) QueryContext(_ context.Context, args []driver.NamedValue) (driver.Rows, error) {
	fmt.Println("QueryContext1")
	argValues := make([]any, len(args))
	for i, arg := range args {
		argValues[i] = arg.Value
	}
	fmt.Println("QueryContext2")

	resultPromise := s.stmtObj.Call("bind", argValues...).Call("all")
	fmt.Println("QueryContext3")

	rowsObj, err := jsutil.AwaitPromise(resultPromise)
	fmt.Println("QueryContext4")

	if err != nil {
		return nil, err
	}
	fmt.Println("QueryContext5")

	if !rowsObj.Get("success").Bool() {
		return nil, errors.New("d1: failed to query")
	}
	fmt.Println("QueryContext6")

	return &rows{
		rowsObj: rowsObj.Get("results"),
	}, nil
}
