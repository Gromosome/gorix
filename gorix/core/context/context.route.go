package context

type Method string
type Path string
type BasePath string
type RouteHandler func(c *Context) (any, error)

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

type RouteInfo struct {
	Method     Method
	Path       string
	Handler    string
	Module     string
	Controller string
}
