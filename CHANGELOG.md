
<a name="v0.3.2"></a>
## [v0.3.2](https://github.com/xoctopus/x/compare/v0.3.1...v0.3.2)

> 2026-02-13

### Chore

* add no checkers warning


<a name="v0.3.1"></a>
## [v0.3.1](https://github.com/xoctopus/x/compare/v0.3.0...v0.3.1)

> 2026-02-13

### Chore

* alias bdd.TB as testing.TB

### Doc

* update CHANGELOG


<a name="v0.3.0"></a>
## [v0.3.0](https://github.com/xoctopus/x/compare/v0.2.12...v0.3.0)

> 2026-02-05

### Feat

* **contextx:** generic context provider
* **slicex:** range slice elements to converted values


<a name="v0.2.12"></a>
## [v0.2.12](https://github.com/xoctopus/x/compare/v0.2.11...v0.2.12)

> 2026-01-19

### Chore

* format config and regen

### Ci

* regen project configurations

### Feat

* **testx:** add bdd for Behavior-Driver Development


<a name="v0.2.11"></a>
## [v0.2.11](https://github.com/xoctopus/x/compare/v0.2.10...v0.2.11)

> 2026-01-05

### Ci

* add makefile changelog gen entry and regen changelog

### Doc

* add changelog

### Feat

* **slicex:** unique mapping for keys and values


<a name="v0.2.10"></a>
## [v0.2.10](https://github.com/xoctopus/x/compare/v0.2.9...v0.2.10)

> 2026-01-03

### Feat

* **misc:** defer error collection


<a name="v0.2.9"></a>
## [v0.2.9](https://github.com/xoctopus/x/compare/v0.2.8...v0.2.9)

> 2026-01-02

### Feat

* **stringsx:** random n


<a name="v0.2.8"></a>
## [v0.2.8](https://github.com/xoctopus/x/compare/v0.2.7...v0.2.8)

> 2025-12-31

### Feat

* docx


<a name="v0.2.7"></a>
## [v0.2.7](https://github.com/xoctopus/x/compare/v0.2.6...v0.2.7)

> 2025-12-30

### Chore

* add linter and fixing


<a name="v0.2.6"></a>
## [v0.2.6](https://github.com/xoctopus/x/compare/v0.2.5...v0.2.6)

> 2025-12-28

### Feat

* **slicex:** unique mapping


<a name="v0.2.5"></a>
## [v0.2.5](https://github.com/xoctopus/x/compare/v0.2.4...v0.2.5)

> 2025-12-25

### Feat

* **testx:** add assertion for codex.Error


<a name="v0.2.4"></a>
## [v0.2.4](https://github.com/xoctopus/x/compare/v0.2.3...v0.2.4)

> 2025-12-21

### Fix

* **flagx:** use generic underlyings


<a name="v0.2.3"></a>
## [v0.2.3](https://github.com/xoctopus/x/compare/v0.2.2...v0.2.3)

> 2025-12-16

### Feat

* **flagx:** add flag bit marker


<a name="v0.2.2"></a>
## [v0.2.2](https://github.com/xoctopus/x/compare/v0.2.1...v0.2.2)

> 2025-12-08

### Feat

* remove quoted option value


<a name="v0.2.1"></a>
## [v0.2.1](https://github.com/xoctopus/x/compare/v0.2.0...v0.2.1)

> 2025-12-02

### Chore

* remove gomega dependencies
* bump dependencies and adaption
* use codex.Error
* bump dependencies

### Ci

* use latest go ci env
* update Makefile
* update Makefile

### Feat

* **contextx:** non-override option
* **initializer:** add universal initializer
* **reflectx:** DerefPointer
* **slicex:** equivalent slice
* **syncx:** new syncx.Set with initial keys and add once_override

### Fix

* t.Helper in ExpectPanic
* **syncx:** NewSet instead pointer return

### Refact

* remove too mach dependencies

### Test

* **syncx:** fix unit test


<a name="v0.2.0"></a>
## [v0.2.0](https://github.com/xoctopus/x/compare/v0.1.2...v0.2.0)

> 2025-10-23

### Chore

* remove gen tools
* remove gomega and mapx
* clean code
* **ci:** pretty Makefile

### Ci

* add codecov.yaml
* add codecov.yaml
* update Makefile
* remove golangci-lint and update modules
* update Makefile

### Doc

* **contextx:** add doc and benchmark testing

### Feat

