package controllers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/goravel/framework/contracts/filesystem"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/contracts/validation"
)

type capturedResponse struct {
	kind         string
	status       int
	jsonBody     any
	redirectURL  string
	viewTemplate string
	viewData     any
}

type fakeHTTPContext struct {
	ctx  context.Context
	req  *fakeHTTPRequest
	resp *fakeHTTPResponse
}

func newFakeHTTPContext(req *fakeHTTPRequest) *fakeHTTPContext {
	return &fakeHTTPContext{
		ctx:  context.Background(),
		req:  req,
		resp: &fakeHTTPResponse{},
	}
}

func (c *fakeHTTPContext) Deadline() (time.Time, bool)    { return c.ctx.Deadline() }
func (c *fakeHTTPContext) Done() <-chan struct{}            { return c.ctx.Done() }
func (c *fakeHTTPContext) Err() error                      { return c.ctx.Err() }
func (c *fakeHTTPContext) Value(key any) any               { return c.ctx.Value(key) }
func (c *fakeHTTPContext) Context() context.Context        { return c.ctx }
func (c *fakeHTTPContext) WithContext(ctx context.Context) { c.ctx = ctx }
func (c *fakeHTTPContext) WithValue(key any, value any)    { c.ctx = context.WithValue(c.ctx, key, value) }

func (c *fakeHTTPContext) Request() contractshttp.ContextRequest   { return c.req }
func (c *fakeHTTPContext) Response() contractshttp.ContextResponse { return c.resp }

type fakeHTTPRequest struct {
	queryParams map[string]string
	headers     map[string]string
	bindResult  any
	bindErr     error
}

