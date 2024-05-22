package must_test

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sincospro/x/misc/must"
)

func ExampleNoError() {
	must.NoError(nil)

	defer func() {
		fmt.Println(recover())
	}()
	must.NoError(errors.New("any"))

	// Output:
	// any
}

func ExampleOK() {
	must.OK(true)

	defer func() {
		fmt.Println(recover())
	}()
	must.OK(false)

	// Output:
	// must ok
}
