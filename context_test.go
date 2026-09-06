package echoext

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func newContext(t *testing.T) Context {
	t.Helper()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	return Context{e.NewContext(req, httptest.NewRecorder())}
}

func TestValue(t *testing.T) {
	c := newContext(t)
	c.Set("int", 42)
	c.Set("str", "hello")
	c.Set("bool", true)
	c.Set("f64", 3.5)
	c.Set("u8", uint8(7))

	if got := c.Value[int]("int"); got != 42 {
		t.Errorf("Value[int] = %d, want 42", got)
	}
	if got := c.Value[string]("str"); got != "hello" {
		t.Errorf("Value[string] = %q, want %q", got, "hello")
	}
	if got := c.Value[bool]("bool"); !got {
		t.Error("Value[bool] = false, want true")
	}
	if got := c.Value[float64]("f64"); got != 3.5 {
		t.Errorf("Value[float64] = %v, want 3.5", got)
	}
	if got := c.Value[uint8]("u8"); got != 7 {
		t.Errorf("Value[uint8] = %d, want 7", got)
	}

	// Wrong type and missing key both yield the zero value, matching the
	// behaviour of the getters this replaced.
	if got := c.Value[string]("int"); got != "" {
		t.Errorf("Value[string] on an int = %q, want empty", got)
	}
	if got := c.Value[int]("absent"); got != 0 {
		t.Errorf("Value[int] on a missing key = %d, want 0", got)
	}
	if got := c.Value[int64]("int"); got != 0 {
		t.Errorf("Value[int64] on an int = %d, want 0 (no implicit widening)", got)
	}
}

func TestParamNum(t *testing.T) {
	c := newContext(t)
	c.SetParamNames("i", "neg", "big", "bad", "empty")
	c.SetParamValues("42", "-7", "300", "abc", "")

	if got := c.ParamNum[int]("i"); got != 42 {
		t.Errorf("ParamNum[int] = %d, want 42", got)
	}
	if got := c.ParamNum[int64]("neg"); got != -7 {
		t.Errorf("ParamNum[int64] = %d, want -7", got)
	}
	if got := c.ParamNum[int16]("big"); got != 300 {
		t.Errorf("ParamNum[int16] = %d, want 300", got)
	}

	// Out of range for the target type -> 0, as the old bitSize-limited
	// strconv calls did.
	if got := c.ParamNum[int8]("big"); got != 0 {
		t.Errorf("ParamNum[int8](300) = %d, want 0", got)
	}
	if got := c.ParamNum[uint8]("big"); got != 0 {
		t.Errorf("ParamNum[uint8](300) = %d, want 0", got)
	}

	// A negative value is not a valid unsigned number.
	if got := c.ParamNum[uint]("neg"); got != 0 {
		t.Errorf("ParamNum[uint](-7) = %d, want 0", got)
	}
	if got := c.ParamNum[uint64]("neg"); got != 0 {
		t.Errorf("ParamNum[uint64](-7) = %d, want 0", got)
	}

	if got := c.ParamNum[int]("bad"); got != 0 {
		t.Errorf("ParamNum[int](abc) = %d, want 0", got)
	}
	if got := c.ParamNum[int]("empty"); got != 0 {
		t.Errorf("ParamNum[int]() = %d, want 0", got)
	}
	if got := c.ParamNum[int]("absent"); got != 0 {
		t.Errorf("ParamNum[int] on a missing param = %d, want 0", got)
	}
}

func TestParamNumBoundaries(t *testing.T) {
	c := newContext(t)
	c.SetParamNames("maxI8", "minI8", "maxU64", "overU64")
	c.SetParamValues("127", "-128", "18446744073709551615", "18446744073709551616")

	if got := c.ParamNum[int8]("maxI8"); got != 127 {
		t.Errorf("ParamNum[int8](127) = %d, want 127", got)
	}
	if got := c.ParamNum[int8]("minI8"); got != -128 {
		t.Errorf("ParamNum[int8](-128) = %d, want -128", got)
	}
	if got := c.ParamNum[uint64]("maxU64"); got != 1<<64-1 {
		t.Errorf("ParamNum[uint64](max) = %d, want %d", got, uint64(1<<64-1))
	}
	if got := c.ParamNum[uint64]("overU64"); got != 0 {
		t.Errorf("ParamNum[uint64](overflow) = %d, want 0", got)
	}
}

// TestNamedTypes covers the ~ in the Integer constraint.
func TestNamedTypes(t *testing.T) {
	type UserID int64

	c := newContext(t)
	c.SetParamNames("id")
	c.SetParamValues("99")

	if got := c.ParamNum[UserID]("id"); got != 99 {
		t.Errorf("ParamNum[UserID] = %d, want 99", got)
	}
}

// TestContextIsEchoContext pins the promotion: a Context satisfies
// echo.Context, so it can be passed anywhere echo expects one.
func TestContextIsEchoContext(t *testing.T) {
	var _ echo.Context = Context{}

	c := newContext(t)
	if c.Request().Method != http.MethodGet {
		t.Errorf("promoted Request() = %q, want GET", c.Request().Method)
	}
}

func TestRoutingEndToEnd(t *testing.T) {
	e := echo.New()
	g := &Group{e.Group("")}

	var (
		gotID   int64
		gotUser string
	)

	mw := func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			c.Set("user", "alice")
			return next(c)
		}
	}

	g.GET("/users/:id", func(c Context) error {
		gotID = c.ParamNum[int64]("id")
		gotUser = c.Value[string]("user")
		return c.NoContent(http.StatusOK)
	}, mw)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/123", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if gotID != 123 {
		t.Errorf("id = %d, want 123", gotID)
	}
	if gotUser != "alice" {
		t.Errorf("user = %q, want alice; middleware state did not survive", gotUser)
	}
}
