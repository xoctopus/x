package internal

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/xoctopus/x/misc/must"
)

func NewPackage(path string) *types.Package {
	if path == "" {
		return nil
	}
	path = UnwrapPkgPath(path)

	if v, ok := gPackagesCache.Load(path); ok {
		return v.(*types.Package)
	}

	pkgs, err := packages.Load(&packages.Config{
		Overlay: make(map[string][]byte),
		Tests:   true,
		Mode:    65535,
	}, path)
	must.NoErrorWrap(err, "failed to load packages from %s", path)

	// TODO enable check loading package error; move NewPackage to pkgx
	// if len(pkgs[0].Errors) != 0 {
	// 	return nil
	// }

	pkg := pkgs[0].Types
	gPackagesCache.Store(path, pkg)

	return pkg
}
