package typex

import "reflect"

func DerefPointer(t Type) Type {
	elem := t
	for elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	return elem
}
