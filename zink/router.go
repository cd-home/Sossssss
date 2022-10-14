package zink

type Router interface {
	Handle(request *Request)
}
