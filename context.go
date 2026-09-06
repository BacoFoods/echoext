package echoext

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Integer is the set of integer types Context.ParamNum can parse into.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Context extends echo.Context with type-safe accessors. Every echo.Context
// method is promoted from the embedded interface, so a Context can be used
// anywhere an echo.Context is expected.
type Context struct {
	echo.Context
}

// Value returns the value stored under key as a T. It returns the zero value
// of T when the key is absent or holds a different type.
//
//	userID := c.Value[int]("user_id")
//	name := c.Value[string]("user_name")
func (c Context) Value[T any](key string) T {
	v, _ := c.Get(key).(T)
	return v
}

// ParamNum returns the path parameter name parsed as a T. It returns 0 when
// the parameter is missing, unparseable, or out of range for T.
//
//	id := c.ParamNum[int64]("id")
func (c Context) ParamNum[T Integer](name string) T {
	return parseNum[T](c.Param(name))
}

// parseNum parses s into T, returning 0 rather than an error on failure so
// callers can use the result directly. Values that do not survive the round
// trip through T (that is, they overflow it) are treated as failures.
func parseNum[T Integer](s string) T {
	var zero T

	// zero-1 wraps around to the maximum value on unsigned types, so this
	// distinguishes signed from unsigned without reflection.
	if zero-1 < zero {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}

		if n := T(v); int64(n) == v {
			return n
		}

		return 0
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}

	if n := T(v); uint64(n) == v {
		return n
	}

	return 0
}

// BindValidate binds the request body into i and then validates it.
func (c Context) BindValidate(i any) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	return c.Validate(i)
}
