package universe

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/xoctopus/x/misc/must"
)

const (
	prefix = "xoctopus__"
	suffix = "__xoctopus"
)

func unwrap(path string) string {
	if !strings.HasSuffix(path, suffix) || !strings.HasPrefix(path, prefix) {
		return path
	}
	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimSuffix(path, suffix)
	path = strings.ReplaceAll(path, "__", ".")
	path = strings.ReplaceAll(path, "_", "/")
	return path
}

func wrap(path string) string {
	if !strings.Contains(path, ".") && !strings.Contains(path, "/") {
		return path
	}
	must.BeTrueWrap(!strings.Contains(path, "_"), "unsupported pkg path which contains `_`")
	path = strings.ReplaceAll(path, ".", "__")
	path = strings.ReplaceAll(path, "/", "_")
	path = prefix + path + suffix
	return path
}

type Package interface {
	ID() string
	Path() string
	Name() string
}

func NewPackage(path string) Package {
	if path == "" {
		return nil
	}

	path = unwrap(path)

	if v, ok := gPackages.Load(path); ok {
		return v.(Package)
	}

	pkgs, err := packages.Load(&packages.Config{
		Overlay: make(map[string][]byte),
		Tests:   true,
		Mode:    0b11111111111111111,
	}, path)
	must.NoErrorWrap(err, "failed to load packages from %s", path)

	// TODO enable check loading package error; move NewPackage to pkgx
	// if len(pkgs[0].Errors) != 0 {
	// 	return nil
	// }

	pkg := pkgs[0].Types
	p := &_pkg_{Package: pkg, id: wrap(pkg.Path())}
	gPackages.Store(path, p)

	return p
}

type _pkg_ struct {
	id string
	*types.Package
}

func (p *_pkg_) ID() string {
	return p.id
}
