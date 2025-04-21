package echoext

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type Context interface {
	echo.Context
	GetInt64(key string) int64
	GetString(key string) string
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetUint(key string) uint
	GetUint64(key string) uint64
	GetInt32(key string) int32
	GetInt16(key string) int16
	GetInt8(key string) int8
	GetUint32(key string) uint32
	GetUint16(key string) uint16
	GetUint8(key string) uint8
}

var _ Context = (*context)(nil)

type context struct {
	parent echo.Context
}

// Attachment implements Context.
func (c *context) Attachment(file string, name string) error {
	return c.parent.Attachment(file, name)
}

// Bind implements Context.
func (c *context) Bind(i interface{}) error {
	return c.parent.Bind(i)
}

// Blob implements Context.
func (c *context) Blob(code int, contentType string, b []byte) error {
	return c.parent.Blob(code, contentType, b)
}

// Cookie implements Context.
func (c *context) Cookie(name string) (*http.Cookie, error) {
	return c.parent.Cookie(name)
}

// Cookies implements Context.
func (c *context) Cookies() []*http.Cookie {
	return c.parent.Cookies()
}

// Echo implements Context.
func (c *context) Echo() *echo.Echo {
	return c.parent.Echo()
}

// Error implements Context.
func (c *context) Error(err error) {
	c.parent.Error(err)
}

// File implements Context.
func (c *context) File(file string) error {
	return c.parent.File(file)
}

// FormFile implements Context.
func (c *context) FormFile(name string) (*multipart.FileHeader, error) {
	return c.parent.FormFile(name)
}

// FormParams implements Context.
func (c *context) FormParams() (url.Values, error) {
	return c.parent.FormParams()
}

// FormValue implements Context.
func (c *context) FormValue(name string) string {
	return c.parent.FormValue(name)
}

// Get implements Context.
func (c *context) Get(key string) interface{} {
	return c.parent.Get(key)
}

// GetUint16 implements Context.
func (c *context) GetUint16(key string) uint16 {
	if v, ok := c.parent.Get(key).(uint16); ok {
		return v
	}

	return 0
}

// GetUint8 implements Context.
func (c *context) GetUint8(key string) uint8 {
	if v, ok := c.parent.Get(key).(uint8); ok {
		return v
	}

	return 0
}

// HTML implements Context.
func (c *context) HTML(code int, html string) error {
	return c.parent.HTML(code, html)
}

// HTMLBlob implements Context.
func (c *context) HTMLBlob(code int, b []byte) error {
	return c.parent.HTMLBlob(code, b)
}

// Handler implements Context.
func (c *context) Handler() echo.HandlerFunc {
	return c.parent.Handler()
}

// Inline implements Context.
func (c *context) Inline(file string, name string) error {
	return c.parent.Inline(file, name)
}

// IsTLS implements Context.
func (c *context) IsTLS() bool {
	return c.parent.IsTLS()
}

// IsWebSocket implements Context.
func (c *context) IsWebSocket() bool {
	return c.parent.IsWebSocket()
}

// JSON implements Context.
func (c *context) JSON(code int, i interface{}) error {
	return c.parent.JSON(code, i)
}

// JSONBlob implements Context.
func (c *context) JSONBlob(code int, b []byte) error {
	return c.parent.JSONBlob(code, b)
}

// JSONP implements Context.
func (c *context) JSONP(code int, callback string, i interface{}) error {
	return c.parent.JSONP(code, callback, i)
}

// JSONPBlob implements Context.
func (c *context) JSONPBlob(code int, callback string, b []byte) error {
	return c.parent.JSONPBlob(code, callback, b)
}

// JSONPretty implements Context.
func (c *context) JSONPretty(code int, i interface{}, indent string) error {
	return c.parent.JSONPretty(code, i, indent)
}

// Logger implements Context.
func (c *context) Logger() echo.Logger {
	return c.parent.Logger()
}

// MultipartForm implements Context.
func (c *context) MultipartForm() (*multipart.Form, error) {
	return c.parent.MultipartForm()
}

// NoContent implements Context.
func (c *context) NoContent(code int) error {
	return c.parent.NoContent(code)
}

// Param implements Context.
func (c *context) Param(name string) string {
	return c.parent.Param(name)
}

// ParamNames implements Context.
func (c *context) ParamNames() []string {
	return c.parent.ParamNames()
}

// ParamValues implements Context.
func (c *context) ParamValues() []string {
	return c.parent.ParamValues()
}

// Path implements Context.
func (c *context) Path() string {
	return c.parent.Path()
}

// QueryParam implements Context.
func (c *context) QueryParam(name string) string {
	return c.parent.QueryParam(name)
}