func (r *fakeHTTPRequest) Query(key string, defaultValue ...string) string {
	if v, ok := r.queryParams[key]; ok {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (r *fakeHTTPRequest) Header(key string, defaultValue ...string) string {
	if v, ok := r.headers[key]; ok {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (r *fakeHTTPRequest) Headers() http.Header  { return http.Header{} }
func (r *fakeHTTPRequest) Ip() string            { return "127.0.0.1" }
func (r *fakeHTTPRequest) Origin() *http.Request { return nil }
func (r *fakeHTTPRequest) Host() string          { return "" }
func (r *fakeHTTPRequest) Url() string           { return "" }
func (r *fakeHTTPRequest) FullUrl() string       { return "" }
func (r *fakeHTTPRequest) Path() string          { return "" }
func (r *fakeHTTPRequest) Method() string        { return "" }

func (r *fakeHTTPRequest) Bind(obj any) error {
	if r.bindErr != nil {
		return r.bindErr
	}
	if r.bindResult == nil {
		return nil
	}
	b, err := json.Marshal(r.bindResult)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func (*fakeHTTPRequest) Cookie(string, ...string) string                                             { panic("not needed") }
func (*fakeHTTPRequest) Name() string                                                                { panic("not needed") }
func (*fakeHTTPRequest) OriginPath() string                                                          { panic("not needed") }
func (*fakeHTTPRequest) Info() contractshttp.Info                                                    { panic("not needed") }
func (*fakeHTTPRequest) All() map[string]any                                                         { panic("not needed") }
func (*fakeHTTPRequest) BindQuery(any) error                                                         { panic("not needed") }
func (*fakeHTTPRequest) Route(string) string                                                         { panic("not needed") }
func (*fakeHTTPRequest) RouteInt(string) int                                                         { panic("not needed") }
func (*fakeHTTPRequest) RouteInt64(string) int64                                                     { panic("not needed") }
func (*fakeHTTPRequest) QueryInt(string, ...int) int                                                 { panic("not needed") }
func (*fakeHTTPRequest) QueryInt64(string, ...int64) int64                                           { panic("not needed") }
func (*fakeHTTPRequest) QueryBool(string, ...bool) bool                                              { panic("not needed") }
func (*fakeHTTPRequest) QueryArray(string) []string                                                  { panic("not needed") }
func (*fakeHTTPRequest) QueryMap(string) map[string]string                                           { panic("not needed") }
func (*fakeHTTPRequest) Queries() map[string]string                                                  { panic("not needed") }
func (*fakeHTTPRequest) HasSession() bool                                                            { panic("not needed") }
func (*fakeHTTPRequest) Session() session.Session                                                    { panic("not needed") }
func (*fakeHTTPRequest) SetSession(session.Session) contractshttp.ContextRequest                     { panic("not needed") }
func (*fakeHTTPRequest) Input(string, ...string) string                                              { panic("not needed") }
func (*fakeHTTPRequest) InputArray(string, ...[]string) []string                                     { panic("not needed") }
func (*fakeHTTPRequest) InputMap(string, ...map[string]any) map[string]any                           { panic("not needed") }
func (*fakeHTTPRequest) InputMapArray(string, ...[]map[string]any) []map[string]any                  { panic("not needed") }
func (*fakeHTTPRequest) InputInt(string, ...int) int                                                 { panic("not needed") }
func (*fakeHTTPRequest) InputInt64(string, ...int64) int64                                           { panic("not needed") }
func (*fakeHTTPRequest) InputBool(string, ...bool) bool                                              { panic("not needed") }
func (*fakeHTTPRequest) File(string) (filesystem.File, error)                                        { panic("not needed") }
func (*fakeHTTPRequest) Files(string) ([]filesystem.File, error)                                     { panic("not needed") }
func (*fakeHTTPRequest) Abort(...int)                                                                 { panic("not needed") }
func (*fakeHTTPRequest) AbortWithStatus(int)                                                         { panic("not needed") }
func (*fakeHTTPRequest) AbortWithStatusJson(int, any)                                                { panic("not needed") }
func (*fakeHTTPRequest) Next()                                                                        { panic("not needed") }
func (*fakeHTTPRequest) Validate(map[string]string, ...validation.Option) (validation.Validator, error) {
	panic("not needed")
}
func (*fakeHTTPRequest) ValidateRequest(contractshttp.FormRequest) (validation.Errors, error) {
	panic("not needed")
}

type fakeHTTPResponse struct {
	captured capturedResponse
}

func (r *fakeHTTPResponse) Json(code int, obj any) contractshttp.AbortableResponse {
	r.captured = capturedResponse{kind: "json", status: code, jsonBody: obj}
	return &fakeAbortableResponse{}
}

func (r *fakeHTTPResponse) Redirect(code int, location string) contractshttp.AbortableResponse {
	r.captured = capturedResponse{kind: "redirect", status: code, redirectURL: location}
	return &fakeAbortableResponse{}
}

func (r *fakeHTTPResponse) View() contractshttp.ResponseView {
	return &fakeResponseView{resp: r}
}

func (*fakeHTTPResponse) Cookie(contractshttp.Cookie) contractshttp.ContextResponse { panic("not needed") }
func (*fakeHTTPResponse) Data(int, string, []byte) contractshttp.AbortableResponse  { panic("not needed") }
func (*fakeHTTPResponse) Download(string, string) contractshttp.Response            { panic("not needed") }
func (*fakeHTTPResponse) File(string) contractshttp.Response                        { panic("not needed") }
func (*fakeHTTPResponse) Header(string, string) contractshttp.ContextResponse       { panic("not needed") }
func (*fakeHTTPResponse) NoContent(...int) contractshttp.AbortableResponse           { panic("not needed") }
func (*fakeHTTPResponse) Origin() contractshttp.ResponseOrigin                      { panic("not needed") }
func (*fakeHTTPResponse) String(int, string, ...any) contractshttp.AbortableResponse { panic("not needed") }
func (*fakeHTTPResponse) Success() contractshttp.ResponseStatus                     { panic("not needed") }
func (*fakeHTTPResponse) Status(int) contractshttp.ResponseStatus                   { panic("not needed") }
func (*fakeHTTPResponse) Stream(int, func(contractshttp.StreamWriter) error) contractshttp.Response {
	panic("not needed")
}
func (*fakeHTTPResponse) Writer() http.ResponseWriter                        { panic("not needed") }
func (*fakeHTTPResponse) WithoutCookie(string) contractshttp.ContextResponse { panic("not needed") }
func (*fakeHTTPResponse) Flush()                                              { panic("not needed") }

type fakeResponseView struct {
	resp *fakeHTTPResponse
}

func (v *fakeResponseView) Make(view string, data ...any) contractshttp.Response {
	var d any
	if len(data) > 0 {
		d = data[0]
	}
	v.resp.captured = capturedResponse{kind: "view", viewTemplate: view, viewData: d}
	return &fakeAbortableResponse{}
}

func (*fakeResponseView) First([]string, ...any) contractshttp.Response { panic("not needed") }

type fakeAbortableResponse struct{}

func (*fakeAbortableResponse) Render() error { return nil }
func (*fakeAbortableResponse) Abort() error  { return nil }
