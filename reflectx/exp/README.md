# About `DeepCopy`

This feature is experimental. Please avoid to import this package in critical
code.

Currently, `DeepCopy` cannot correctly handle composite types with circular
references. For example, the following case:

```golang
package x

type T struct {
	pointer *T
	slice   []T
}
```