// QueryParams implements Context.
func (c *context) QueryParams() url.Values {
	return c.parent.QueryParams()
}

// QueryString implements Context.
func (c *context) QueryString() string {
	return c.parent.QueryString()
}

// RealIP implements Context.
func (c *context) RealIP() string {
	return c.parent.RealIP()
}

// Redirect implements Context.
func (c *context) Redirect(code int, url string) error {
	return c.parent.Redirect(code, url)
}

// Render implements Context.
func (c *context) Render(code int, name string, data interface{}) error {
	return c.parent.Render(code, name, data)
}

// Request implements Context.
func (c *context) Request() *http.Request {
	return c.parent.Request()
}

// Reset implements Context.
func (c *context) Reset(r *http.Request, w http.ResponseWriter) {
	c.parent.Reset(r, w)
}

// Response implements Context.
func (c *context) Response() *echo.Response {
	return c.parent.Response()
}

// Scheme implements Context.
func (c *context) Scheme() string {
	return c.parent.Scheme()
}

// Set implements Context.
func (c *context) Set(key string, val interface{}) {
	c.parent.Set(key, val)
}

// SetCookie implements Context.
func (c *context) SetCookie(cookie *http.Cookie) {
	c.parent.SetCookie(cookie)
}

// SetHandler implements Context.
func (c *context) SetHandler(h echo.HandlerFunc) {
	c.parent.SetHandler(h)
}

// SetLogger implements Context.
func (c *context) SetLogger(l echo.Logger) {
	c.parent.SetLogger(l)
}

// SetParamNames implements Context.
func (c *context) SetParamNames(names ...string) {
	c.parent.SetParamNames(names...)
}

// SetParamValues implements Context.
func (c *context) SetParamValues(values ...string) {
	c.parent.SetParamValues(values...)
}

// SetPath implements Context.
func (c *context) SetPath(p string) {
	c.parent.SetPath(p)
}

// SetRequest implements Context.
func (c *context) SetRequest(r *http.Request) {
	c.parent.SetRequest(r)
}

// SetResponse implements Context.
func (c *context) SetResponse(r *echo.Response) {
	c.parent.SetResponse(r)
}

// Stream implements Context.
func (c *context) Stream(code int, contentType string, r io.Reader) error {
	return c.parent.Stream(code, contentType, r)
}

// String implements Context.
func (c *context) String(code int, s string) error {
	return c.parent.String(code, s)
}

// Validate implements Context.
func (c *context) Validate(i interface{}) error {
	return c.parent.Validate(i)
}

// XML implements Context.
func (c *context) XML(code int, i interface{}) error {
	return c.parent.XML(code, i)
}

// XMLBlob implements Context.
func (c *context) XMLBlob(code int, b []byte) error {
	return c.parent.XMLBlob(code, b)
}

// XMLPretty implements Context.
func (c *context) XMLPretty(code int, i interface{}, indent string) error {
	return c.parent.XMLPretty(code, i, indent)
}

func (c *context) GetInt64(key string) int64 {
	if v, ok := c.parent.Get(key).(int64); ok {
		return v
	}

	return 0
}

func (c *context) GetString(key string) string {
	if v, ok := c.parent.Get(key).(string); ok {
		return v
	}

	return ""
}

func (c *context) GetBool(key string) bool {
	if v, ok := c.parent.Get(key).(bool); ok {
		return v
	}

	return false
}

func (c *context) GetFloat64(key string) float64 {
	if v, ok := c.parent.Get(key).(float64); ok {
		return v
	}

	return 0
}

func (c *context) GetInt(key string) int {
	if v, ok := c.parent.Get(key).(int); ok {
		return v
	}

	return 0
}

func (c *context) GetUint(key string) uint {
	if v, ok := c.parent.Get(key).(uint); ok {
		return v
	}

	return 0
}

func (c *context) GetUint64(key string) uint64 {
	if v, ok := c.parent.Get(key).(uint64); ok {
		return v
	}

	return 0
}

func (c *context) GetInt32(key string) int32 {
	if v, ok := c.parent.Get(key).(int32); ok {
		return v
	}

	return 0
}

func (c *context) GetInt16(key string) int16 {
	if v, ok := c.parent.Get(key).(int16); ok {
		return v
	}

	return 0
}

func (c *context) GetInt8(key string) int8 {
	if v, ok := c.parent.Get(key).(int8); ok {
		return v
	}

	return 0
}

func (c *context) GetUint32(key string) uint32 {
	if v, ok := c.parent.Get(key).(uint32); ok {
		return v
	}

	return 0
}
