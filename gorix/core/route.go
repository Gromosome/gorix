package core

type Method string
type Path string
type BasePath string

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
