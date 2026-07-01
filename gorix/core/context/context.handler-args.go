package context

type BindSource string

const (
	BindSourceParams  BindSource = "params"
	BindSourceQuery   BindSource = "query"
	BindSourceBody    BindSource = "body"
	BindSourceHeaders BindSource = "headers"
)

type BindingArg interface {
	BindSource() BindSource
}

type Params[T any] struct {
	Value T
}

func (Params[T]) BindSource() BindSource {
	return BindSourceParams
}

type Query[T any] struct {
	Value T
}

func (Query[T]) BindSource() BindSource {
	return BindSourceQuery
}

type Body[T any] struct {
	Value T
}

func (Body[T]) BindSource() BindSource {
	return BindSourceBody
}

type Headers[T any] struct {
	Value T
}

func (Headers[T]) BindSource() BindSource {
	return BindSourceHeaders
}