* upgrade modules
* **iterx:** a extended pkg for std iter
* **mapx:** removed, use syncx.Map
* **must:** assertion with F, Wrap, WrapF
* **must:** add must.Success must.OK
* **reflectx:** use stringsx validations
* **slicex:** add slicex.Unique
* **slicex:** remove slicex
* **stringsx:** common reg validations
* **stringsx:** common reg expressions
* **syncx:** typed sync.map
* **syncx:** update syncx Map/Set/Pool
* **testx:** add matcher Succeed/Failed
* **testx:** add assertion Succeed and Failed
* **testx:** a unit testing tool like gomega
* **textx:** support ptr implements of TextMarshaler

### Fix

* **reflectx:** fix CanElem
* **reflectx:** improve key's validations
* **stringsx:** fix flag name regex
* **testx:** fix Matcher BeNil

### Test

* add unit test
* **reflectx:** replace to testx
* **reflectx:** use ExceptPanic
* **reflectx:** add unit tests
* **stringsx:** add single letter test case


<a name="v0.1.2"></a>
## [v0.1.2](https://github.com/xoctopus/x/compare/v0.1.1...v0.1.2)

> 2025-09-21

### Chore

* **ci:** update go version in ci.yml

### Feat

* **slicex:** slice mapping

### Fix

* **testx:** recover assertion

### Test

* **reflectx:** pretty unit tests


<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/xoctopus/x/compare/v0.1.0...v0.1.1)

> 2025-09-21

### Ci

* remove lint step in ci flow
* remove golangci-lint ci flow and bump go to 1.24
* add golangci-lint ci flow

### Feat

* **reflectx:** structure tag parsing with flag has options


<a name="v0.1.0"></a>
## v0.1.0

> 2025-06-01

### Build

* update ci.yml and Makefile
* fix Makefile on `report` entry

### Chore

* upgrade go.mod dependencies
* typo
* upgrade dependencies
* update readme
* add make and ci

### Ci

* upgrade actions/checkout and actions/setup-go
* tags trigger ci workflow

### Doc

* reference badge

### Docs

* update README

### Feat

* upgrade dependencies
* **contextx:** cleanup and add context composer
* **contextx:** context composer
* **contextx:** generic context injector and loader
* **mapx:** add mapx.Exists
* **mapx:** add enhanced map implements with map[k]v and sync.Map
* **misc:** strings helper
* **misc:** must util funcs
* **reflectx:** indirectNew to allocate and return the deepest indirect reflect value
* **reflectx:** fix struct tag parsing
* **reflectx:** hack unexported struct field
* **reflectx:** type assert and casting
* **reflectx:** add reflectx.DeepCopy and reflectx.Clone for value copy and remove reflectx.Set
* **reflectx:** reflectx
* **reflectx:** improve IsZero, more stricter
* **reflectx:** check if reflect.Type can call Elem without panic
* **reflectx:** add reflectx.DeepCopy and Clone for value copy
* **reflectx:** add Set function to shallow copy reflect.Value
* **reflextx:** add `IsNumeric` `IsInteger` and `IsFloat` functions
* **resultx:** add resultx for universal function result handler
* **resultx:** add result Succeed and Failed methods
* **resultx:** function results handler
* **textx:** add MarshalURL and UnmarshalURL for conversion between url.Value and struct
* **textx:** add simple text (de)serializer
* **textx:** use fmt.Sscan to parse numeric types from string to support multi numerical bases input
* **typex:** support generic types
* **typex:** typex for univeral reflect and go type abstract
* **typex:** generic supported
* **universe:** add universe type parse by reflect.Type and type string

### Fix

* **mapx:** simplify Map.LoadEq and Map.LoadEqs
* **reflectx:** reflectx.IndirectNew
* **reflectx:** fix ParseTagKeyAndFlags
* **textx:** use ParseFlags for tag parsing
* **textx:** scan MarshalText implements from both interface and addr.interface
* **textx:** fix textx.UnmarshalText
* **textx:** fix text.UnmarshalText
* **textx:** fix textx.UnmarshalText to check if the value can be set

### Refactor

* move typex to xoctpus/typex
* remove universe
* move x/typex/internal x/internal
* rename module, migrated to xoctopus
* **reflectx:** add struct tag parsing to Flags
* **stringsx:** move x/misc/stringsx => x/stringsx

### Style

* code format

### Test

* add tests for fuctions with multi type parameters
* fix must unit test
* **reflectx:** update unit test
* **reflextx:** fix unit test
* **resultx:** add unit test
* **resultx:** update unit tests
* **textx:** fix texts unit tests

