package zine

type Router interface {
	Handle(request *Request)
}
