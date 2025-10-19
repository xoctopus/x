package main

import (
	"context"
	"os"
	"path/filepath"

	_ "github.com/xoctopus/genx/devpkg/errorx"
	"github.com/xoctopus/genx/pkg/genx"

	"github.com/xoctopus/x/misc/must"
)

func main() {
	cwd := must.NoErrorV(os.Getwd())

	must.NoError(
		genx.NewContext(&genx.Args{
			Entrypoint: []string{
				filepath.Join(cwd, "reflectx"),
				filepath.Join(cwd, "textx"),
			},
		}).Execute(context.Background(), genx.Get()...),
	)
}
