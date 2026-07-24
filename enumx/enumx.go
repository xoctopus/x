package enumx

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// _U represents the union of supported underlying types for an enumeration.
type _U interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type EncodingArshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type DatabaseValuer interface {
	driver.Valuer
	sql.Scanner
}

// Enum defines an enumeration interface for mapping between `int` and `string`.
// It implements encoding.TextMarshaler encoding.TextUnmarshaler for network
// interaction, and driver.Valuer sql.Scanner for database interaction.
type Enum[E _U] interface {
	// Values presents enum value list
	Values() []E
	// String returns the string identifier (key) of the enumeration.
	String() string
	// Text returns enum description for presents as EnumLabel
	Text() string
	// IsZero check if v is valid
	IsZero() bool

	EncodingArshaler

	DatabaseValuer
}

// CanBeEnum defines an interface for types that can provide value list of their
// enumeration values.
type CanBeEnum interface {
	EnumValues() []any
}

// DriverValueOffset acts as an adaptor between code and database.
// It provides an offset to adjust the enumeration's value during converting
type DriverValueOffset interface {
	Offset() int
}

// Scan parses the database value (src) into an integer and adjusts by offset.
func Scan(src any, offset int) (int, error) {
	switch v := src.(type) {
	case []byte:
		if len(v) == 0 {
			return 0, nil
		}
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return offset, err
		}
		return int(i) - offset, nil
	case string:
		if len(v) == 0 {
			return 0, nil
		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return offset, err
		}
		return int(i) - offset, nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(v).Int()) - offset, nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(v).Uint()) - offset, nil
	default:
		return 0, nil
	}
}

// ParseErrorFor new parsing error for the specific enum type E.
func ParseErrorFor[E _U](from string) error {
	return &parseError[E]{from: from}
}

type parseError[E _U] struct {
	from string
}

func (e *parseError[E]) Error() string {
	return fmt.Sprintf("failed to parse `%s` to %s", e.from, reflect.TypeFor[E]())
}

func (e *parseError[E]) Is(err error) bool {
	target, ok := errors.AsType[*parseError[E]](err)
	return ok && target.from == e.from
}
